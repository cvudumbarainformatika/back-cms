package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
