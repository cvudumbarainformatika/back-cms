package controllers

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	"6281350125649",
	"6281381295959",
	"6281554545012",
	"628155088994",
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

	// Let's re-fetch with ImageURL
	var beritaWithImg struct {
		Title    string `db:"title"`
		Slug     string `db:"slug"`
		ImageURL string `db:"image_url"`
		Excerpt  string `db:"excerpt"`
	}
	ctrl.DB.Get(&beritaWithImg, "SELECT title, slug, image_url, excerpt FROM berita WHERE id = ?", id)
	
	readMoreLink := fmt.Sprintf("%s/berita/%s", baseURL, beritaWithImg.Slug)
	
	fullImageURL := ""
	if ctrl.Config.Env == "local" {
		fullImageURL = "https://hobiternak.com/wp-content/uploads/2017/08/bebek-unggulan.jpg"
	} else if beritaWithImg.ImageURL != "" {
		if strings.HasPrefix(beritaWithImg.ImageURL, "/") {
			fullImageURL = baseURL + beritaWithImg.ImageURL
		} else {
			fullImageURL = beritaWithImg.ImageURL
		}
	} else {
		fullImageURL = baseURL + "/images/default-news.jpg"
	}

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
	go func(to []string, title, url, img string) {
		fmt.Printf("[WA] Starting broadcast for %d recipients using 360dialog template...\n", len(to))
		for i, number := range to {
			msgNum := i + 1
			if err := ctrl.WAService.SendArtikel(number, title, url, img); err != nil {
				fmt.Printf("[WA ERROR] [%d/%d] to %s: %v\n", msgNum, len(to), number, err)
				ctrl.WAService.LogEvent(number, "failed")
			} else {
				fmt.Printf("[WA] [%d/%d] Sent to %s\n", msgNum, len(to), number)
				ctrl.WAService.LogEvent(number, "sent")
			}
			
			// Delay random 5-10 seconds to avoid spam flagging
			if msgNum < len(to) {
				delay := 5 + rand.Intn(6)
				time.Sleep(time.Duration(delay) * time.Second)
			}
		}
		fmt.Printf("[WA] Broadcast completed for %d recipients.\n", len(to))
	}(targetRecipients, beritaWithImg.Title, readMoreLink, fullImageURL)

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
		ID     int    `json:"id"`
		Status string `json:"status"`
	}

	query := "SELECT id, status FROM agenda WHERE id = ?"
	err := ctrl.DB.QueryRow(query, id).Scan(&agenda.ID, &agenda.Status)
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

	// 3. Fetch full data for Agenda template
	var ag struct {
		Title           string    `db:"title"`
		IsOnline        bool      `db:"is_online"`
		Location        string    `db:"location"`
		Date            time.Time `db:"date"`
		Fee             string    `db:"fee"`
		Quota           int       `db:"quota"`
		RegistrationURL string    `db:"registration_url"`
		ImageURL        string    `db:"image_url"`
		Slug            string    `db:"slug"`
	}
	err = ctrl.DB.Get(&ag, "SELECT title, is_online, location, date, fee, quota, registration_url, image_url, slug FROM agenda WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch agenda data: " + err.Error()})
		return
	}

	readMoreLink := fmt.Sprintf("%s/agenda/%s", baseURL, ag.Slug)
	
	mode := "Luring"
	if ag.IsOnline {
		mode = "Daring"
	}
	
	location := ag.Location
	if ag.IsOnline && location == "" {
		location = "Zoom/Cloud Meeting"
	}
	
	timeStr := ag.Date.Format("Monday, 02 January 2006 pukul 15:04 MST")
	// Replace English names with Indonesian names manually
	timeStr = strings.NewReplacer(
		"Monday", "Senin", "Tuesday", "Selasa", "Wednesday", "Rabu",
		"Thursday", "Kamis", "Friday", "Jumat", "Saturday", "Sabtu", "Sunday", "Minggu",
		"January", "Januari", "February", "Februari", "March", "Maret", "May", "Mei",
		"June", "Juni", "July", "Juli", "August", "Agustus", "October", "Oktober", "December", "Desember",
	).Replace(timeStr)

	fee := ag.Fee
	if fee == "" || fee == "0" {
		fee = "Gratis"
	}
	
	quota := strconv.Itoa(ag.Quota)
	if ag.Quota == 0 {
		quota = "Tanpa Batas"
	}

	regURL := ag.RegistrationURL
	if regURL == "" {
		regURL = readMoreLink
	}

	// Image URL
	fullImageURL := ""
	if ctrl.Config.Env == "local" {
		fullImageURL = "https://hobiternak.com/wp-content/uploads/2017/08/bebek-unggulan.jpg"
	} else if ag.ImageURL != "" {
		if strings.HasPrefix(ag.ImageURL, "/") {
			fullImageURL = baseURL + ag.ImageURL
		} else {
			fullImageURL = ag.ImageURL
		}
	} else {
		fullImageURL = baseURL + "/images/default-agenda.jpg"
	}

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
	go func(to []string) {
		fmt.Printf("[WA-AGENDA] Starting broadcast for %d recipients using 360dialog template...\n", len(to))
		for i, number := range to {
			msgNum := i + 1
			err := ctrl.WAService.SendAgenda(number, ag.Title, mode, location, timeStr, fee, quota, regURL, fullImageURL)
			
			if err != nil {
				fmt.Printf("[WA-AGENDA ERROR] [%d/%d] to %s: %v\n", msgNum, len(to), number, err)
				ctrl.WAService.LogEvent(number, "failed")
			} else {
				fmt.Printf("[WA-AGENDA] [%d/%d] Sent to %s\n", msgNum, len(to), number)
				ctrl.WAService.LogEvent(number, "sent")
			}

			// Delay random 5-10 seconds to avoid spam flagging
			if msgNum < len(to) {
				delay := 5 + rand.Intn(6)
				time.Sleep(time.Duration(delay) * time.Second)
			}
		}
		fmt.Printf("[WA-AGENDA] Broadcast completed for %d recipients.\n", len(to))
	}(targetRecipients)

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
func (ctrl *BroadcastController) GetEmailLogs(c *gin.Context) {
	sinceStr := c.Query("since")
	var sinceTime time.Time
	if sinceStr != "" {
		if sec, err := strconv.ParseInt(sinceStr, 10, 64); err == nil {
			sinceTime = time.Unix(sec, 0)
		}
	}

	logPath := "/var/log/mail/mail.log" // Reaches mapped volume in production
	// Fallback for local testing
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		logPath = "logs/mail.log"
	}

	file, err := os.Open(logPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success":  []string{},
			"deferred": []string{},
			"message":  "Log file not found or inaccessible",
		})
		return
	}
	defer file.Close()

	// Regexes
	sentRegex := regexp.MustCompile(`to=<([^>]+)>,.*status=sent`)
	deferredRegex := regexp.MustCompile(`to=<([^>]+)>,.*status=deferred`)
	// Postfix RFC3339 timestamp format: 2026-04-02T12:10:30.807970+00:00
	timestampRegex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})`)

	sentMap := make(map[string]bool)
	deferredMap := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Time filtering
		if !sinceTime.IsZero() {
			if tsMatches := timestampRegex.FindStringSubmatch(line); len(tsMatches) > 1 {
				// We use a shorter layout for Parse "2006-01-02T15:04:05"
				lineTime, err := time.Parse("2006-01-02T15:04:05", tsMatches[1])
				if err == nil {
					// We need to account for the year/location if log doesn't have it
					// but RFC3339 usually has everything. Postfix here uses UTC.
					if lineTime.Before(sinceTime.UTC()) {
						continue
					}
				}
			}
		}

		if matches := sentRegex.FindStringSubmatch(line); len(matches) > 1 {
			email := matches[1]
			sentMap[email] = true
			delete(deferredMap, email)
		} else if matches := deferredRegex.FindStringSubmatch(line); len(matches) > 1 {
			email := matches[1]
			if !sentMap[email] {
				deferredMap[email] = true
			}
		}
	}

	sentEmails := make([]string, 0, len(sentMap))
	for email := range sentMap {
		sentEmails = append(sentEmails, email)
	}

	deferredEmails := make([]string, 0, len(deferredMap))
	for email := range deferredMap {
		deferredEmails = append(deferredEmails, email)
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  sentEmails,
		"deferred": deferredEmails,
	})
}
func (ctrl *BroadcastController) GetWhatsAppLogs(c *gin.Context) {
	sinceStr := c.Query("since")
	var sinceTime time.Time
	if sinceStr != "" {
		if sec, err := strconv.ParseInt(sinceStr, 10, 64); err == nil {
			sinceTime = time.Unix(sec, 0)
		}
	}

	logPath := "logs/whatsapp.log"
	file, err := os.Open(logPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": []string{},
			"failed":  []string{},
			"message": "WA Log file not found",
		})
		return
	}
	defer file.Close()

	// Regex for WhatsApp log: 2026-04-04T12:10:30+00:00 status=sent to=62812...
	// RFC3339 timestamp regex
	timestampRegex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})`)
	statusRegex := regexp.MustCompile(`status=(\w+)`)
	toRegex := regexp.MustCompile(`to=([^\s]+)`)

	sentList := make(map[string]bool)
	failedList := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Time filter
		if !sinceTime.IsZero() {
			if tsMatches := timestampRegex.FindStringSubmatch(line); len(tsMatches) > 1 {
				lineTime, err := time.Parse("2006-01-02T15:04:05", tsMatches[1])
				if err == nil {
					if lineTime.Before(sinceTime.UTC()) {
						continue
					}
				}
			}
		}

		statusMatches := statusRegex.FindStringSubmatch(line)
		toMatches := toRegex.FindStringSubmatch(line)

		if len(statusMatches) > 1 && len(toMatches) > 1 {
			status := statusMatches[1]
			to := toMatches[1]

			if status == "sent" {
				sentList[to] = true
				delete(failedList, to)
			} else {
				if !sentList[to] {
					failedList[to] = true
				}
			}
		}
	}

	success := make([]string, 0, len(sentList))
	for to := range sentList {
		success = append(success, to)
	}

	failed := make([]string, 0, len(failedList))
	for to := range failedList {
		failed = append(failed, to)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"failed":  failed,
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
