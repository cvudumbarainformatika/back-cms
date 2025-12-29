package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Agenda represents an event/agenda
type Agenda struct {
	ID              int64      `db:"id" json:"id"`
	Slug            string     `db:"slug" json:"slug"`
	Title           string     `db:"title" json:"title"`
	Description     string     `db:"description" json:"description"`
	Type            string     `db:"type" json:"type"`
	Date            time.Time  `db:"date" json:"date"`
	EndDate         *time.Time `db:"end_date" json:"end_date"`
	IsOnline        bool       `db:"is_online" json:"is_online"`
	Location        string     `db:"location" json:"location"`
	SKP             *float64   `db:"skp" json:"skp"`
	Quota           *int       `db:"quota" json:"quota"`
	RegistrationURL string     `db:"registration_url" json:"registration_url"`
	ImageURL        string     `db:"image_url" json:"image_url"`
	Fee             *float64   `db:"fee" json:"fee"`
	Status          string     `db:"status" json:"status"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

// Create creates a new agenda record
func (a *Agenda) Create(db *sqlx.DB) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	query := `
		INSERT INTO agenda (slug, title, description, type, date, end_date, is_online, location, skp, quota, registration_url, image_url, fee, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, a.Slug, a.Title, a.Description, a.Type, a.Date, a.EndDate, a.IsOnline, a.Location, a.SKP, a.Quota, a.RegistrationURL, a.ImageURL, a.Fee, a.Status, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = id
	return nil
}

// FindBySlug finds an agenda by slug (excluding deleted)
func FindAgendaBySlug(db *sqlx.DB, slug string) (*Agenda, error) {
	agenda := &Agenda{}
	query := `
		SELECT id, slug, title, description, type, date, end_date, is_online, location, skp, quota, registration_url, image_url, fee, status, created_at, updated_at, deleted_at 
		FROM agenda 
		WHERE slug = ? AND deleted_at IS NULL
	`
	err := db.Get(agenda, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return agenda, nil
}

// FindByID finds an agenda by ID (excluding deleted)
func FindAgendaByID(db *sqlx.DB, id int64) (*Agenda, error) {
	agenda := &Agenda{}
	query := `
		SELECT id, slug, title, description, type, date, end_date, is_online, location, skp, quota, registration_url, image_url, fee, status, created_at, updated_at, deleted_at 
		FROM agenda 
		WHERE id = ? AND deleted_at IS NULL
	`
	err := db.Get(agenda, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return agenda, nil
}

// GetAllAgenda retrieves all agenda with filters and pagination
func GetAllAgenda(db *sqlx.DB, filters map[string]interface{}, offset int, limit int) ([]Agenda, int64, error) {
	var agendas []Agenda

	query := `SELECT id, slug, title, description, type, date, end_date, is_online, location, skp, quota, registration_url, image_url, fee, status, created_at, updated_at, deleted_at FROM agenda WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM agenda WHERE deleted_at IS NULL`

	args := []interface{}{}

	if agendaType, ok := filters["type"].(string); ok && agendaType != "" {
		query += ` AND type = ?`
		countQuery += ` AND type = ?`
		args = append(args, agendaType)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query += ` AND status = ?`
		countQuery += ` AND status = ?`
		args = append(args, status)
	}
	if upcoming, ok := filters["upcoming"].(bool); ok {
		if upcoming {
			query += ` AND date >= CURDATE()`
			countQuery += ` AND date >= CURDATE()`
		} else {
			query += ` AND date < CURDATE()`
			countQuery += ` AND date < CURDATE()`
		}
	}

	// Get total count
	var total int64
	err := db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Add sorting and pagination
	query += ` ORDER BY date ASC, created_at DESC LIMIT ? OFFSET ?`
	paginationArgs := append(args, limit, offset)

	err = db.Select(&agendas, query, paginationArgs...)
	if err != nil {
		return nil, 0, err
	}

	return agendas, total, nil
}

// Update updates an agenda record
func (a *Agenda) Update(db *sqlx.DB) error {
	a.UpdatedAt = time.Now()
	query := `
		UPDATE agenda 
		SET slug = ?, title = ?, description = ?, type = ?, date = ?, end_date = ?, is_online = ?, location = ?, skp = ?, quota = ?, registration_url = ?, image_url = ?, fee = ?, status = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := db.Exec(query, a.Slug, a.Title, a.Description, a.Type, a.Date, a.EndDate, a.IsOnline, a.Location, a.SKP, a.Quota, a.RegistrationURL, a.ImageURL, a.Fee, a.Status, a.UpdatedAt, a.ID)
	return err
}

// Delete soft deletes an agenda record
func (a *Agenda) Delete(db *sqlx.DB) error {
	now := time.Now()
	a.DeletedAt = &now
	a.UpdatedAt = now
	query := `UPDATE agenda SET deleted_at = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, a.DeletedAt, a.UpdatedAt, a.ID)
	return err
}

// AgendaRegistration represents a user registration for an agenda
type AgendaRegistration struct {
	ID         int64     `db:"id" json:"id"`
	AgendaID   int64     `db:"agenda_id" json:"agenda_id"`
	UserID     int64     `db:"user_id" json:"user_id"`
	Status     string    `db:"status" json:"status"`
	RegisteredAt time.Time `db:"registered_at" json:"registered_at"`
}

// CreateRegistration creates a new agenda registration
func (ar *AgendaRegistration) Create(db *sqlx.DB) error {
	ar.RegisteredAt = time.Now()
	query := `INSERT INTO agenda_registrations (agenda_id, user_id, status, registered_at) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, ar.AgendaID, ar.UserID, ar.Status, ar.RegisteredAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ar.ID = id
	return nil
}

// GetRegistrations retrieves all registrations for an agenda
func GetAgendaRegistrations(db *sqlx.DB, agendaID int64) ([]AgendaRegistration, error) {
	var registrations []AgendaRegistration
	query := `SELECT id, agenda_id, user_id, status, registered_at FROM agenda_registrations WHERE agenda_id = ? ORDER BY registered_at DESC`
	err := db.Select(&registrations, query, agendaID)
	return registrations, err
}
