package controllers

import (
	"fmt"
	"net/http"

	services "github.com/cvudumbarainformatika/backend/app/Services"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type BroadcastController struct {
	MailService *services.MailService
	DB          *sqlx.DB
	Config      config.AppConfig
}

func NewBroadcastController(mailService *services.MailService, db *sqlx.DB, appConfig config.AppConfig) *BroadcastController {
	return &BroadcastController{
		MailService: mailService,
		DB:          db,
		Config:      appConfig,
	}
}

// Hardcoded recipients for initial phase
var recipients = []string{
	"pharyyady@gmail.com",
	"cvudumbarainformatika@gmail.com",
	"vichoirul@gmail.com",
	"ariezlnd69@gmail.com",
	"andi.meka.025@gmail.com",
	"vpluzt@gmail.com",
	"fafnir573@gmail.com",
}

// getRecipients determines who to email based on query param
// ?target=all -> Fetch from DB
// Default -> Use hardcoded test list
func (ctrl *BroadcastController) getRecipients(c *gin.Context) ([]string, error) {
	target := c.Query("target")
	if target == "all" {
		var emails []string
		query := "SELECT DISTINCT email FROM pdpi_members WHERE email IS NOT NULL AND email != ''"
		err := ctrl.DB.Select(&emails, query)
		if err != nil {
			return nil, err
		}
		return emails, nil
	}
	return recipients, nil
}

func (ctrl *BroadcastController) BroadcastBerita(c *gin.Context) {
	id := c.Param("id")

	// 1. Fetch Berita
	var berita struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		Status   string `json:"status"`
		Slug     string `json:"slug"`
		ImageURL string `json:"image_url" db:"image_url"`
		Excerpt  string `json:"excerpt"`
	}

	query := "SELECT id, title, content, status, slug, image_url, excerpt FROM berita WHERE id = ?"
	err := ctrl.DB.QueryRow(query, id).Scan(&berita.ID, &berita.Title, &berita.Content, &berita.Status, &berita.Slug, &berita.ImageURL, &berita.Excerpt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita not found"})
		return
	}

	// 2. Validate Status
	if berita.Status != "published" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only published news can be broadcasted"})
		return
	}

	// 3. Determine Recipients
	targetRecipients, err := ctrl.getRecipients(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipients: " + err.Error()})
		return
	}

	if len(targetRecipients) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No recipients found"})
		return
	}

	// 4. Prepare Email Content
	baseURL := ctrl.Config.BaseURL
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	subject := fmt.Sprintf("[PDPI News] %s", berita.Title)
	readMoreLink := fmt.Sprintf("%s/berita/%s", baseURL, berita.Slug)

	imageURL := berita.ImageURL
	if imageURL == "" {
		imageURL = fmt.Sprintf("%s/images/default-news.jpg", baseURL)
	} else if len(imageURL) > 0 && imageURL[0] == '/' {
		imageURL = fmt.Sprintf("%s%s", baseURL, imageURL)
	}

	body := getEmailTemplate(berita.Title, imageURL, berita.Excerpt, readMoreLink, "Baca Berita Selengkapnya")

	// 5. Send Email via Background Goroutine (Fire and Forget)
	// Prevents timeout when sending to many recipients
	go func(targets []string, subj, content string) {
		fmt.Printf("Starting broadcast for %d recipients...\n", len(targets))
		if err := ctrl.MailService.SendEmail(targets, subj, content); err != nil {
			fmt.Printf("Error sending broadcast: %v\n", err)
		} else {
			fmt.Printf("Broadcast completed for %d recipients.\n", len(targets))
		}
	}(targetRecipients, subject, body)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Broadcast process started in background",
		"recipient_count": len(targetRecipients),
		"mode":            c.Query("target"),
	})
}

