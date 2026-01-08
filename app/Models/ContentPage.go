package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ContentAuthor struct {
	Name   string `json:"name"`
	To     string `json:"to,omitempty"`
	Avatar struct {
		Src string `json:"src"`
	} `json:"avatar,omitempty"`
}

// ContentAuthors handles JSON marshaling for authors
type ContentAuthors []ContentAuthor

func (a ContentAuthors) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *ContentAuthors) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

type ContentPage struct {
	ID          int64          `db:"id" json:"id"`
	Slug        string         `db:"slug" json:"slug"`
	Title       string         `db:"title" json:"title"`
	Description string         `db:"description" json:"description"`
	Body        string         `db:"body" json:"body"` // Markdown content
	HTML        string         `db:"html" json:"html"` // HTML content (from WYSIWYG)
	Date        time.Time      `db:"date" json:"date"`
	ImageSrc    string         `db:"image_src" json:"image_src"`
	BadgeLabel  string         `db:"badge_label" json:"badge_label"`
	Authors     ContentAuthors `db:"authors" json:"authors"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}
