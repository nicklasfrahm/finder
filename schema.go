package finder

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// File represents an indexed file or folder in the file system.
type File struct {
	ID           string         `db:"id"`
	ParentID     sql.NullString `db:"parent_id"`
	Hash         sql.NullString `db:"hash"`
	Path         string         `db:"path"`
	Size         int64          `db:"size"`
	IsDir        bool           `db:"is_dir"`
	MIMEType     sql.NullString `db:"mime_type"`
	TimeModified time.Time      `db:"time_modified"`
	TimeCreated  time.Time      `db:"time_created"`
	TimeIndexed  time.Time      `db:"time_indexed"`
}

// NewFile creates a new File instance.
func NewFile(path string) *File {
	return &File{
		ID:   uuid.New().String(),
		Path: path,
	}
}

// Persist writes the file indexing information to the database.
func (file *File) Persist(db *sqlx.DB) error {
	_, err := db.NamedExec(`
    INSERT INTO files (
      id,
      parent_id,
      hash,
      path,
      size,
      is_dir,
      mime_type,
      time_modified,
      time_created,
      time_indexed
    )
    VALUES (
      :id,
      :parent_id,
      :hash,
      :path,
      :size,
      :is_dir,
      :mime_type,
      :time_modified,
      :time_created,
      :time_indexed
    )
    ON CONFLICT (path) DO UPDATE SET
      parent_id = :parent_id,
      hash = :hash,
      path = :path,
      size = :size,
      is_dir = :is_dir,
      mime_type = :mime_type,
      time_modified = :time_modified,
      time_created = :time_created,
      time_indexed = :time_indexed
  `, file)
	return err
}

// CreateSchema sets up the database tables to store the index
// of files and folders and additional application state.
func CreateSchema(db *sqlx.DB) error {
	// Create the tables.
	if _, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS files (
      id TEXT NOT NULL,
      parent_id TEXT,
      hash TEXT,
      path TEXT UNIQUE NOT NULL,
      size INTEGER NOT NULL,
      is_dir INTEGER NOT NULL,
      mime_type TEXT,
      time_modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      time_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      time_indexed TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY(id),
      FOREIGN KEY(parent_id) REFERENCES files(id)
    )
  `); err != nil {
		return err
	}

	return nil
}
