package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// User represents a user in the system
type User struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"` // Don't expose password in responses
	Role      string    `db:"role" json:"role"`
	Status    string    `db:"status" json:"status"` // active, inactive, suspended
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CreateUser creates a new user in the database
func (u *User) Create(db *sqlx.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	if u.Status == "" {
		u.Status = "active"
	}
	if u.Role == "" {
		u.Role = "user"
	}

	query := `
		INSERT INTO users (name, email, password, role, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, u.Name, u.Email, u.Password, u.Role, u.Status, u.CreatedAt, u.UpdatedAt)
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

// FindByEmail finds a user by email
func FindByEmail(db *sqlx.DB, email string) (*User, error) {
	user := &User{}
	query := `SELECT id, name, email, password, role, status, created_at, updated_at FROM users WHERE email = ?`
	err := db.Get(user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// FindByID finds a user by ID
func FindByID(db *sqlx.DB, id int64) (*User, error) {
	user := &User{}
	query := `SELECT id, name, email, password, role, status, created_at, updated_at FROM users WHERE id = ?`
	err := db.Get(user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetAll retrieves all users with pagination
func GetAll(db *sqlx.DB, offset int, limit int) ([]User, int64, error) {
	var users []User

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE status != 'deleted'`
	if err := db.Get(&total, countQuery); err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `SELECT id, name, email, password, role, status, created_at, updated_at FROM users WHERE status != 'deleted' ORDER BY created_at DESC LIMIT ? OFFSET ?`
	err := db.Select(&users, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update updates a user's information
func (u *User) Update(db *sqlx.DB) error {
	u.UpdatedAt = time.Now()
	query := `
		UPDATE users 
		SET name = ?, email = ?, role = ?, status = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, u.Name, u.Email, u.Role, u.Status, u.UpdatedAt, u.ID)
	return err
}

// UpdatePassword updates a user's password
func (u *User) UpdatePassword(db *sqlx.DB, newPassword string) error {
	u.UpdatedAt = time.Now()
	query := `UPDATE users SET password = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, newPassword, u.UpdatedAt, u.ID)
	return err
}

// Delete marks a user as deleted (soft delete)
func (u *User) Delete(db *sqlx.DB) error {
	u.Status = "deleted"
	u.UpdatedAt = time.Now()
	query := `UPDATE users SET status = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, u.Status, u.UpdatedAt, u.ID)
	return err
}
