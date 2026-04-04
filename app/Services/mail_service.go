package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
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

	// Create new dialer
	d := gomail.NewDialer(s.Config.Host, s.Config.Port, dialerUser, dialerPass)

	// Bypass certificate verification for self-hosted mailserver (Docker internal network)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Use implicit SSL for port 465 (SMTPS)
	if s.Config.Port == 465 {
		d.SSL = true
	}

	// Updated Anti-Spam logic: Send individually with random delays and batch sleep
	total := len(to)
	batchSize := 40

	for i, recipient := range to {
		msgNum := i + 1
		m := gomail.NewMessage()
		m.SetHeader("From", fmt.Sprintf("Admin PDPI <%s>", from))
		m.SetHeader("To", recipient)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)

		log.Printf("[MAIL] [%d/%d] Sending to: %s", msgNum, total, recipient)

		if err := d.DialAndSend(m); err != nil {
			log.Printf("[MAIL ERROR] Failed to send to %s: %v", recipient, err)
			// Continue to next recipient despite error
		}

		// 1. Random delay 3-6 seconds after EVERY message
		if msgNum < total {
			delay := 3 + rand.Intn(4) // 3, 4, 5, or 6 seconds
			
			// 2. Batch system: After 40 messages, sleep 3-6 minutes
			if msgNum % batchSize == 0 {
				sleepMinutes := 3 + rand.Intn(4) // 3, 4, 5, or 6 minutes
				totalBatches := (total + batchSize - 1) / batchSize
				currentBatch := msgNum / batchSize
				
				log.Printf("[MAIL] Batch %d/%d selesai, istirahat %d menit dulu untuh cegah spam...", 
					currentBatch, totalBatches, sleepMinutes)
				
				time.Sleep(time.Duration(sleepMinutes) * time.Minute)
			} else {
				// Regular message delay
				time.Sleep(time.Duration(delay) * time.Second)
			}
		}
	}

	log.Printf("[MAIL] Broadcast selesai dikirim ke %d penerima.", total)
	return nil
}
