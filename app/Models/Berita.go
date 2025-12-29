package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Berita represents a news article
type Berita struct {
	ID          int64      `db:"id" json:"id"`
	Slug        string     `db:"slug" json:"slug"`
	Title       string     `db:"title" json:"title"`
	Excerpt     string     `db:"excerpt" json:"excerpt"`
	Content     string     `db:"content" json:"content"`
	ImageURL    string     `db:"image_url" json:"image_url"`
	Category    string     `db:"category" json:"category"`
	Author      string     `db:"author" json:"author"`
	Status      string     `db:"status" json:"status"`
	Views       int64      `db:"views" json:"views"`
	PublishedAt *time.Time `db:"published_at" json:"published_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	Tags        []string   `db:"-" json:"tags,omitempty"`
}

// Create creates a new berita record
func (b *Berita) Create(db *sqlx.DB) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	if b.Status == "published" && b.PublishedAt == nil {
		now := time.Now()
		b.PublishedAt = &now
	}

	query := `
		INSERT INTO berita (slug, title, excerpt, content, image_url, category, author, status, views, published_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, b.Slug, b.Title, b.Excerpt, b.Content, b.ImageURL, b.Category, b.Author, b.Status, b.Views, b.PublishedAt, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}

// FindBySlug finds a berita by slug (excluding deleted)
func FindBeritaBySlug(db *sqlx.DB, slug string) (*Berita, error) {
	berita := &Berita{}
	query := `
		SELECT id, slug, title, excerpt, content, image_url, category, author, status, views, published_at, created_at, updated_at, deleted_at 
		FROM berita 
		WHERE slug = ? AND deleted_at IS NULL
	`
	err := db.Get(berita, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return berita, nil
}

// FindByID finds a berita by ID (excluding deleted)
func FindBeritaByID(db *sqlx.DB, id int64) (*Berita, error) {
	berita := &Berita{}
	query := `
		SELECT id, slug, title, excerpt, content, image_url, category, author, status, views, published_at, created_at, updated_at, deleted_at 
		FROM berita 
		WHERE id = ? AND deleted_at IS NULL
	`
	err := db.Get(berita, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return berita, nil
}

// GetAllBerita retrieves all berita with filters and pagination
func GetAllBerita(db *sqlx.DB, filters map[string]interface{}, offset int, limit int) ([]Berita, int64, error) {
	var berita []Berita

	query := `SELECT id, slug, title, excerpt, content, image_url, category, author, status, views, published_at, created_at, updated_at, deleted_at FROM berita WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM berita WHERE deleted_at IS NULL`

	// Build WHERE clause based on filters
	args := []interface{}{}
	if category, ok := filters["category"].(string); ok && category != "" {
		query += ` AND category = ?`
		countQuery += ` AND category = ?`
		args = append(args, category)
	}
	if author, ok := filters["author"].(string); ok && author != "" {
		query += ` AND author = ?`
		countQuery += ` AND author = ?`
		args = append(args, author)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query += ` AND status = ?`
		countQuery += ` AND status = ?`
		args = append(args, status)
	}

	// Get total count
	var total int64
	err := db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Add sorting and pagination
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	paginationArgs := append(args, limit, offset)

	err = db.Select(&berita, query, paginationArgs...)
	if err != nil {
		return nil, 0, err
	}

	return berita, total, nil
}

// Update updates a berita record
func (b *Berita) Update(db *sqlx.DB) error {
	b.UpdatedAt = time.Now()
	query := `
		UPDATE berita 
		SET slug = ?, title = ?, excerpt = ?, content = ?, image_url = ?, category = ?, author = ?, status = ?, published_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := db.Exec(query, b.Slug, b.Title, b.Excerpt, b.Content, b.ImageURL, b.Category, b.Author, b.Status, b.PublishedAt, b.UpdatedAt, b.ID)
	return err
}

// Delete soft deletes a berita record
func (b *Berita) Delete(db *sqlx.DB) error {
	now := time.Now()
	b.DeletedAt = &now
	b.UpdatedAt = now
	query := `UPDATE berita SET deleted_at = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, b.DeletedAt, b.UpdatedAt, b.ID)
	return err
}

// GetCategories retrieves all unique categories
func GetBeritaCategories(db *sqlx.DB) ([]string, error) {
	var categories []string
	query := `SELECT DISTINCT category FROM berita WHERE deleted_at IS NULL AND category IS NOT NULL ORDER BY category`
	err := db.Select(&categories, query)
	return categories, err
}
