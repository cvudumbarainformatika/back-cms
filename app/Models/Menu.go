package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Menu represents a navigation menu item
type Menu struct {
	ID        int64     `db:"id" json:"id"`
	Label     string    `db:"label" json:"label"`
	Slug      string    `db:"slug" json:"slug"`
	To        string    `db:"to" json:"to"`
	Icon      string    `db:"icon" json:"icon"`
	ParentID  *int64    `db:"parent_id" json:"parent_id"`
	Position  string    `db:"position" json:"position"`
	Order     int       `db:"order" json:"order"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	IsFixed   bool      `db:"is_fixed" json:"is_fixed"`
	Roles     string    `db:"roles" json:"roles"` // JSON array
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Children  []Menu    `db:"-" json:"children,omitempty"`
}

// Create creates a new menu record
func (m *Menu) Create(db *sqlx.DB) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	query := "INSERT INTO menus (label, slug, `to`, icon, parent_id, position, `order`, is_active, is_fixed, roles, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := db.Exec(query, m.Label, m.Slug, m.To, m.Icon, m.ParentID, m.Position, m.Order, m.IsActive, m.IsFixed, m.Roles, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

// FindByID finds a menu by ID
func FindMenuByID(db *sqlx.DB, id int64) (*Menu, error) {
	menu := &Menu{}
	query := "SELECT id, label, slug, `to`, icon, parent_id, position, `order`, is_active, is_fixed, roles, created_at, updated_at FROM menus WHERE id = ?"
	err := db.Get(menu, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return menu, nil
}

// GetMenusByPosition retrieves all menus by position
func GetMenusByPosition(db *sqlx.DB, position string) ([]Menu, error) {
	var menus []Menu
	query := "SELECT id, label, slug, `to`, icon, parent_id, position, `order`, is_active, is_fixed, roles, created_at, updated_at FROM menus WHERE position = ? ORDER BY `order` ASC"
	err := db.Select(&menus, query, position)
	return menus, err
}

// GetAllMenus retrieves all menus
func GetAllMenus(db *sqlx.DB) ([]Menu, error) {
	var menus []Menu
	query := "SELECT id, label, slug, `to`, icon, parent_id, position, `order`, is_active, is_fixed, roles, created_at, updated_at FROM menus WHERE is_active = TRUE ORDER BY position ASC, `order` ASC"
	err := db.Select(&menus, query)
	return menus, err
}

// Update updates a menu record
func (m *Menu) Update(db *sqlx.DB) error {
	m.UpdatedAt = time.Now()
	query := "UPDATE menus SET label = ?, slug = ?, `to` = ?, icon = ?, parent_id = ?, position = ?, `order` = ?, is_active = ?, is_fixed = ?, roles = ?, updated_at = ? WHERE id = ?"
	_, err := db.Exec(query, m.Label, m.Slug, m.To, m.Icon, m.ParentID, m.Position, m.Order, m.IsActive, m.IsFixed, m.Roles, m.UpdatedAt, m.ID)
	return err
}

// Delete deletes a menu record
func (m *Menu) Delete(db *sqlx.DB) error {
	query := "DELETE FROM menus WHERE id = ?"
	_, err := db.Exec(query, m.ID)
	return err
}
