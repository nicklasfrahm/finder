package finder

import (
	"crypto/sha512"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
)

// Index builds an index of all files and folders in the specified directory.
func Index(folder string) error {
	// Open the SQLite database.
	db, err := sqlx.Open("sqlite3", "finder.sqlite3")
	if err != nil {
		return err
	}

	// Set up database schema.
	if err := CreateSchema(db); err != nil {
		return err
	}

	// Attempt to resume previous indexing operation.
	checkpoint := Resume(db)
	if checkpoint != "" {
		fmt.Println("Resuming indexing from checkpoint:", checkpoint)
		folder = checkpoint
	}

	// Find total folder size.
	total := int64(0)
	if err := filepath.WalkDir(folder, func(relativePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() {
			// Get file or folder information.
			info, err := entry.Info()
			if err != nil {
				return err
			}

			total += info.Size()
		}

		return err
	}); err != nil {
		return err
	}
	fmt.Printf("Folder size: %s\n", humanizeBytes(total))

	indexingProgress := progressbar.NewOptions64(total,
		progressbar.OptionSetDescription("Indexing files..."),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
	)

	// Walk directories to build index.
	if err := filepath.WalkDir(folder, func(relativePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get absolute path.
		path, err := filepath.Abs(relativePath)
		if err != nil {
			return err
		}
		file := NewFile(path)

		// Get file or folder information.
		info, err := entry.Info()
		if err != nil {
			return err
		}

		// Set file or folder information.
		file.Size = info.Size()
		file.IsDir = info.IsDir()
		file.TimeModified = info.ModTime()
		file.TimeIndexed = time.Now()

		// Get additional information for files.
		if !file.IsDir {
			// Detect MIME type.
			mtype, err := mimetype.DetectFile(path)
			if err != nil {
				return err
			}
			file.MIMEType.String = mtype.String()

			// Get file hash.
			hash, err := Hash(path)
			if err != nil {
				return err
			}
			file.Hash.String = hash

			// TODO: Get file creation time for images from EXIF.

			// Update progress.
			indexingProgress.Add64(file.Size)
		}

		// TODO: Set file or folder parent.

		// Persist file or folder information.
		if err := file.Persist(db); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	// Update progress.
	indexingProgress.Set64(total)
	fmt.Println()

	return nil
}

// Resume attempts to resume an indexing operation.
// It returns the path of the last indexed file or folder.
func Resume(db *sqlx.DB) string {
	// TODO: Implement this.

	return ""
}

// Hash returns the hash of the file.
func Hash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create hash without loading entire file into memory.
	// io.Copy will read the file in chunks of 128KiB and
	// update the hash for every chunk.
	hash := sha512.New()
	buf := make([]byte, 128*1024)
	if _, err := io.CopyBuffer(hash, file, buf); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// humanizeBytes returns a human-readable string of the specified size.
func humanizeBytes(size int64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	unit := math.Floor(math.Log2(float64(size)) / 10)

	return fmt.Sprintf("%.2f %s", float64(size)/math.Pow(2, unit*10), units[int(unit)])
}
