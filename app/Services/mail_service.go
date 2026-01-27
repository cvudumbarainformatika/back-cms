package services

import (
	"fmt"
	"log"

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

	m := gomail.NewMessage()
	from := s.Config.User
	if from == "" {
		from = "noreply@localhost"
	}
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.Config.Host, s.Config.Port, s.Config.User, s.Config.Password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
