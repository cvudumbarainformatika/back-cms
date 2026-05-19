package models

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Thumbnail struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	YoutubeURL  string    `db:"youtube_url" json:"youtube_url"`
	Description string    `db:"description" json:"description"`

	// Deprecated (single category) - tetap ada untuk sementara
	Category   string   `db:"category" json:"category,omitempty"`

	// Multi Category (yang baru dipakai)
	Categories []string `db:"-" json:"categories"`

	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// ====================== HELPER FUNCTIONS ======================

// Simpan kategori ke tabel junction
func (t *Thumbnail) SaveCategories(db *sqlx.DB) error {
	if len(t.Categories) == 0 {
		return nil
	}

	// Hapus semua kategori lama dulu
	_, err := db.Exec(`DELETE FROM thumbnail_categories WHERE thumbnail_id = ?`, t.ID)
	if err != nil {
		return err
	}

	// Insert kategori baru
	query := `INSERT INTO thumbnail_categories (thumbnail_id, category) VALUES (?, ?)`
	for _, cat := range t.Categories {
		cat = strings.ToUpper(strings.TrimSpace(cat))
		if cat != "" {
			_, err = db.Exec(query, t.ID, cat)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ====================== MAIN FUNCTIONS ======================

func GetAllThumbnails(db *sqlx.DB, category string, search string, offset int, limit int) ([]Thumbnail, int64, error) {
	var items []Thumbnail
	var total int64

	baseQuery := `FROM thumbnails t`
	whereClause := ""
	var args []interface{}

	conditions := []string{}

	if category != "" {
		conditions = append(conditions, `EXISTS (SELECT 1 FROM thumbnail_categories tc WHERE tc.thumbnail_id = t.id AND tc.category = ?)`)
		args = append(args, strings.ToUpper(category))
	}
	if search != "" {
		conditions = append(conditions, `(t.title LIKE ? OR t.description LIKE ?)`)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	if len(conditions) > 0 {
		whereClause = ` WHERE ` + strings.Join(conditions, ` AND `)
	}

	// Count total
	countQuery := `SELECT COUNT(DISTINCT t.id) ` + baseQuery + whereClause
	err := db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Get data + join categories
	query := `
		SELECT t.*,
		       COALESCE((
		           SELECT GROUP_CONCAT(tc.category)
		           FROM thumbnail_categories tc 
		           WHERE tc.thumbnail_id = t.id
		       ), '') as categories_str
		` + baseQuery + whereClause + ` 
		ORDER BY t.created_at DESC 
		LIMIT ? OFFSET ?`

	paginationArgs := append(args, limit, offset)

	rows, err := db.Queryx(query, paginationArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Thumbnail
		var categoriesStr string

		err = rows.Scan(
			&t.ID, &t.Title, &t.Category, &t.YoutubeURL, &t.Description,
			&t.CreatedAt, &t.UpdatedAt, &categoriesStr,
		)
		if err != nil {
			continue
		}

		// Parse categories string ke slice
		if categoriesStr != "" {
			t.Categories = strings.Split(categoriesStr, ",")
		}

		items = append(items, t)
	}

	return items, total, nil
}

func (t *Thumbnail) Create(db *sqlx.DB) error {
	query := `INSERT INTO thumbnails (title, category, youtube_url, description, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	res, err := db.Exec(query, t.Title, "", t.YoutubeURL, t.Description, now, now)
	if err != nil {
		return err
	}

	id, _ := res.LastInsertId()
	t.ID = id

	// Simpan multi categories
	return t.SaveCategories(db)
}

func (t *Thumbnail) Update(db *sqlx.DB) error {
	query := `UPDATE thumbnails SET title = ?, category = '', youtube_url = ?, description = ?, updated_at = ? WHERE id = ?`
	
	_, err := db.Exec(query, t.Title, t.YoutubeURL, t.Description, time.Now(), t.ID)
	if err != nil {
		return err
	}

	return t.SaveCategories(db)
}

func DeleteThumbnail(db *sqlx.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM thumbnails WHERE id = ?`, id)
	return err
}

// Delete semua thumbnail dalam satu kategori (sesuai request client sebelumnya)
func DeleteThumbnailsByCategory(db *sqlx.DB, category string) error {
	category = strings.ToUpper(strings.TrimSpace(category))
	
	_, err := db.Exec(`
		DELETE FROM thumbnails 
		WHERE id IN (SELECT thumbnail_id FROM thumbnail_categories WHERE category = ?)`, category)
	return err
}

func GetThumbnailCategories(db *sqlx.DB) ([]string, error) {
	var cats []string
	err := db.Select(&cats, `
		SELECT DISTINCT category 
		FROM thumbnail_categories 
		ORDER BY category ASC`)
	return cats, err
}

// Bonus: Ambil satu thumbnail beserta kategorinya
func FindThumbnailByID(db *sqlx.DB, id int64) (*Thumbnail, error) {
	var t Thumbnail
	var categoriesStr string

	err := db.QueryRowx(`
		SELECT t.*, 
		       COALESCE((
		           SELECT GROUP_CONCAT(tc.category)
		           FROM thumbnail_categories tc 
		           WHERE tc.thumbnail_id = t.id
		       ), '') as categories_str
		FROM thumbnails t 
		WHERE t.id = ?`, id).Scan(
		&t.ID, &t.Title, &t.Category, &t.YoutubeURL, &t.Description,
		&t.CreatedAt, &t.UpdatedAt, &categoriesStr,
	)
	if err != nil {
		return nil, err
	}

	if categoriesStr != "" {
		t.Categories = strings.Split(categoriesStr, ",")
	}

	return &t, nil
}