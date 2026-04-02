package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/cvudumbarainformatika/backend/config"

	"gopkg.in/gomail.v2"
)

type MailService struct {
	Config config.MailConfig
}

func NewMailService(cfg config.MailConfig) *MailService {
	return &MailService{
		Config: cfg,
	}
}

func (s *MailService) SendEmail(to []string, subject string, body string) error {
	// Dry Run / Log Mode Check
	// Only log if Host is not configured
	if s.Config.Host == "" {
		log.Printf("[MAIL LOG MODE] Would send email to: %v\nSubject: %s\nBody Length: %d bytes", to, subject, len(body))
		return nil
	}

	from := s.Config.User
	if from == "" {
		from = "noreply@localhost"
	}

	dialerUser := s.Config.User
	dialerPass := s.Config.Password

	// If using Port 25 (Internal Relay), disable Auth to avoid STARTTLS requirement
	if s.Config.Port == 25 {
		dialerUser = ""
		dialerPass = ""
	}

	d := gomail.NewDialer(s.Config.Host, s.Config.Port, dialerUser, dialerPass)

	// Bypass certificate verification for self-hosted mailserver (Docker internal network)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Use implicit SSL for port 465 (SMTPS)
	if s.Config.Port == 465 {
		d.SSL = true
	}

	// Batch recipients to avoid SMTP "too many recipients" error (452 4.5.3)
	// Common limits are 100-500, but 50 is a very safe baseline.
	batchSize := 50
	for i := 0; i < len(to); i += batchSize {
		end := i + batchSize
		if end > len(to) {
			end = len(to)
		}
		batch := to[i:end]

		m := gomail.NewMessage()
		// Set sender with display name
		m.SetHeader("From", fmt.Sprintf("Admin PDPI <%s>", from))
		// Use BCC for privacy - recipients won't see each other's emails
		m.SetHeader("To", from)      // Set To to sender (admin)
		m.SetHeader("Bcc", batch...) // This batch of recipients in BCC (hidden)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)

		log.Printf("[MAIL] Sending batch %d/%d (%d recipients)", (i/batchSize)+1, (len(to)+batchSize-1)/batchSize, len(batch))

		if err := d.DialAndSend(m); err != nil {
			return fmt.Errorf("failed to send email batch starting at index %d: %v", i, err)
		}

		// Small delay between batches to avoid server throttling/rate limiting
		if end < len(to) {
			time.Sleep(800 * time.Millisecond)
		}
	}

	return nil
}
