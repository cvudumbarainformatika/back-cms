package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type TypeArtikel struct {
	ID          int       `db:"id" json:"id"`
	TypeArtikel string    `db:"typeartikel" json:"typeartikel"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Ambil semua data
func GetAllTypeArtikel(db *sqlx.DB) ([]TypeArtikel, error) {
	var items []TypeArtikel

	query := `
		SELECT 
			id,
			typeartikel,
			created_at,
			updated_at
		FROM typeartikel
		ORDER BY id DESC
	`

	err := db.Select(&items, query)

	return items, err
}

// Tambah data baru
func (t *TypeArtikel) Create(db *sqlx.DB) error {
	query := `
		INSERT INTO typeartikel (
			typeartikel,
			created_at,
			updated_at
		)
		VALUES (?, NOW(), NOW())
	`

	result, err := db.Exec(query, t.TypeArtikel)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = int(id)

	return nil
}
