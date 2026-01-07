package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

// Homepage represents the homepage content configuration
type Homepage struct {
	ID        int64     `db:"id" json:"id"`
	Content   JSON      `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// JSON is a wrapper for handling JSON data in database
type JSON map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return json.Unmarshal(nil, j) // Should probably error out or just ignore
	}
	return json.Unmarshal(bytes, j)
}

// GetHomepage retrieves the homepage content
func GetHomepage(db *sqlx.DB) (*Homepage, error) {
	var homepage Homepage
	// We assume there's always one record with ID 1, or just the latest one
	query := `SELECT id, content, created_at, updated_at FROM homepage ORDER BY id DESC LIMIT 1`
	err := db.Get(&homepage, query)
	return &homepage, err
}

// UpdateHomepage updates or creates the homepage content
func UpdateHomepage(db *sqlx.DB, content map[string]interface{}) (*Homepage, error) {
	// Convert content to JSON
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	// Upsert query (PostgreSQL specific, for MySQL use ON DUPLICATE KEY UPDATE)
	// Since we are using MySQL (implied by previous context or standard setups), let's use MySQL syntax.
	// But let's check if we just want to insert a new version every time or update single row.
	// Let's assume single row single source of truth structure for CMS simplicity?
	// actually standard is to UPDATE id=1 or INSERT if not exists.

	// Check if exists
	var count int
	db.Get(&count, "SELECT COUNT(*) FROM homepage WHERE id = 1")

	var query string
	if count > 0 {
		query = `UPDATE homepage SET content = ?, updated_at = NOW() WHERE id = 1`
		_, err = db.Exec(query, contentJSON)
	} else {
		// Force ID 1
		query = `INSERT INTO homepage (id, content, created_at, updated_at) VALUES (1, ?, NOW(), NOW())`
		_, err = db.Exec(query, contentJSON)
	}

	if err != nil {
		return nil, err
	}

	return GetHomepage(db)
}
