package finder

import (
	"crypto/sha512"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/context"
)

// IndexJob represents a job to index a file or folder.
type IndexJob struct {
	Path  string
	Entry fs.DirEntry
}

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
	total, err := FolderSize(folder)
	if err != nil {
		return err
	}
	fmt.Printf("Folder size: %s\n", humanizeBytes(total))

	progress := progressbar.NewOptions64(total,
		progressbar.OptionSetDescription("Indexing files..."),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
	)

	// Create pool of workers.
	jobs := make(chan IndexJob, runtime.NumCPU())
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case job := <-jobs:
					if err := IndexFile(db, &job, progress); err != nil {
						fmt.Println(err)
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Walk directories to build index.
	if err := filepath.WalkDir(folder, func(relativePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Process indexing job.
		jobs <- IndexJob{
			Path:  relativePath,
			Entry: entry,
		}

		return nil
	}); err != nil {
		return err
	}

	// Shut down workers and close channels.
	cancel()
	close(jobs)

	// Update progress.
	progress.Set64(total)
	fmt.Println()

	return nil
}

// IndexFile indexes a file for a given IndexJob.
func IndexFile(db *sqlx.DB, job *IndexJob, progress *progressbar.ProgressBar) error {
	// Get absolute path.
	path, err := filepath.Abs(job.Path)
	if err != nil {
		return err
	}
	file := NewFile(path)

	// Get file or folder information.
	info, err := job.Entry.Info()
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
		// TODO: Pipe file content into mimetype detection
		// and hasher to avoid loading entire file into memory twice.

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
		progress.Add64(file.Size)
	}

	// TODO: Set file or folder parent.

	// Persist file or folder information.
	return file.Persist(db)
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

// FolderSize returns the size of the specified folder.
func FolderSize(folder string) (int64, error) {
	size := int64(0)

	// Walk directory to get individual file sizes.
	if err := filepath.WalkDir(folder, func(relativePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() {
			// Get file information.
			info, err := entry.Info()
			if err != nil {
				return err
			}

			// Add file size to total.
			size += info.Size()
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return size, nil
}

// humanizeBytes returns a human-readable string of the specified size.
func humanizeBytes(size int64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	unit := math.Floor(math.Log2(float64(size)) / 10)

	return fmt.Sprintf("%.2f %s", float64(size)/math.Pow(2, unit*10), units[int(unit)])
}
