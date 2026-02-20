package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cvudumbarainformatika/backend/config"
)

// PDPIService handles communication with PDPI API (now Supabase)
type PDPIService struct {
	config *config.PDPIConfig
	client *http.Client
}

// NewPDPIService creates a new PDPI service instance
func NewPDPIService(cfg *config.PDPIConfig) *PDPIService {
	return &PDPIService{
		config: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// SupabaseMember represents a member from Supabase public_member_directory view
type SupabaseMember struct {
	ID                    string  `json:"id"`
	NPA                   *string `json:"npa"`
	NPANumeric            *int64  `json:"npa_numeric"`
	Nama                  string  `json:"nama"`
	Foto                  *string `json:"foto"`
	Gelar                 *string `json:"gelar"`
	Gelar2                *string `json:"gelar2"`
	Email                 *string `json:"email"`
	NoHP                  *string `json:"no_hp"`
	NIK                   *string `json:"nik"`
	JenisKelamin          *string `json:"jenis_kelamin"`
	TempatLahir           *string `json:"tempat_lahir"`
	TglLahir              *string `json:"tgl_lahir"` // Date string YYYY-MM-DD
	AlamatRumah           *string `json:"alamat_rumah"`
	Cabang                *string `json:"cabang"`
	Provinsi              *string `json:"provinsi"`
	KotaKabupaten         *string `json:"kota_kabupaten"`
	KotaKabupatenKantor   *string `json:"kota_kabupaten_kantor"`
	ProvinsiKantor        *string `json:"provinsi_kantor"`
	Status                *string `json:"status"`
	Alumni                *string `json:"alumni"`
	ThnLulus              *int64  `json:"thn_lulus"`
	TempatTugas           *string `json:"tempat_tugas"`
	TempatPraktek1        *string `json:"tempat_praktek_1"`
	TempatPraktek1Tipe    *string `json:"tempat_praktek_1_tipe"`
	TempatPraktek1Tipe2   *string `json:"tempat_praktek_1_tipe_2"`
	TempatPraktek1Alkes   *string `json:"tempat_praktek_1_alkes"`
	TempatPraktek1Alkes2  *string `json:"tempat_praktek_1_alkes_2"`
	TempatPraktek2        *string `json:"tempat_praktek_2"`
	TempatPraktek2Tipe    *string `json:"tempat_praktek_2_tipe"`
	TempatPraktek2Tipe2   *string `json:"tempat_praktek_2_tipe_2"`
	TempatPraktek2Alkes   *string `json:"tempat_praktek_2_alkes"`
	TempatPraktek2Alkes2  *string `json:"tempat_praktek_2_alkes_2"`
	KotaKabupatenPraktek2 *string `json:"kota_kabupaten_praktek_2"`
	ProvinsiPraktek2      *string `json:"provinsi_praktek_2"`
	TempatPraktek3        *string `json:"tempat_praktek_3"`
	TempatPraktek3Tipe    *string `json:"tempat_praktek_3_tipe"`
	TempatPraktek3Tipe2   *string `json:"tempat_praktek_3_tipe_2"`
	TempatPraktek3Alkes   *string `json:"tempat_praktek_3_alkes"`
	TempatPraktek3Alkes2  *string `json:"tempat_praktek_3_alkes_2"`
	KotaKabupatenPraktek3 *string `json:"kota_kabupaten_praktek_3"`
	ProvinsiPraktek3      *string `json:"provinsi_praktek_3"`
	Subspesialis          *string `json:"subspesialis"`
	GelarFISR             *string `json:"gelar_fisr"`
	NoSTR                 *string `json:"no_str"`
	STRBerlakuSampai      *string `json:"str_berlaku_sampai"` // Date string YYYY-MM-DD
	NoSIP                 *string `json:"no_sip"`
	SIPBerlakuSampai      *string `json:"sip_berlaku_sampai"` // Date string YYYY-MM-DD
	CreatedAt             *string `json:"created_at"`         // Timestamp string
	UpdatedAt             *string `json:"updated_at"`         // Timestamp string
}

// MembersFilter represents filters for GetMembers
type MembersFilter struct {
	Page     int
	Limit    int
	Cabang   string
	Provinsi string
	Status   string
	Search   string
}

// FetchMembersFromSupabase retrieves all members from Supabase REST API
func (s *PDPIService) FetchMembersFromSupabase() ([]SupabaseMember, error) {
	// Endpoint: /public_member_directory?select=*&order=npa_numeric.asc.nullslast
	var allMembers []SupabaseMember
	limit := 1000
	offset := 0

	for {
		endpoint := fmt.Sprintf("/public_member_directory?select=*&order=npa_numeric.asc.nullslast&limit=%d&offset=%d", limit, offset)
		url := s.getBaseURL() + endpoint

		fmt.Printf("[PDPI Sync] Fetching page %d (offset %d): %s\n", (offset/limit)+1, offset, url)

		// Debug API Key (Safely)
		maskedKey := "N/A"
		if len(s.config.APIKey) > 8 {
			maskedKey = s.config.APIKey[:4] + "..." + s.config.APIKey[len(s.config.APIKey)-4:]
		}
		fmt.Printf("[PDPI Debug] API Key: %s (Length: %d)\n", maskedKey, len(s.config.APIKey))

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set Authentication Header
		req.Header.Set("apikey", s.config.APIKey)
		// req.Header.Set("Authorization", "Bearer "+s.config.APIKey) // Removed based on user feedback
		req.Header.Set("Content-Type", "application/json")

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to execute request: %w", err)
		}
		defer resp.Body.Close()

		// Read response body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// Check for error response
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Supabase API error (status %d) at URL %s: %s", resp.StatusCode, url, string(bodyBytes))
		}

		// Parse success response array
		var members []SupabaseMember
		if err := json.Unmarshal(bodyBytes, &members); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		if len(members) == 0 {
			break
		}

		allMembers = append(allMembers, members...)

		if len(members) < limit {
			break
		}

		offset += limit
	}

	return allMembers, nil
}

