package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// DynamicContent represents a dynamic content page
type DynamicContent struct {
	ID          int64      `db:"id" json:"id"`
	Slug        string     `db:"slug" json:"slug"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Body        string     `db:"body" json:"body"`
	HTML        string     `db:"html" json:"html"`
	Date        *time.Time `db:"date" json:"date"`
	Image       string     `db:"image" json:"image"` // JSON
	Authors     string     `db:"authors" json:"authors"` // JSON array
	Badge       string     `db:"badge" json:"badge"` // JSON
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// Create creates a new dynamic content record
func (dc *DynamicContent) Create(db *sqlx.DB) error {
	dc.CreatedAt = time.Now()
	dc.UpdatedAt = time.Now()

	query := `
		INSERT INTO dynamic_contents (slug, title, description, body, html, date, image, authors, badge, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, dc.Slug, dc.Title, dc.Description, dc.Body, dc.HTML, dc.Date, dc.Image, dc.Authors, dc.Badge, dc.CreatedAt, dc.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	dc.ID = id
	return nil
}

// FindBySlug finds a dynamic content by slug
func FindDynamicContentBySlug(db *sqlx.DB, slug string) (*DynamicContent, error) {
	content := &DynamicContent{}
	query := `
		SELECT id, slug, title, description, body, html, date, image, authors, badge, created_at, updated_at 
		FROM dynamic_contents 
		WHERE slug = ?
	`
	err := db.Get(content, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return content, nil
}

// FindByID finds a dynamic content by ID
func FindDynamicContentByID(db *sqlx.DB, id int64) (*DynamicContent, error) {
	content := &DynamicContent{}
	query := `
		SELECT id, slug, title, description, body, html, date, image, authors, badge, created_at, updated_at 
		FROM dynamic_contents 
		WHERE id = ?
	`
	err := db.Get(content, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return content, nil
}

// Update updates a dynamic content record
func (dc *DynamicContent) Update(db *sqlx.DB) error {
	dc.UpdatedAt = time.Now()
	query := `
		UPDATE dynamic_contents 
		SET slug = ?, title = ?, description = ?, body = ?, html = ?, date = ?, image = ?, authors = ?, badge = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, dc.Slug, dc.Title, dc.Description, dc.Body, dc.HTML, dc.Date, dc.Image, dc.Authors, dc.Badge, dc.UpdatedAt, dc.ID)
	return err
}

// Delete deletes a dynamic content record
func (dc *DynamicContent) Delete(db *sqlx.DB) error {
	query := `DELETE FROM dynamic_contents WHERE id = ?`
	_, err := db.Exec(query, dc.ID)
	return err
}
