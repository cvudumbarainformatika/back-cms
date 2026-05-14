package services

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate        string `xml:"pubDate"`
	ContentEncoded string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	Enclosure      struct {
		URL string `xml:"url,attr"`
	} `xml:"enclosure"`
}

type RSSChannel struct {
	Items []RSSItem `xml:"item"`
}

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

type RSSService struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func NewRSSService(db *sqlx.DB, rdb *redis.Client) *RSSService {
	return &RSSService{DB: db, Redis: rdb}
}

func (s *RSSService) FetchAndStoreFeeds() error {
	sources := []struct {
		Name string
		URL  string
	}{
		{"CNN Indonesia Health", "https://www.cnnindonesia.com/gaya-hidup/rss"},
		{"Antara News Health", "https://www.antaranews.com/rss/terkini.xml"},
		{"Detik Health", "https://health.detik.com/rss"},
	}

	keywords := []string{"kesehatan", "rumah sakit", "dokter", "medis", "pasien", "penyakit", "obat"}

	for _, source := range sources {
		log.Printf("[RSS] Fetching from %s...", source.Name)
		items, err := s.fetchRSS(source.URL)
		if err != nil {
			log.Printf("[RSS ERROR] Failed to fetch from %s: %v", source.Name, err)
			continue
		}

		for _, item := range items {
			// Filter by keywords
			match := false
			fullText := strings.ToLower(item.Title + " " + item.Description)
			for _, kw := range keywords {
				if strings.Contains(fullText, kw) {
					match = true
					break
				}
			}

			if !match {
				continue
			}

			// Parse date (try multiple formats)
			dateFormats := []string{time.RFC1123Z, time.RFC1123, "Mon, 02 Jan 2006 15:04:05 MST", "2006-01-02 15:04:05"}
			var pubDate time.Time
			for _, fmtStr := range dateFormats {
				if t, err := time.Parse(fmtStr, item.PubDate); err == nil {
					pubDate = t
					break
				}
			}

			thumbURL := item.Enclosure.URL
			if thumbURL == "" {
				re := regexp.MustCompile(`<img[^>]+src="([^">]+)"`)
				match := re.FindStringSubmatch(item.Description)
				if len(match) > 1 {
					thumbURL = match[1]
				}
			}

			// Use ContentEncoded if available and longer than Description
			bestDescription := item.Description
			if item.ContentEncoded != "" && len(item.ContentEncoded) > len(item.Description) {
				bestDescription = item.ContentEncoded
			}

			extNews := &models.ExternalNews{
				Title:        item.Title,
				Source:       source.Name,
				URL:          item.Link,
				Description:  sql.NullString{String: bestDescription, Valid: bestDescription != ""},
				ThumbnailURL: sql.NullString{String: thumbURL, Valid: thumbURL != ""},
				PublishedAt:  &pubDate,
				Status:       "pending",
			}

			_, err := models.CreateExternalNews(context.Background(), s.DB, extNews)
			if err != nil {
				// Ignore duplicate key errors (URL is unique)
				if !strings.Contains(err.Error(), "Duplicate entry") {
					log.Printf("[RSS ERROR] Failed to save item %s: %v", item.Title, err)
				}
			}
		}
	}

	return nil
}

func (s *RSSService) fetchRSS(url string) ([]RSSItem, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set User-Agent to avoid being blocked by media sites
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}

	return feed.Channel.Items, nil
}
