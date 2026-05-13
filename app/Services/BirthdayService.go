package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/jmoiron/sqlx"
)

type BirthdayService struct {
	DB          *sqlx.DB
	MailService *MailService
	WAService   *WhatsAppService
}

func NewBirthdayService(db *sqlx.DB, mailService *MailService, waService *WhatsAppService) *BirthdayService {
	return &BirthdayService{
		DB:          db,
		MailService: mailService,
		WAService:   waService,
	}
}

// CheckAndSendGreetings finds members celebrating birthday today and sends greetings
func (s *BirthdayService) CheckAndSendGreetings() error {
	now := time.Now()
	month := int(now.Month())
	day := now.Day()

	fmt.Printf("[%s] Checking for birthdays today (%02d-%02d)...\n", now.Format("2006-01-02 15:04:05"), day, month)

	var members []struct {
		Nama   string  `db:"nama"`
		Email  *string `db:"email"`
		NoHP   *string `db:"no_hp"`
		Gelar  *string `db:"gelar"`
		Gelar2 *string `db:"gelar2"`
	}

	// Query for MySQL
	query := `
		SELECT nama, email, no_hp, gelar, gelar2
		FROM pdpi_members 
		WHERE MONTH(tgl_lahir) = ? AND DAY(tgl_lahir) = ?
	`
	err := s.DB.Select(&members, query, month, day)
	if err != nil {
		return fmt.Errorf("failed to fetch birthday members: %w", err)
	}

	if len(members) == 0 {
		fmt.Println("[BIRTHDAY] No members having birthday today.")
		return nil
	}

	fmt.Printf("[BIRTHDAY] Found %d members having birthday today. Starting sequential broadcast...\n", len(members))

	// Run the broadcast in a single background goroutine
	go func(targets []struct {
		Nama   string  `db:"nama"`
		Email  *string `db:"email"`
		NoHP   *string `db:"no_hp"`
		Gelar  *string `db:"gelar"`
		Gelar2 *string `db:"gelar2"`
	}) {
		total := len(targets)
		for i, m := range targets {
			msgNum := i + 1
			
			// 1. Send Email
			if m.Email != nil && *m.Email != "" {
				subject := fmt.Sprintf("Selamat Ulang Tahun, %s! 🎂", m.Nama)
				body := s.getBirthdayEmailTemplate(m.Nama)
				fmt.Printf("[BIRTHDAY-EMAIL] [%d/%d] Sending to %s...\n", msgNum, total, *m.Email)
				if err := s.MailService.SendEmail([]string{*m.Email}, subject, body); err != nil {
					fmt.Printf("[BIRTHDAY-EMAIL ERROR] to %s: %v\n", *m.Email, err)
				}
			}

			// 2. Send WhatsApp
			if m.NoHP != nil && *m.NoHP != "" {
				normalized := utils.NormalizePhoneNumber(*m.NoHP)
				if normalized != "" {
					// Prepare parameters from database
					gelarStr := ""
					if m.Gelar != nil {
						gelarStr = *m.Gelar
					}
					gelar2Str := ""
					if m.Gelar2 != nil {
						gelar2Str = *m.Gelar2
					}
					
					fmt.Printf("[BIRTHDAY-WA] [%d/%d] Sending to %s...\n", msgNum, total, normalized)
					if err := s.WAService.SendUltah(normalized, gelarStr, m.Nama, gelar2Str); err != nil {
						fmt.Printf("[BIRTHDAY-WA ERROR] to %s: %v\n", m.Nama, err)
						s.WAService.LogEvent(normalized, "failed")
					} else {
						s.WAService.LogEvent(normalized, "sent")
					}
				}
			}

			// 3. Anti-Spam Delays
			if msgNum < total {
				delay := 3 + rand.Intn(4) // 3-6 seconds
				
				if msgNum % 40 == 0 {
					sleepMinutes := 3 + rand.Intn(4) // 3-6 minutes
					fmt.Printf("[BIRTHDAY] Batch %d selesai, istirahat %d menit untuk anti-spam...\n", 
						msgNum/40, sleepMinutes)
					time.Sleep(time.Duration(sleepMinutes) * time.Minute)
				} else {
					time.Sleep(time.Duration(delay) * time.Second)
				}
			}
		}
		fmt.Printf("[BIRTHDAY] Completed sending greetings to %d members.\n", total)
	}(members)

	return nil
}

func (s *BirthdayService) getBirthdayEmailTemplate(nama string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: 'Helvetica Neue', Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 20px auto; border: 1px solid #eee; border-radius: 10px; overflow: hidden; }
        .header { background-color: #2563eb; color: white; padding: 40px 20px; text-align: center; }
        .content { padding: 30px; background-color: #ffffff; text-align: center; }
        .footer { background-color: #f9fafb; padding: 20px; text-align: center; font-size: 12px; color: #6b7280; }
        h1 { margin: 0; font-size: 28px; }
        p { font-size: 16px; margin-bottom: 20px; }
        .highlight { color: #2563eb; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Selamat Ulang Tahun! 🎂</h1>
        </div>
        <div class="content">
            <p>Halo <span class="highlight">%s</span>,</p>
            <p>Segenap keluarga besar <strong>Perhimpunan Dokter Paru Indonesia (PDPI)</strong> mengucapkan selamat hari ulang tahun.</p>
            <p>Semoga di usia yang baru ini, Anda selalu diberikan kesehatan, kebahagiaan, dan keberkahan dalam menjalankan tugas mulia sebagai dokter paru.</p>
            <p>Teruslah menginspirasi dan memberikan yang terbaik bagi kesehatan masyarakat Indonesia.</p>
            <br>
            <p>Salam hangat,</p>
            <p><strong>Pengurus Pusat PDPI</strong></p>
        </div>
        <div class="footer">
            <p>&copy; 2026 Perhimpunan Dokter Paru Indonesia</p>
        </div>
    </div>
</body>
</html>
	`, nama)
}
