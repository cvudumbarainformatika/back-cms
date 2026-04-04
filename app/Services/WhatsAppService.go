package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cvudumbarainformatika/backend/config"
)

type WhatsAppService struct {
	config config.ZuwindaConfig
	client *http.Client
}

func NewWhatsAppService(cfg config.ZuwindaConfig) *WhatsAppService {
	return &WhatsAppService{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// LogEvent records a WhatsApp sending attempt to logs/whatsapp.log
func (s *WhatsAppService) LogEvent(to, status string) {
	logPath := "logs/whatsapp.log"
	
	// Ensure directory exists (basic check)
	_ = os.MkdirAll("logs", 0755)

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[WA-LOG-ERROR] Failed to open log file: %v\n", err)
		return
	}
	defer f.Close()

	// Format: 2026-04-04T12:10:30+00:00 status=sent to=62812...
	timestamp := time.Now().Format(time.RFC3339)
	logLine := fmt.Sprintf("%s status=%s to=%s\n", timestamp, status, to)

	if _, err := f.WriteString(logLine); err != nil {
		fmt.Printf("[WA-LOG-ERROR] Failed to write to log file: %v\n", err)
	}
}

// SendMessage sends a WhatsApp text message using Zuwinda Cloud API
func (s *WhatsAppService) SendMessage(to, message string) error {
	return s.send(to, "text", message, "")
}

// SendImageMessage sends a WhatsApp image message using Zuwinda Cloud API
func (s *WhatsAppService) SendImageMessage(to, message, urlFile string) error {
	return s.send(to, "image", message, urlFile)
}

// send is a generic helper to send different types of messages
func (s *WhatsAppService) send(to, messageType, content, urlFile string) error {
	// Zuwinda API Endpoint
	url := fmt.Sprintf("%s/messaging/whatsapp/message", s.config.BaseURL)

	// Prepare payload
	payload := map[string]interface{}{
		"accountId":   s.config.AccountID,
		"to":          to,
		"messageType": messageType,
		"content":     content,
	}

	if messageType == "image" && urlFile != "" {
		payload["url_file"] = urlFile
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Access-Key", s.config.AccessKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Zuwinda API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