// FetchMemberByNPA retrieves a specific member by NPA from Supabase
func (s *PDPIService) FetchMemberByNPA(npa string) (*SupabaseMember, error) {
	// Filter: npa=eq.VALUE
	endpoint := fmt.Sprintf("/public_member_directory?npa=eq.%s&limit=1", npa)
	url := s.getBaseURL() + endpoint

	fmt.Printf("[PDPI Sync] Fetching member by NPA: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", s.config.APIKey)
	// req.Header.Set("Authorization", "Bearer "+s.config.APIKey) // Removed
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Supabase API error (status %d) at URL %s", resp.StatusCode, url)
	}

	var members []SupabaseMember
	if err := json.Unmarshal(bodyFromResponse(resp), &members); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("member not found")
	}

	return &members[0], nil
}

// FetchMemberByEmail retrieves a specific member by Email from Supabase
func (s *PDPIService) FetchMemberByEmail(email string) (*SupabaseMember, error) {
	// Filter: email=eq.VALUE
	endpoint := fmt.Sprintf("/public_member_directory?email=eq.%s&limit=1", email)
	url := s.getBaseURL() + endpoint

	fmt.Printf("[PDPI Sync] Fetching member by Email: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", s.config.APIKey)
	// req.Header.Set("Authorization", "Bearer "+s.config.APIKey) // Removed
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Supabase API error (status %d) at URL %s", resp.StatusCode, url)
	}

	var members []SupabaseMember
	if err := json.Unmarshal(bodyFromResponse(resp), &members); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("member not found")
	}

	return &members[0], nil
}

// Helper to read body safely
func bodyFromResponse(resp *http.Response) []byte {
	body, _ := io.ReadAll(resp.Body)
	return body
}

// getBaseURL handles sanitization and common config fixes
func (s *PDPIService) getBaseURL() string {
	baseURL := strings.TrimRight(s.config.BaseURL, "/")

	// Auto-fix common mistake: user copied Edge Function URL (/functions/v1) instead of REST URL (/rest/v1)
	if strings.Contains(baseURL, "/functions/v1") {
		baseURL = strings.Replace(baseURL, "/functions/v1", "/rest/v1", 1)
		fmt.Printf("[PDPI WARN] Auto-corrected BaseURL: replaced /functions/v1/ with /rest/v1/\n")
	}

	return baseURL
}
