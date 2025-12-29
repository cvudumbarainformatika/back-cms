package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Direktori represents a directory entry (hospital, clinic, institution)
type Direktori struct {
	ID               int64      `db:"id" json:"id"`
	Name             string     `db:"name" json:"name"`
	Type             string     `db:"type" json:"type"`
	Address          string     `db:"address" json:"address"`
	Phone            string     `db:"phone" json:"phone"`
	Email            string     `db:"email" json:"email"`
	Website          string     `db:"website" json:"website"`
	City             string     `db:"city" json:"city"`
	Province         string     `db:"province" json:"province"`
	HasRespirologist bool       `db:"has_respirologist" json:"has_respirologist"`
	Facilities       string     `db:"facilities" json:"facilities"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// Create creates a new direktori record
func (d *Direktori) Create(db *sqlx.DB) error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()

	query := `
		INSERT INTO direktori (name, type, address, phone, email, website, city, province, has_respirologist, facilities, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, d.Name, d.Type, d.Address, d.Phone, d.Email, d.Website, d.City, d.Province, d.HasRespirologist, d.Facilities, d.CreatedAt, d.UpdatedAt)
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

// FindByID finds a direktori by ID (excluding deleted)
func FindDirektoriByID(db *sqlx.DB, id int64) (*Direktori, error) {
	direktori := &Direktori{}
	query := `
		SELECT id, name, type, address, phone, email, website, city, province, has_respirologist, facilities, created_at, updated_at, deleted_at 
		FROM direktori 
		WHERE id = ? AND deleted_at IS NULL
	`
	err := db.Get(direktori, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return direktori, nil
}

// GetAllDirektori retrieves all direktori with filters and pagination
func GetAllDirektori(db *sqlx.DB, filters map[string]interface{}, offset int, limit int) ([]Direktori, int64, error) {
	var direktori []Direktori

	query := `SELECT id, name, type, address, phone, email, website, city, province, has_respirologist, facilities, created_at, updated_at, deleted_at FROM direktori WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM direktori WHERE deleted_at IS NULL`

	args := []interface{}{}

	if province, ok := filters["province"].(string); ok && province != "" {
		query += ` AND province = ?`
		countQuery += ` AND province = ?`
		args = append(args, province)
	}
	if dirType, ok := filters["type"].(string); ok && dirType != "" {
		query += ` AND type = ?`
		countQuery += ` AND type = ?`
		args = append(args, dirType)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		query += ` AND (name LIKE ? OR address LIKE ?)`
		countQuery += ` AND (name LIKE ? OR address LIKE ?)`
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
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

	err = db.Select(&direktori, query, paginationArgs...)
	if err != nil {
		return nil, 0, err
	}

	return direktori, total, nil
}

// Update updates a direktori record
func (d *Direktori) Update(db *sqlx.DB) error {
	d.UpdatedAt = time.Now()
	query := `
		UPDATE direktori 
		SET name = ?, type = ?, address = ?, phone = ?, email = ?, website = ?, city = ?, province = ?, has_respirologist = ?, facilities = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := db.Exec(query, d.Name, d.Type, d.Address, d.Phone, d.Email, d.Website, d.City, d.Province, d.HasRespirologist, d.Facilities, d.UpdatedAt, d.ID)
	return err
}

// Delete soft deletes a direktori record
func (d *Direktori) Delete(db *sqlx.DB) error {
	now := time.Now()
	d.DeletedAt = &now
	d.UpdatedAt = now
	query := `UPDATE direktori SET deleted_at = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, d.DeletedAt, d.UpdatedAt, d.ID)
	return err
}
