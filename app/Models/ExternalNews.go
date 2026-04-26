package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type ExternalNews struct {
	ID           int            `db:"id" json:"id"`
	Title        string         `db:"title" json:"title"`
	Source       string         `db:"source" json:"source"`
	URL          string         `db:"url" json:"url"`
	Description  sql.NullString `db:"description" json:"description"`
	ThumbnailURL sql.NullString `db:"thumbnail_url" json:"thumbnail_url"`
	PublishedAt  *time.Time     `db:"published_at" json:"published_at"`
	FullContent  sql.NullString `db:"full_content" json:"full_content"`
	Status       string         `db:"status" json:"status"`
	ScrapedAt    *time.Time     `db:"scraped_at" json:"scraped_at"`
	IsImported   bool           `db:"is_imported" json:"is_imported"`
	ImportedSlug sql.NullString `db:"imported_slug" json:"imported_slug"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`
}

func GetAllExternalNews(ctx context.Context, db *sqlx.DB, limit int, offset int, isImported string) ([]ExternalNews, int, error) {
	var items []ExternalNews
	var total int

	whereClause := ""
	args := []interface{}{}

	if isImported == "1" || isImported == "true" {
		whereClause = "WHERE en.is_imported = 1"
	} else if isImported == "0" || isImported == "false" {
		whereClause = "WHERE en.is_imported = 0"
	}

	orderBy := "en.published_at DESC"
	if isImported == "1" || isImported == "true" {
		orderBy = "en.updated_at DESC"
	}

	query := fmt.Sprintf(`
		SELECT en.*, b.slug as imported_slug 
		FROM external_news en 
		LEFT JOIN berita b ON b.title COLLATE utf8mb4_unicode_ci = en.title COLLATE utf8mb4_unicode_ci AND b.deleted_at IS NULL
		%s
		ORDER BY %s LIMIT ? OFFSET ?
	`, whereClause, orderBy)
	
	args = append(args, limit, offset)

	err := db.SelectContext(ctx, &items, query, args...)
	if err != nil {
		return nil, 0, err
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM external_news en %s", whereClause)
	err = db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func CreateExternalNews(ctx context.Context, db *sqlx.DB, item *ExternalNews) (int64, error) {
	query := `INSERT INTO external_news (title, source, url, description, thumbnail_url, published_at, is_imported, status) 
			  VALUES (:title, :source, :url, :description, :thumbnail_url, :published_at, :is_imported, :status)
			  ON DUPLICATE KEY UPDATE 
			  title = VALUES(title), 
			  description = VALUES(description), 
			  thumbnail_url = VALUES(thumbnail_url), 
			  published_at = VALUES(published_at)`
	
	res, err := db.NamedExecContext(ctx, query, item)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil || lastID == 0 {
		// If duplicate key update didn't insert, we might need to get the ID by URL
		var id int64
		err = db.GetContext(ctx, &id, "SELECT id FROM external_news WHERE url = ?", item.URL)
		return id, err
	}

	return lastID, nil
}

func UpdateExternalNewsContent(ctx context.Context, db *sqlx.DB, id int, content string) error {
	query := `UPDATE external_news SET full_content = ?, status = 'scraped', scraped_at = NOW() WHERE id = ?`
	_, err := db.ExecContext(ctx, query, content, id)
	return err
}

func GetPendingExternalNews(ctx context.Context, db *sqlx.DB, limit int) ([]ExternalNews, error) {
	var items []ExternalNews
	query := `SELECT * FROM external_news WHERE status = 'pending' ORDER BY published_at DESC LIMIT ?`
	err := db.SelectContext(ctx, &items, query, limit)
	return items, err
}

func MarkAsImported(ctx context.Context, db *sqlx.DB, id int) error {
	_, err := db.ExecContext(ctx, "UPDATE external_news SET is_imported = 1 WHERE id = ?", id)
	return err
}
