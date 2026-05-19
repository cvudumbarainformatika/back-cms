package models

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Thumbnail struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Category    string    `db:"category" json:"category"`
	YoutubeURL  string    `db:"youtube_url" json:"youtube_url"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func GetAllThumbnails(db *sqlx.DB, category string, search string, offset int, limit int) ([]Thumbnail, int64, error) {
	var items []Thumbnail
	var total int64
	
	baseQuery := `FROM thumbnails`
	whereClause := ""
	var args []interface{}

	conditions := []string{}
	if category != "" {
		conditions = append(conditions, `category = ?`)
		args = append(args, category)
	}
	if search != "" {
		conditions = append(conditions, `(title LIKE ? OR description LIKE ?)`)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	if len(conditions) > 0 {
		whereClause = ` WHERE ` + strings.Join(conditions, ` AND `)
	}

	// Get total count
	err := db.Get(&total, `SELECT COUNT(*) `+baseQuery+whereClause, args...)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated items
	query := `SELECT * ` + baseQuery + whereClause + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	paginationArgs := append(args, limit, offset)

	err = db.Select(&items, query, paginationArgs...)
	return items, total, err
}

func FindThumbnailByID(db *sqlx.DB, id int64) (*Thumbnail, error) {
	var item Thumbnail
	err := db.Get(&item, `SELECT * FROM thumbnails WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (t *Thumbnail) Create(db *sqlx.DB) error {
	query := `INSERT INTO thumbnails (title, category, youtube_url, description) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, t.Title, t.Category, t.YoutubeURL, t.Description)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	t.ID = id
	return nil
}

func (t *Thumbnail) Update(db *sqlx.DB) error {
	query := `UPDATE thumbnails SET title = ?, category = ?, youtube_url = ?, description = ? WHERE id = ?`
	_, err := db.Exec(query, t.Title, t.Category, t.YoutubeURL, t.Description, t.ID)
	return err
}

func DeleteThumbnail(db *sqlx.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM thumbnails WHERE id = ?`, id)
	return err
}

func DeleteThumbnailsByCategory(db *sqlx.DB, category string) error {
	_, err := db.Exec(`DELETE FROM thumbnails WHERE category = ?`, category)
	return err
}

func GetThumbnailCategories(db *sqlx.DB) ([]string, error) {
	var cats []string
	err := db.Select(&cats, `SELECT DISTINCT category FROM thumbnails ORDER BY category ASC`)
	return cats, err
}
