package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ScraperService struct{}

func NewScraperService() *ScraperService {
	return &ScraperService{}
}

func (s *ScraperService) ScrapeDetail(url string) (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var content string
	var selector string

	// Determine selector based on URL
	if strings.Contains(url, "cnnindonesia.com") {
		selector = ".detail-text, .detail__body-text"
	} else if strings.Contains(url, "antaranews.com") {
		selector = ".post-content, .article-content"
	} else if strings.Contains(url, "detik.com") {
		selector = ".detail__body-text, .itp_bodycontent"
	} else {
		// Fallback to common content selectors
		selector = "article, .content, .entry-content"
	}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		// Remove unwanted elements
		s.Find("script, style, iframe, .parallax, .sisipan, .ad-box, .ads").Remove()
		
		html, err := s.Html()
		if err == nil {
			content += html
		}
	})

	if content == "" {
		return "", fmt.Errorf("could not find content with selector %s", selector)
	}

	return strings.TrimSpace(content), nil
}
