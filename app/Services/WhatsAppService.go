package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cvudumbarainformatika/backend/config"
)

type WhatsAppService struct {
	config config.WABA360Config
	client *http.Client
}

func NewWhatsAppService(cfg config.WABA360Config) *WhatsAppService {
	return &WhatsAppService{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// LogEvent records a WhatsApp sending attempt
func (s *WhatsAppService) LogEvent(to, status string) {
	logPath := "logs/whatsapp.log"
	_ = os.MkdirAll("logs", 0755)

	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[WA-LOG-ERROR] Failed to open log file: %v\n", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format(time.RFC3339)
	logLine := fmt.Sprintf("%s status=%s to=%s\n", timestamp, status, to)
	_, _ = f.WriteString(logLine)
}

// WhatsApp Template Structures
type WABAParameter struct {
	Type          string      `json:"type"`
	ParameterName string      `json:"parameter_name,omitempty"`
	Text          string      `json:"text,omitempty"`
	Image         *WABAImage `json:"image,omitempty"`
}

type WABAImage struct {
	Link string `json:"link"`
}

type WABAComponent struct {
	Type       string          `json:"type"`
	Parameters []WABAParameter `json:"parameters"`
}

type WABATemplate struct {
	Name       string          `json:"name"`
	Language   map[string]string `json:"language"`
	Components []WABAComponent `json:"components"`
}

type WABAPayload struct {
	MessagingProduct string       `json:"messaging_product"`
	To               string       `json:"to"`
	Type             string       `json:"type"`
	Template         WABATemplate `json:"template"`
}

// SendTemplateMessage sends a template message using 360dialog WABA API
func (s *WhatsAppService) SendTemplateMessage(to, templateName string, bodyParams []WABAParameter, imageURL string) error {
	components := []WABAComponent{}

	// Clean imageURL from HTML escaped characters like &amp;
	imageURL = html.UnescapeString(imageURL)

	// 1. Add Header if image is provided
	if imageURL != "" {
		components = append(components, WABAComponent{
			Type: "header",
			Parameters: []WABAParameter{
				{
					Type: "image",
					Image: &WABAImage{Link: imageURL},
				},
			},
		})
	}

	// 2. Add Body Parameters
	components = append(components, WABAComponent{
		Type:       "body",
		Parameters: bodyParams,
	})

	payload := WABAPayload{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: WABATemplate{
			Name:     templateName,
			Language: map[string]string{"code": "id"},
			Components: components,
		},
	}

	return s.execute(payload)
}

// SendArtikel sends article/news notification
func (s *WhatsAppService) SendArtikel(to, title, url, imageURL string) error {
	// Using positional parameters for artikel_inter as per user snippet
	params := []WABAParameter{
		{Type: "text", Text: title},
		{Type: "text", Text: url},
	}
	return s.SendTemplateMessage(to, "artikel_inter", params, imageURL)
}

// SendAgenda sends agenda/event notification
func (s *WhatsAppService) SendAgenda(to, title, mode, location, timeStr, fee, quota, url, imageURL string) error {
	// Using positional parameters for agenda_gmbr as per user snippet
	params := []WABAParameter{
		{Type: "text", Text: title},
		{Type: "text", Text: mode},
		{Type: "text", Text: location},
		{Type: "text", Text: timeStr},
		{Type: "text", Text: fee},
		{Type: "text", Text: quota},
		{Type: "text", Text: url},
	}
	return s.SendTemplateMessage(to, "agenda_gmbr", params, imageURL)
}

// SendGreeting sends birthday/greeting message
func (s *WhatsAppService) SendGreeting(to, title, url, imageURL string) error {
	// Greetings uses same template as artikel as per user request
	return s.SendArtikel(to, title, url, imageURL)
}

// SendUltah sends birthday greeting using ultah template
func (s *WhatsAppService) SendUltah(to, gelar, nama, gelar2 string) error {
	// Remove trailing "." from gelar
	gelar = strings.TrimSuffix(gelar, ".")

	// Using positional parameters for ultah template
	params := []WABAParameter{
		{Type: "text", Text: gelar},
		{Type: "text", Text: nama},
		{Type: "text", Text: gelar2},
	}
	return s.SendTemplateMessage(to, "ultah2", params, "")
}

func (s *WhatsAppService) execute(payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.config.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("D360-API-KEY", s.config.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("[360dialog Debug] Payload: %s\n", string(jsonData))
	fmt.Printf("[360dialog Debug] Response (%d): %s\n", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("360dialog API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// Legacy compatibility (Optional: redirect to templates or keep empty)
func (s *WhatsAppService) SendMessage(to, message string) error {
	// For 360dialog, we mostly use templates. Mapping to a generic template if needed.
	return fmt.Errorf("direct SendMessage not supported for 360dialog template-only setup")
}

func (s *WhatsAppService) SendImageMessage(to, message, urlFile string) error {
	return fmt.Errorf("direct SendImageMessage not supported for 360dialog template-only setup")
}
