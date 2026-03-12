package services

import (
	"fmt"
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
		Nama  string  `db:"nama"`
		Email *string `db:"email"`
		NoHP  *string `db:"no_hp"`
	}

	// Query for MySQL
	query := `
		SELECT nama, email, no_hp 
		FROM pdpi_members 
		WHERE MONTH(tgl_lahir) = ? AND DAY(tgl_lahir) = ?
	`
	err := s.DB.Select(&members, query, month, day)
	if err != nil {
		return fmt.Errorf("failed to fetch birthday members: %w", err)
	}

	if len(members) == 0 {
		fmt.Println("No members having birthday today.")
		return nil
	}

	fmt.Printf("Found %d members having birthday today.\n", len(members))

	for _, m := range members {
		// 1. Send Email
		if m.Email != nil && *m.Email != "" {
			subject := fmt.Sprintf("Selamat Ulang Tahun, %s! 🎂", m.Nama)
			body := s.getBirthdayEmailTemplate(m.Nama)
			go func(to string, subj, content string) {
				if err := s.MailService.SendEmail([]string{to}, subj, content); err != nil {
					fmt.Printf("Error sending birthday email to %s: %v\n", to, err)
				}
			}(*m.Email, subject, body)
		}

		// 2. Send WhatsApp
		if m.NoHP != nil && *m.NoHP != "" {
			normalized := utils.NormalizePhoneNumber(*m.NoHP)
			if normalized != "" {
				message := fmt.Sprintf("Halo %s 👋,\n\nKami segenap keluarga besar *PDPI (Perhimpunan Dokter Paru Indonesia)* mengucapkan:\n\n✨ *Selamat Ulang Tahun!* ✨\n\nSemoga panjang umur, sehat selalu, dan sukses dalam menjalankan tugas mulia bagi bangsa dan sesama.\n\nSalam Hangat,\n*PDPI Pusat*", m.Nama)
				
				go func(to, msg string) {
					// Add slight delay between WA messages if there are many
					time.Sleep(1 * time.Second)
					if err := s.WAService.SendMessage(to, msg); err != nil {
						fmt.Printf("Error sending birthday WA to %s: %v\n", m.Nama, err)
					}
				}(normalized, message)
			}
		}
	}

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
