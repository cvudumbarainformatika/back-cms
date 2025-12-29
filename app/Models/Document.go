package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Document represents a user document (STR, SIP, etc.)
type Document struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"user_id"`
	Name       string     `db:"name" json:"name"`
	Type       string     `db:"type" json:"type"`
	ValidUntil *time.Time `db:"valid_until" json:"valid_until"`
	Status     string     `db:"status" json:"status"`
	FileURL    string     `db:"file_url" json:"file_url"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// Create creates a new document record
func (d *Document) Create(db *sqlx.DB) error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()

	query := `
		INSERT INTO documents (user_id, name, type, valid_until, status, file_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, d.UserID, d.Name, d.Type, d.ValidUntil, d.Status, d.FileURL, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	d.ID = id
	return nil
}

// FindByID finds a document by ID (excluding deleted)
func FindDocumentByID(db *sqlx.DB, id int64) (*Document, error) {
	document := &Document{}
	query := `
		SELECT id, user_id, name, type, valid_until, status, file_url, created_at, updated_at, deleted_at 
		FROM documents 
		WHERE id = ? AND deleted_at IS NULL
	`
	err := db.Get(document, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return document, nil
}

// GetUserDocuments retrieves all documents for a user
func GetUserDocuments(db *sqlx.DB, userID int64, filters map[string]interface{}, offset int, limit int) ([]Document, int64, error) {
	var documents []Document

	query := `SELECT id, user_id, name, type, valid_until, status, file_url, created_at, updated_at, deleted_at FROM documents WHERE user_id = ? AND deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM documents WHERE user_id = ? AND deleted_at IS NULL`

	args := []interface{}{userID}
	countArgs := []interface{}{userID}

	if docType, ok := filters["type"].(string); ok && docType != "" {
		query += ` AND type = ?`
		countQuery += ` AND type = ?`
		args = append(args, docType)
		countArgs = append(countArgs, docType)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query += ` AND status = ?`
		countQuery += ` AND status = ?`
		args = append(args, status)
		countArgs = append(countArgs, status)
	}

	// Get total count
	var total int64
	err := db.Get(&total, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// Add sorting and pagination
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	paginationArgs := append(args, limit, offset)

	err = db.Select(&documents, query, paginationArgs...)
	if err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}

// Update updates a document record
func (d *Document) Update(db *sqlx.DB) error {
	d.UpdatedAt = time.Now()
	query := `
		UPDATE documents 
		SET name = ?, type = ?, valid_until = ?, status = ?, file_url = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := db.Exec(query, d.Name, d.Type, d.ValidUntil, d.Status, d.FileURL, d.UpdatedAt, d.ID)
	return err
}

// Delete soft deletes a document record
func (d *Document) Delete(db *sqlx.DB) error {
	now := time.Now()
	d.DeletedAt = &now
	d.UpdatedAt = now
	query := `UPDATE documents SET deleted_at = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, d.DeletedAt, d.UpdatedAt, d.ID)
	return err
}
