package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Pengurus represents a leadership/staff member
type Pengurus struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Position  string     `db:"position" json:"position"`
	Bidang    string     `db:"bidang" json:"bidang"`
	Level     string     `db:"level" json:"level"`
	Periode   string     `db:"periode" json:"periode"`
	Email     string     `db:"email" json:"email"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// Create creates a new pengurus record
func (p *Pengurus) Create(db *sqlx.DB) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	query := `
		INSERT INTO pengurus (name, position, bidang, level, periode, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, p.Name, p.Position, p.Bidang, p.Level, p.Periode, p.Email, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = id
	return nil
}

// FindByID finds a pengurus by ID (excluding deleted)
func FindPengurusByID(db *sqlx.DB, id int64) (*Pengurus, error) {
	pengurus := &Pengurus{}
	query := `
		SELECT id, name, position, bidang, level, periode, email, created_at, updated_at, deleted_at 
		FROM pengurus 
		WHERE id = ? AND deleted_at IS NULL
	`
	err := db.Get(pengurus, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return pengurus, nil
}

// GetAllPengurus retrieves all pengurus with filters and pagination
func GetAllPengurus(db *sqlx.DB, filters map[string]interface{}, offset int, limit int) ([]Pengurus, int64, error) {
	var pengurus []Pengurus

	query := `SELECT id, name, position, bidang, level, periode, email, created_at, updated_at, deleted_at FROM pengurus WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM pengurus WHERE deleted_at IS NULL`

	args := []interface{}{}

	if level, ok := filters["level"].(string); ok && level != "" {
		query += ` AND level = ?`
		countQuery += ` AND level = ?`
		args = append(args, level)
	}
	if periode, ok := filters["periode"].(string); ok && periode != "" {
		query += ` AND periode = ?`
		countQuery += ` AND periode = ?`
		args = append(args, periode)
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

	err = db.Select(&pengurus, query, paginationArgs...)
	if err != nil {
		return nil, 0, err
	}

	return pengurus, total, nil
}

// Update updates a pengurus record
func (p *Pengurus) Update(db *sqlx.DB) error {
	p.UpdatedAt = time.Now()
	query := `
		UPDATE pengurus 
		SET name = ?, position = ?, bidang = ?, level = ?, periode = ?, email = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := db.Exec(query, p.Name, p.Position, p.Bidang, p.Level, p.Periode, p.Email, p.UpdatedAt, p.ID)
	return err
}

// Delete soft deletes a pengurus record
func (p *Pengurus) Delete(db *sqlx.DB) error {
	now := time.Now()
	p.DeletedAt = &now
	p.UpdatedAt = now
	query := `UPDATE pengurus SET deleted_at = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, p.DeletedAt, p.UpdatedAt, p.ID)
	return err
}
