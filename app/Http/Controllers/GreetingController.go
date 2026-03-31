package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	requests "github.com/cvudumbarainformatika/backend/app/Http/Requests"
	models "github.com/cvudumbarainformatika/backend/app/Models"
	services "github.com/cvudumbarainformatika/backend/app/Services"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// GreetingController handles greeting operations
type GreetingController struct {
	db     *sqlx.DB
	mail   *services.MailService
	wa     *services.WhatsAppService
	config *config.Config
}

// NewGreetingController creates a new GreetingController instance
func NewGreetingController(db *sqlx.DB, mail *services.MailService, wa *services.WhatsAppService, cfg *config.Config) *GreetingController {
	return &GreetingController{
		db:     db,
		mail:   mail,
		wa:     wa,
		config: cfg,
	}
}

// GetList returns paginated list of greetings
func (gc *GreetingController) GetList(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	search := c.Query("search")

	offset := (page - 1) * limit
	greetings, total, err := models.GetAllGreetings(gc.db, offset, limit, search)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch greetings", nil)
		return
	}

	pagination := utils.OffsetPaginate(greetings, page, limit, total)
	utils.Success(c, http.StatusOK, "Greetings fetched successfully", gin.H{
		"items":      pagination.Data,
		"pagination": pagination.Meta,
	})
}

// GetByID returns a single greeting by ID
func (gc *GreetingController) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid greeting ID", nil)
		return
	}

	greeting, err := models.FindGreetingByID(gc.db, id)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to fetch greeting", nil)
		return
	}

	if greeting == nil {
		utils.Error(c, http.StatusNotFound, "greeting_not_found", "Greeting not found", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Greeting retrieved successfully", greeting)
}

// Create creates a new greeting
func (gc *GreetingController) Create(c *gin.Context) {
	var req requests.CreateGreetingRequest
	if err := req.Validate(c); err != nil {
		return
	}

	greeting := &models.Greeting{
		Title:    req.Title,
		Content:  req.Content,
		ImageURL: req.ImageURL,
		IsActive: req.IsActive,
	}

	if err := greeting.Create(gc.db); err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to create greeting: "+err.Error(), nil)
		return
	}

	utils.Success(c, http.StatusCreated, "Greeting created successfully", greeting)
}

// Update updates a greeting
func (gc *GreetingController) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid greeting ID", nil)
		return
	}

	var req requests.UpdateGreetingRequest
	if err := req.Validate(c); err != nil {
		return
	}

	greeting, err := models.FindGreetingByID(gc.db, id)
	if err != nil || greeting == nil {
		utils.Error(c, http.StatusNotFound, "greeting_not_found", "Greeting not found", nil)
		return
	}

	greeting.Title = req.Title
	greeting.Content = req.Content
	greeting.ImageURL = req.ImageURL
	greeting.IsActive = req.IsActive

	if err := greeting.Update(gc.db); err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to update greeting", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Greeting updated successfully", greeting)
}

// Delete soft deletes a greeting
func (gc *GreetingController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid greeting ID", nil)
		return
	}

	greeting, err := models.FindGreetingByID(gc.db, id)
	if err != nil || greeting == nil {
		utils.Error(c, http.StatusNotFound, "greeting_not_found", "Greeting not found", nil)
		return
	}

	if err := greeting.Delete(gc.db); err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to delete greeting", nil)
		return
	}

	utils.Success(c, http.StatusOK, "Greeting deleted successfully", nil)
}

// SendWA sends greeting via WhatsApp
func (gc *GreetingController) SendWA(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid greeting ID", nil)
		return
	}

	greeting, err := models.FindGreetingByID(gc.db, id)
	if err != nil || greeting == nil {
		utils.Error(c, http.StatusNotFound, "greeting_not_found", "Greeting not found", nil)
		return
	}

	target := c.DefaultQuery("target", "test")
	recipients, err := gc.getWARecipients(target)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "error", "Failed to fetch recipients", nil)
		return
	}

	// Prepare message
	message := fmt.Sprintf("*%s*\n\n%s", greeting.Title, greeting.Content)
	
	// Send in background
	go func(to []string, msg, img string) {
		// Use placeholder image if in local environment (Zuwinda needs public URL)
		if gc.config.App.Env == "local" && img != "" {
			img = "https://hobiternak.com/wp-content/uploads/2017/08/bebek-unggulan.jpg"
		} else if img != "" && strings.HasPrefix(img, "/") {
			// Handle relative URL in production
			baseURL := strings.TrimSuffix(gc.config.App.BaseURL, "/")
			img = baseURL + img
		}
		log.Printf("[WA] Preparing to send to recipients. Image URL: %s", img)

		for _, number := range to {
			var err error
			if img != "" {
				err = gc.wa.SendImageMessage(number, msg, img)
				log.Printf("[WA] Sending image to %s: err=%v", number, err)
			} else {
				err = gc.wa.SendMessage(number, msg)
				log.Printf("[WA] Sending text to %s: err=%v", number, err)
			}
			// Jeda random antara 3 sampai 6 detik untuk menghindari deteksi spam
			jitter := rand.Intn(3000) // 0-3000ms
			time.Sleep(time.Duration(3000+jitter) * time.Millisecond)
		}
	}(recipients, message, greeting.ImageURL)

	utils.Success(c, http.StatusOK, "WhatsApp sending started in background", gin.H{
		"recipient_count": len(recipients),
	})
}

