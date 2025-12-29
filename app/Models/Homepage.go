package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Homepage represents the homepage configuration (single row)
type Homepage struct {
	ID        int64      `db:"id" json:"id"`
	Hero      string     `db:"hero" json:"hero"` // JSON
	Stats     string     `db:"stats" json:"stats"` // JSON array
	Features  string     `db:"features" json:"features"` // JSON array
	SEO       string     `db:"seo" json:"seo"` // JSON
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

// GetHomepage retrieves the homepage configuration
func GetHomepage(db *sqlx.DB) (*Homepage, error) {
	homepage := &Homepage{}
	query := `SELECT id, hero, stats, features, seo, updated_at FROM homepage LIMIT 1`
	err := db.Get(homepage, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return homepage, nil
}

// Update updates the homepage configuration
func (h *Homepage) Update(db *sqlx.DB) error {
	h.UpdatedAt = time.Now()
	
	// Check if record exists
	existing := &Homepage{}
	query := `SELECT id FROM homepage LIMIT 1`
	err := db.Get(existing, query)
	
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	
	if err == sql.ErrNoRows {
		// Insert new record
		query := `INSERT INTO homepage (hero, stats, features, seo, updated_at) VALUES (?, ?, ?, ?, ?)`
		result, err := db.Exec(query, h.Hero, h.Stats, h.Features, h.SEO, h.UpdatedAt)
		if err != nil {
			return err
		}
		
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		h.ID = id
		return nil
	}
	
	// Update existing record
	query = `UPDATE homepage SET hero = ?, stats = ?, features = ?, seo = ?, updated_at = ? WHERE id = ?`
	_, err = db.Exec(query, h.Hero, h.Stats, h.Features, h.SEO, h.UpdatedAt, existing.ID)
	return err
}

// ProfilOrganisasi represents the organization profile
type ProfilOrganisasi struct {
	ID        int64      `db:"id" json:"id"`
	VisiMisi  string     `db:"visi_misi" json:"visi_misi"`
	Sejarah   string     `db:"sejarah" json:"sejarah"`
	AdArt     string     `db:"ad_art" json:"ad_art"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}

// GetProfilOrganisasi retrieves the organization profile
func GetProfilOrganisasi(db *sqlx.DB) (*ProfilOrganisasi, error) {
	profil := &ProfilOrganisasi{}
	query := `SELECT id, visi_misi, sejarah, ad_art, created_at, updated_at FROM profil_organisasi LIMIT 1`
	err := db.Get(profil, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return profil, nil
}

// Update updates the organization profile
func (p *ProfilOrganisasi) Update(db *sqlx.DB) error {
	p.UpdatedAt = time.Now()
	
	// Check if record exists
	existing := &ProfilOrganisasi{}
	query := `SELECT id FROM profil_organisasi LIMIT 1`
	err := db.Get(existing, query)
	
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	
	if err == sql.ErrNoRows {
		// Insert new record
		p.CreatedAt = time.Now()
		query := `INSERT INTO profil_organisasi (visi_misi, sejarah, ad_art, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
		result, err := db.Exec(query, p.VisiMisi, p.Sejarah, p.AdArt, p.CreatedAt, p.UpdatedAt)
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
	
	// Update existing record
	query = `UPDATE profil_organisasi SET visi_misi = ?, sejarah = ?, ad_art = ?, updated_at = ? WHERE id = ?`
	_, err = db.Exec(query, p.VisiMisi, p.Sejarah, p.AdArt, p.UpdatedAt, existing.ID)
	return err
}
