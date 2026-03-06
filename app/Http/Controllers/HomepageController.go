package controllers

import (
	"net/http"

	requests "github.com/cvudumbarainformatika/backend/app/Http/Requests"
	models "github.com/cvudumbarainformatika/backend/app/Models"
	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type HomepageController struct {
	db *sqlx.DB
}

func NewHomepageController(db *sqlx.DB) *HomepageController {
	return &HomepageController{db: db}
}

func (hc *HomepageController) Get(c *gin.Context) {
	homepage, err := models.GetHomepage(hc.db)
	if err != nil {
		// If not found, return default structure or empty
		defaultContent := gin.H{
			"theme": "theme-2",
			"hero": gin.H{
				"title":       "Leading Respiratory Science.",
				"label":       "PERCOBAAN",
				"description": "Wadah profesional kesehatan paru dan respirasi untuk kemajuan sains, pelayanan medis, dan edukasi masyarakat.",
				"event_tag":   "PDPI",
				"event_title": "Coba Judul",
				"event_desc":  "Inovasi Penanganan PPOK & Asma",
				"images":      []string{},
			},
			"about": gin.H{
				"description": "Berdiri sejak 8 September 1973\n• Ikatan Dokter Paru Indonesia (IDPI): rapat 20 dokter ahli penyakit paru di Indonesia, termasuk Dr. Rasmin Rasjid memperjuangkan pembentukan Bagian Pulmonologi di FKUI.\n• 1988: nama IDPI diubah menjadi Perhimpunan Dokter Paru Indonesia (PDPI) dalam Kongres Nasional V-IDPI, disesuaikan dengan Muktamar IDI ke-20 di Surabaya\n• 2023 HUT PDPI 50 tahun\n• 2025 HUT PDPI 52 tahun",
				"photo":       "",
				"members":     []interface{}{},
			},
			"stats": []interface{}{},
			"footer": gin.H{
				"contact": gin.H{
					"phone": "+62 21 568 1149 ext 101 (sekretariat)\n+62 21 568 1149 ext 106 (kolegium)\n+62 21 568 1149 ext 108 (web-support)",
					"email": "sekretariat@pdpi.or.id\nkolegium@pdpi.or.id",
				},
				"socials": gin.H{
					"instagram": "official_pdpi",
					"youtube":   "Media PDPI",
				},
				"columns": []gin.H{
					{
						"title": "ALAMAT RUMAH PDPI",
						"items": []string{
							"Alamat Lengkap :",
							"Jl. Cipinang Bunder No.19",
							"Cipinang Pulogadung – Jakarta",
							"• Kode Pos : 13240",
							"• Telepon : (021) 22474845",
							"• Email : sekjen_pdpi@ymail.com",
							"• Website : www.klikpdpi.com",
						},
					},
					{
						"title": "ORGANISASI",
						"items": []string{
							"1. Badan Legislatif: KONAS, PIK, KONKER",
							"2. Majelis Kehormatan: 9 orang",
							"3. Dewan Pengawas",
							"4. Dewan Etik, Hukum & Pembelaan: 7 orang",
							"5. Dewan Eksekutif: Pengurus Harian",
							"6. Dewan Pendidikan",
							"7. Pengurus Cabang",
							"(AD-ART BAB VII Pasal 8)",
						},
					},
				},
			},
			"seo": gin.H{
				"title":       "PDPI",
				"description": "Perhimpunan Dokter Paru Indonesia",
			},
		}
		utils.Success(c, http.StatusOK, "Default homepage data (DB empty)", defaultContent)
		return
	}

	utils.Success(c, http.StatusOK, "Homepage data fetched", homepage.Content)
}

func (hc *HomepageController) Update(c *gin.Context) {
	var req requests.UpdateHomepageRequest
	if err := req.Validate(c); err != nil {
		return
	}

	// Prepare content map
	content := map[string]interface{}{
		"theme":  req.Theme,
		"hero":   req.Hero,
		"about":  req.About,
		"stats":  req.Stats,
		"footer": req.Footer,
		"seo":    req.SEO,
	}

	homepage, err := models.UpdateHomepage(hc.db, content)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "database_error", "Failed to update homepage", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Homepage updated successfully", homepage.Content)
}