// SendEmail sends greeting via Email
func (gc *GreetingController) SendEmail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid greeting ID", nil)
		return
	}

	greeting, err := models.FindGreetingByID(gc.db, id)
	if err != nil || greeting == nil {
		utils.Error(c, http.StatusNotFound, "greeting_not_found", "Greeting not found", nil)
		return
	}

	target := c.DefaultQuery("target", "test")
	recipients, err := gc.getEmailRecipients(target)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "error", "Failed to fetch recipients", nil)
		return
	}

	subject := greeting.Title
	body := gc.getGreetingEmailTemplate(greeting.Title, greeting.ImageURL, greeting.Content)

	// Send in background
	go func(to []string, subj, content string) {
		if err := gc.mail.SendEmail(to, subj, content); err != nil {
			log.Printf("Failed to send greeting email: %v", err)
		}
	}(recipients, subject, body)

	utils.Success(c, http.StatusOK, "Email sending started in background", gin.H{
		"recipient_count": len(recipients),
	})
}

// Helper methods (extracted from BroadcastController logic)

func (gc *GreetingController) getWARecipients(target string) ([]string, error) {
	if target == "all" {
		var phones []string
		query := "SELECT DISTINCT no_hp FROM pdpi_members WHERE no_hp IS NOT NULL AND no_hp != ''"
		err := gc.db.Select(&phones, query)
		if err != nil {
			return nil, err
		}

		var normalized []string
		for _, p := range phones {
			norm := utils.NormalizePhoneNumber(p)
			if norm != "" {
				normalized = append(normalized, norm)
			}
		}
		return normalized, nil
	}

	// Hardcoded test numbers from BroadcastController
	return []string{
		"6281237660656", 
		"6282334148314", 
		"6285736336536", 
		"6282324141494",
		"6281350125649",
		"6281381295959",
		"6281554545012",
		"628155088994",
	}, nil
}

func (gc *GreetingController) getEmailRecipients(target string) ([]string, error) {
	if target == "all" {
		var emails []string
		query := "SELECT DISTINCT email FROM pdpi_members WHERE email IS NOT NULL AND email != ''"
		err := gc.db.Select(&emails, query)
		if err != nil {
			return nil, err
		}
		return emails, nil
	}

	// Hardcoded test emails from BroadcastController
	return []string{"pharyyady@gmail.com", "cvudumbarainformatika@gmail.com", "vichoirul@gmail.com"}, nil
}

func (gc *GreetingController) getGreetingEmailTemplate(title, imageURL, content string) string {
	// Simple elegant template for greetings
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<style>
	body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #f8fafc; margin: 0; padding: 0; }
	.container { max-width: 600px; margin: 40px auto; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 6px rgba(0,0,0,0.05); border: 1px solid #e2e8f0; }
	.header { background: linear-gradient(135deg, #1e40af 0%%, #3b82f6 100%%); padding: 40px 20px; text-align: center; color: #ffffff; }
	.header h1 { margin: 0; font-size: 28px; font-weight: 700; letter-spacing: -0.025em; }
	.hero-image { width: 100%%; height: auto; display: block; max-height: 400px; object-fit: cover; }
	.content { padding: 40px; color: #334155; line-height: 1.7; font-size: 16px; }
	.content p { margin-top: 0; white-space: pre-line; }
	.footer { background-color: #f1f5f9; padding: 24px; text-align: center; font-size: 13px; color: #64748b; border-top: 1px solid #e2e8f0; }
	.logo { font-weight: bold; color: #1e40af; margin-bottom: 8px; font-size: 18px; }
</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>%s</h1>
		</div>
		%s
		<div class="content">
			<p>%s</p>
		</div>
		<div class="footer">
			<div class="logo">PDPI</div>
			<p>&copy; 2026 Perhimpunan Dokter Paru Indonesia. All rights reserved.</p>
			<p>Pesan ini dikirim secara otomatis melalui sistem dashboard anggota.</p>
		</div>
	</div>
</body>
</html>
	`, title, func() string {
		if imageURL != "" {
			return fmt.Sprintf(`<img src="%s" alt="%s" class="hero-image"/>`, imageURL, title)
		}
		return ""
	}(), content)
}
