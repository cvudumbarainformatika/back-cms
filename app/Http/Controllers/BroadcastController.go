package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cvudumbarainformatika/backend/utils"
	services "github.com/cvudumbarainformatika/backend/app/Services"
	"github.com/cvudumbarainformatika/backend/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type BroadcastController struct {
	MailService     *services.MailService
	WAService       *services.WhatsAppService
	BirthdayService *services.BirthdayService
	DB              *sqlx.DB
	Config          config.AppConfig
}

func NewBroadcastController(mailService *services.MailService, waService *services.WhatsAppService, birthdayService *services.BirthdayService, db *sqlx.DB, appConfig config.AppConfig) *BroadcastController {
	return &BroadcastController{
		MailService:     mailService,
		WAService:       waService,
		BirthdayService: birthdayService,
		DB:              db,
		Config:          appConfig,
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

// WhatsApp test recipients
var waTestRecipients = []string{
	"6281237660656",
	"6282334148314",
	"6285736336536",
	"6282324141494",
}

// getRecipients determines who to email based on query param
// ?target=all -> Fetch ALL members from DB (1800+)
// ?target=warmup -> Fetch 50 members + 7 test emails (IP Warm-up)
// Default -> Use hardcoded test list (7 emails)
func (ctrl *BroadcastController) getRecipients(c *gin.Context) ([]string, error) {
	target := c.Query("target")

	if target == "all" {
		var emails []string
		query := "SELECT DISTINCT email FROM pdpi_members WHERE email IS NOT NULL AND email != ''"
		err := ctrl.DB.Select(&emails, query)
		if err != nil {
			return nil, err
		}
		// Merge dengan test emails
		return append(recipients, emails...), nil
	}

	if target == "warmup" {
		var emails []string
		query := "SELECT DISTINCT email FROM pdpi_members WHERE email IS NOT NULL AND email != '' LIMIT 50"
		err := ctrl.DB.Select(&emails, query)
		if err != nil {
			return nil, err
		}
		// Merge 50 members + 7 test emails
		return append(recipients, emails...), nil
	}

	// Default: Test list only (7 emails)
	return recipients, nil
}

// getWARecipients determines who to WhatsApp based on query param
func (ctrl *BroadcastController) getWARecipients(c *gin.Context) ([]string, error) {
	target := c.Query("target")

	if target == "all" {
		var phones []string
		query := "SELECT DISTINCT no_hp FROM pdpi_members WHERE no_hp IS NOT NULL AND no_hp != ''"
		err := ctrl.DB.Select(&phones, query)
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

	// Default: Test numbers requested by user
	return waTestRecipients, nil
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

func (ctrl *BroadcastController) BroadcastBeritaWA(c *gin.Context) {
	id := c.Param("id")

	// 1. Fetch Berita
	var berita struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Status  string `json:"status"`
		Slug    string `json:"slug"`
		Excerpt string `json:"excerpt"`
	}

	query := "SELECT id, title, status, slug, excerpt FROM berita WHERE id = ?"
	err := ctrl.DB.QueryRow(query, id).Scan(&berita.ID, &berita.Title, &berita.Status, &berita.Slug, &berita.Excerpt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita not found"})
		return
	}

	// 2. Validate Status
	if berita.Status != "published" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only published news can be shared via WhatsApp"})
		return
	}

	// 3. Prepare Message Content
	baseURL := ctrl.Config.BaseURL
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	readMoreLink := fmt.Sprintf("%s/berita/%s", baseURL, berita.Slug)
	message := fmt.Sprintf("*[PDPI News]*\n\n*%s*\n\n%s\n\nBaca selengkapnya:\n%s",
		berita.Title,
		berita.Excerpt,
		readMoreLink,
	)

	// 4. Determine Recipients
	targetRecipients, err := ctrl.getWARecipients(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipients: " + err.Error()})
		return
	}

	if len(targetRecipients) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid phone numbers found"})
		return
	}

	// 5. Send WA via Background Goroutine with delay
	go func(to []string, msg string) {
		fmt.Printf("Starting WA broadcast for %d recipients...\n", len(to))
		for i, number := range to {
			if err := ctrl.WAService.SendMessage(number, msg); err != nil {
				fmt.Printf("[%d/%d] Error sending WA to %s: %v\n", i+1, len(to), number, err)
			}
			// Jeda 500ms - 1s untuk menghindari blokir masal oleh WhatsApp
			time.Sleep(800 * time.Millisecond)
		}
		fmt.Printf("WA broadcast completed for %d recipients.\n", len(to))
	}(targetRecipients, message)

	c.JSON(http.StatusOK, gin.H{
		"message":         "WhatsApp broadcast process started in background",
		"recipient_count": len(targetRecipients),
		"target":          c.DefaultQuery("target", "test"),
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

func (ctrl *BroadcastController) BroadcastAgendaWA(c *gin.Context) {
	id := c.Param("id")

	// 1. Fetch Agenda
	var agenda struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Date        string `json:"date"`
		Slug        string `json:"slug"`
	}

	query := "SELECT id, title, description, status, date, slug FROM agenda WHERE id = ?"
	err := ctrl.DB.QueryRow(query, id).Scan(&agenda.ID, &agenda.Title, &agenda.Description, &agenda.Status, &agenda.Date, &agenda.Slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agenda not found"})
		return
	}

	// 2. Validate Status
	if agenda.Status != "published" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only published agenda can be shared via WhatsApp"})
		return
	}

	// 3. Prepare Message Content
	baseURL := ctrl.Config.BaseURL
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	readMoreLink := fmt.Sprintf("%s/agenda/%s", baseURL, agenda.Slug)
	message := fmt.Sprintf("*[PDPI Agenda]*\n\n*%s*\n\nTanggal: %s\n\nLihat detail agenda:\n%s",
		agenda.Title,
		agenda.Date,
		readMoreLink,
	)

	// 4. Determine Recipients
	targetRecipients, err := ctrl.getWARecipients(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipients: " + err.Error()})
		return
	}

	if len(targetRecipients) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid phone numbers found"})
		return
	}

	// 5. Send WA via Background Goroutine with delay
	go func(to []string, msg string) {
		fmt.Printf("Starting WA (Agenda) broadcast for %d recipients...\n", len(to))
		for i, number := range to {
			if err := ctrl.WAService.SendMessage(number, msg); err != nil {
				fmt.Printf("[%d/%d] Error sending WA (Agenda) to %s: %v\n", i+1, len(to), number, err)
			}
			time.Sleep(800 * time.Millisecond)
		}
		fmt.Printf("WA (Agenda) broadcast completed for %d recipients.\n", len(to))
	}(targetRecipients, message)

	c.JSON(http.StatusOK, gin.H{
		"message":         "WhatsApp Agenda broadcast started in background",
		"recipient_count": len(targetRecipients),
		"target":          c.DefaultQuery("target", "test"),
	})
}

func (ctrl *BroadcastController) TriggerBirthdayGreetings(c *gin.Context) {
	if err := ctrl.BirthdayService.CheckAndSendGreetings(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Birthday greetings process triggered successfully",
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
				<a href="%s" class="button" style="color: #ffffff !important; text-decoration: none;">%s</a>
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
