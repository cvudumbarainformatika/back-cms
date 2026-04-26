package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type NewsWorker struct {
	DB             *sqlx.DB
	Redis          *redis.Client
	ScraperService *ScraperService
}

func NewNewsWorker(db *sqlx.DB, rdb *redis.Client) *NewsWorker {
	return &NewsWorker{
		DB:             db,
		Redis:          rdb,
		ScraperService: NewScraperService(),
	}
}

func (w *NewsWorker) Start(ctx context.Context) {
	log.Println("[NewsWorker] Started background worker for detail scraping")
	
	// Start a worker loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("[NewsWorker] Stopping worker...")
				return
			default:
				// BLPOP with timeout to allow graceful shutdown
				result, err := w.Redis.BLPop(ctx, 5*time.Second, "news:scrape_queue").Result()
				if err != nil {
					if err == redis.Nil {
						// Timeout, continue
						continue
					}
					log.Printf("[NewsWorker ERROR] BLPop failed: %v", err)
					time.Sleep(2 * time.Second)
					continue
				}

				if len(result) < 2 {
					continue
				}

				newsID := result[1]
				log.Printf("[NewsWorker] Processing news item ID: %s", newsID)

				// Tambahkan random delay 5-10 detik agar tidak diblokir sumber
				delay := rand.Intn(6) + 5 // 5 to 10 seconds
				log.Printf("[NewsWorker] Waiting for %d seconds before processing...", delay)
				time.Sleep(time.Duration(delay) * time.Second)

				w.processItem(newsID)
			}
		}
	}()
}

func (w *NewsWorker) processItem(newsIDStr string) {
	var id int
	_, err := fmt.Sscanf(newsIDStr, "%d", &id)
	if err != nil {
		log.Printf("[NewsWorker ERROR] Invalid news ID %s: %v", newsIDStr, err)
		return
	}

	// 1. Fetch from DB
	var news models.ExternalNews
	err = w.DB.Get(&news, "SELECT * FROM external_news WHERE id = ?", id)
	if err != nil {
		log.Printf("[NewsWorker ERROR] News ID %d not found in DB: %v", id, err)
		return
	}

	// 2. Scrape content
	log.Printf("[NewsWorker] Scraping detail for: %s", news.URL)
	content, err := w.ScraperService.ScrapeDetail(news.URL)
	if err != nil {
		log.Printf("[NewsWorker ERROR] Failed to scrape %s: %v", news.URL, err)
		
		// Update status to error
		w.DB.Exec("UPDATE external_news SET status = 'error' WHERE id = ?", id)
		return
	}

	// 3. Update DB
	err = models.UpdateExternalNewsContent(context.Background(), w.DB, id, content)
	if err != nil {
		log.Printf("[NewsWorker ERROR] Failed to update content for ID %d: %v", id, err)
		return
	}

	log.Printf("[NewsWorker SUCCESS] Scraped and updated content for ID %d", id)
}