func (ctrl *BroadcastController) BroadcastAgenda(c *gin.Context) {
	id := c.Param("id")

	// 1. Fetch Agenda
	var agenda struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Date        string `json:"date"`
		Slug        string `json:"slug"`
		ImageURL    string `json:"image_url" db:"image_url"`
	}

	query := "SELECT id, title, description, status, date, slug, image_url FROM agenda WHERE id = ?"
	err := ctrl.DB.QueryRow(query, id).Scan(&agenda.ID, &agenda.Title, &agenda.Description, &agenda.Status, &agenda.Date, &agenda.Slug, &agenda.ImageURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agenda not found"})
		return
	}

	// 2. Validate Status
	if agenda.Status != "published" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only published agenda can be broadcasted"})
		return
	}

	// 3. Determine Recipients
	targetRecipients, err := ctrl.getRecipients(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipients: " + err.Error()})
		return
	}

	if len(targetRecipients) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No recipients found"})
		return
	}

	// 4. Prepare Email Content
	baseURL := ctrl.Config.BaseURL
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	subject := fmt.Sprintf("[PDPI Agenda] %s", agenda.Title)
	readMoreLink := fmt.Sprintf("%s/agenda/%s", baseURL, agenda.Slug)

	imageURL := agenda.ImageURL
	if imageURL == "" {
		imageURL = fmt.Sprintf("%s/images/default-agenda.jpg", baseURL)
	} else if len(imageURL) > 0 && imageURL[0] == '/' {
		imageURL = fmt.Sprintf("%s%s", baseURL, imageURL)
	}

	content := fmt.Sprintf("<p><strong>Tanggal:</strong> %s</p><br/>%s", agenda.Date, agenda.Description)
	body := getEmailTemplate(agenda.Title, imageURL, content, readMoreLink, "Lihat Detail Agenda")

	// 5. Send Email via Background Goroutine
	go func(targets []string, subj, content string) {
		fmt.Printf("Starting broadcast for %d recipients...\n", len(targets))
		if err := ctrl.MailService.SendEmail(targets, subj, content); err != nil {
			fmt.Printf("Error sending broadcast: %v\n", err)
		} else {
			fmt.Printf("Broadcast completed for %d recipients.\n", len(targets))
		}
	}(targetRecipients, subject, body)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Broadcast process started in background",
		"recipient_count": len(targetRecipients),
		"mode":            c.Query("target"),
	})
}

func getEmailTemplate(title, imageURL, description, link, buttonText string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<style>
	body { font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
	.container { max-width: 600px; margin: 20px auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
	.header { background-color: #2563eb; padding: 20px; text-align: center; }
	.header h1 { color: #ffffff; margin: 0; font-size: 24px; }
	.hero-image { width: 100%%; height: auto; display: block; max-height: 300px; object-fit: cover; }
	.content { padding: 30px; color: #333333; line-height: 1.6; }
	.content h2 { margin-top: 0; color: #1f2937; }
	.button-container { text-align: center; margin-top: 30px; margin-bottom: 20px; }
	.button { background-color: #2563eb; color: #ffffff; padding: 12px 24px; text-decoration: none; border-radius: 6px; font-weight: bold; display: inline-block; }
	.footer { background-color: #f9fafb; padding: 20px; text-align: center; font-size: 12px; color: #6b7280; border-top: 1px solid #e5e7eb; }
</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>PDPI Update</h1>
		</div>
		<img src="%s" alt="%s" class="hero-image" onerror="this.style.display='none'"/>
		<div class="content">
			<h2>%s</h2>
			<div>%s</div>
			<div class="button-container">
				<a href="%s" class="button">%s</a>
			</div>
		</div>
		<div class="footer">
			<p>&copy; %s Perhimpunan Dokter Paru Indonesia. All rights reserved.</p>
			<p>Anda menerima email ini karena Anda terdaftar sebagai anggota.</p>
		</div>
	</div>
</body>
</html>
	`, imageURL, title, title, description, link, buttonText, "2026") // Start year hardcoded or dynamic
}
