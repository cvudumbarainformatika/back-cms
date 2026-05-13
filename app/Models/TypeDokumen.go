package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// TypeDokumen represents document type master
type TypeDokumen struct {
	ID          int64     `db:"id" json:"id"`
	Typedokumen string    `db:"typedokumen" json:"typedokumen"`
	Flag        string    `db:"flag" json:"flag"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// CREATE
func (t *TypeDokumen) Create(db *sqlx.DB) error {

	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	query := `
		INSERT INTO typedokumen (
			typedokumen,
			flag,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?)
	`

	result, err := db.Exec(
		query,
		t.Typedokumen,
		t.Flag,
		t.CreatedAt,
		t.UpdatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	t.ID = id

	return nil
}

// FIND BY ID
func FindTypeDokumenByID(
	db *sqlx.DB,
	id int64,
) (*TypeDokumen, error) {

	data := &TypeDokumen{}

	query := `
		SELECT
			id,
			typedokumen,
			flag,
			created_at,
			updated_at
		FROM typedokumen
		WHERE id = ?
	`

	err := db.Get(data, query, id)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return data, nil
}

// GET ALL
func GetAllTypeDokumen(
	db *sqlx.DB,
	filters map[string]interface{},
	offset int,
	limit int,
) ([]TypeDokumen, int64, error) {

	var items []TypeDokumen

	query := `
		SELECT
			id,
			typedokumen,
			flag,
			created_at,
			updated_at
		FROM typedokumen
		ORDER BY id DESC
		LIMIT ? OFFSET ?
	`

	err := db.Select(&items, query, limit, offset)

	if err != nil {
		return nil, 0, err
	}

	var total int64

	err = db.Get(&total, `
		SELECT COUNT(*)
		FROM typedokumen
	`)

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// UPDATE
func (t *TypeDokumen) Update(db *sqlx.DB) error {

	t.UpdatedAt = time.Now()

	query := `
		UPDATE typedokumen
		SET
			typedokumen = ?,
			flag = ?,
			updated_at = ?
		WHERE id = ?
	`

	_, err := db.Exec(
		query,
		t.Typedokumen,
		t.Flag,
		t.UpdatedAt,
		t.ID,
	)

	return err
}

// DELETE
func (t *TypeDokumen) Delete(db *sqlx.DB) error {

	query := `
		DELETE FROM typedokumen
		WHERE id = ?
	`

	_, err := db.Exec(query, t.ID)

	return err
}
