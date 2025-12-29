package models

import (
	"database/sql"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/jmoiron/sqlx"
)

// Example represents an example model
// This is a template model showing best practices
type Example struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	Email     string `db:"email" json:"email"`
	Status    string `db:"status" json:"status"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}

// Create inserts a new example into the database
func (e *Example) Create(db *sqlx.DB) error {
	query := `INSERT INTO examples (name, email, status, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?)`
	
	now := utils.GetCurrentTimeString()
	result, err := db.Exec(query, e.Name, e.Email, e.Status, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	e.ID = int(id)
	e.CreatedAt = now
	e.UpdatedAt = now
	return nil
}

// FindExampleByID finds an example by ID
func FindExampleByID(db *sqlx.DB, id int) (*Example, error) {
	var example Example
	query := `SELECT id, name, email, status, created_at, updated_at FROM examples WHERE id = ?`
	
	err := db.Get(&example, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &example, nil
}

// GetExamples retrieves examples with filtering and pagination
func GetExamples(db *sqlx.DB, params utils.FilterParams) ([]Example, int64, error) {
	var examples []Example
	var total int64

	// Base query
	query := `SELECT id, name, email, status, created_at, updated_at FROM examples WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM examples WHERE 1=1`

	var args []interface{}

	// Apply search filter
	if params.Q != "" {
		filter := " AND (name LIKE ? OR email LIKE ?)"
		query += filter
		countQuery += filter
		searchTerm := "%" + params.Q + "%"
		args = append(args, searchTerm, searchTerm)
	}

	// Get total count
	err := db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderBy := "id"
	if params.OrderBy == "name" {
		orderBy = "name"
	}

	sort := "DESC"
	if params.Sort == "asc" {
		sort = "ASC"
	}

	query += " ORDER BY " + orderBy + " " + sort

	// Apply pagination
	offset := (params.Page - 1) * params.PerPage
	query += " LIMIT ? OFFSET ?"
	args = append(args, params.PerPage, offset)

	err = db.Select(&examples, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return examples, total, nil
}

// Update modifies an example record in the database
func (e *Example) Update(db *sqlx.DB) error {
	query := `UPDATE examples SET name = ?, email = ?, status = ?, updated_at = ? WHERE id = ?`
	now := utils.GetCurrentTimeString()
	_, err := db.Exec(query, e.Name, e.Email, e.Status, now, e.ID)
	return err
}

// Delete removes an example from the database by ID
func DeleteExampleByID(db *sqlx.DB, id int) error {
	query := `DELETE FROM examples WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
