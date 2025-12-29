package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Upload represents an uploaded file
type Upload struct {
	ID           int64      `db:"id" json:"id"`
	UploaderID   *int64     `db:"uploader_id" json:"uploader_id"`
	Filename     string     `db:"filename" json:"filename"`
	OriginalName string     `db:"original_name" json:"original_name"`
	MimeType     string     `db:"mime_type" json:"mime_type"`
	Size         int64      `db:"size" json:"size"`
	URL          string     `db:"url" json:"url"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
}

// Create creates a new upload record
func (u *Upload) Create(db *sqlx.DB) error {
	u.CreatedAt = time.Now()

	query := `
		INSERT INTO uploads (uploader_id, filename, original_name, mime_type, size, url, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, u.UploaderID, u.Filename, u.OriginalName, u.MimeType, u.Size, u.URL, u.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

// FindByID finds an upload by ID
func FindUploadByID(db *sqlx.DB, id int64) (*Upload, error) {
	upload := &Upload{}
	query := `
		SELECT id, uploader_id, filename, original_name, mime_type, size, url, created_at 
		FROM uploads 
		WHERE id = ?
	`
	err := db.Get(upload, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return upload, nil
}

// GetAllUploads retrieves all uploads with pagination
func GetAllUploads(db *sqlx.DB, offset int, limit int) ([]Upload, int64, error) {
	var uploads []Upload

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM uploads`
	if err := db.Get(&total, countQuery); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `SELECT id, uploader_id, filename, original_name, mime_type, size, url, created_at FROM uploads ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err := db.Select(&uploads, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return uploads, total, nil
}

// Delete deletes an upload record
func (u *Upload) Delete(db *sqlx.DB) error {
	query := `DELETE FROM uploads WHERE id = ?`
	_, err := db.Exec(query, u.ID)
	return err
}
