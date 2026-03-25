package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Greeting represents a holiday or special occasion greeting
type Greeting struct {
	ID        int64      `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	Content   string     `db:"content" json:"content"`
	ImageURL  string     `db:"image_url" json:"image_url"`
	IsActive  bool       `db:"is_active" json:"is_active"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// Create creates a new greeting record
func (g *Greeting) Create(db *sqlx.DB) error {
	now := time.Now()
	g.CreatedAt = now
	g.UpdatedAt = now
	
	query := `
		INSERT INTO greetings (title, content, image_url, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, g.Title, g.Content, g.ImageURL, g.IsActive, g.CreatedAt, g.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	g.ID = id
	return nil
}

// FindGreetingByID finds a greeting by ID
func FindGreetingByID(db *sqlx.DB, id int64) (*Greeting, error) {
	greeting := &Greeting{}
	query := `SELECT id, title, content, image_url, is_active, created_at, updated_at, deleted_at FROM greetings WHERE id = ? AND deleted_at IS NULL`
	err := db.Get(greeting, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return greeting, nil
}

// GetAllGreetings retrieves all greetings with pagination
func GetAllGreetings(db *sqlx.DB, offset int, limit int, search string) ([]Greeting, int64, error) {
	var greetings []Greeting
	query := `SELECT id, title, content, image_url, is_active, created_at, updated_at, deleted_at FROM greetings WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM greetings WHERE deleted_at IS NULL`
	
	args := []interface{}{}
	if search != "" {
		query += ` AND title LIKE ?`
		countQuery += ` AND title LIKE ?`
		args = append(args, "%"+search+"%")
	}

	var total int64
	err := db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	paginationArgs := append(args, limit, offset)

	err = db.Select(&greetings, query, paginationArgs...)
	if err != nil {
		return nil, 0, err
	}

	return greetings, total, nil
}

// Update updates a greeting record
func (g *Greeting) Update(db *sqlx.DB) error {
	g.UpdatedAt = time.Now()
	query := `
		UPDATE greetings 
		SET title = ?, content = ?, image_url = ?, is_active = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := db.Exec(query, g.Title, g.Content, g.ImageURL, g.IsActive, g.UpdatedAt, g.ID)
	return err
}

// Delete soft deletes a greeting record
func (g *Greeting) Delete(db *sqlx.DB) error {
	now := time.Now()
	g.DeletedAt = &now
	g.UpdatedAt = now
	query := `UPDATE greetings SET deleted_at = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, g.DeletedAt, g.UpdatedAt, g.ID)
	return err
}
