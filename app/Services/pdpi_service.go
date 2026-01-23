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

// PDPIService handles communication with PDPI API
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

// PDPIMember represents a member from PDPI API
type PDPIMember struct {
	ID               string `json:"id"`
	NPA              string `json:"npa"`
	Nama             string `json:"nama"`
	Gelar            string `json:"gelar"`
	Gelar2           string `json:"gelar2"`
	Email            string `json:"email"`
	NoHP             string `json:"no_hp"`
	NIK              string `json:"nik"`
	JenisKelamin     string `json:"jenis_kelamin"`
	TempatLahir      string `json:"tempat_lahir"`
	TglLahir         string `json:"tgl_lahir"`
	AlamatRumah      string `json:"alamat_rumah"`
	Cabang           string `json:"cabang"`
	Provinsi         string `json:"provinsi"`
	KotaKabupaten    string `json:"kota_kabupaten"`
	Status           string `json:"status"`
	Alumni           string `json:"alumni"`
	ThnLulus         int    `json:"thn_lulus"`
	TempatTugas      string `json:"tempat_tugas"`
	TempatPraktek1   string `json:"tempat_praktek_1"`
	TempatPraktek2   string `json:"tempat_praktek_2"`
	Subspesialis     string `json:"subspesialis"`
	NoSTR            string `json:"no_str"`
	STRBerlakuSampai string `json:"str_berlaku_sampai"`
	NoSIP            string `json:"no_sip"`
	SIPBerlakuSampai string `json:"sip_berlaku_sampai"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// PDPILoginRequest represents login request to PDPI API
type PDPILoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// PDPILoginResponse represents login response from PDPI API
type PDPILoginResponse struct {
	Success bool `json:"success"`
	User    struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		NIK      string `json:"nik"`
		Role     string `json:"role"`
		BranchID string `json:"branch_id"`
	} `json:"user"`
	Member  PDPIMember `json:"member"`
	Session struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		ExpiresAt    int64  `json:"expires_at"`
		TokenType    string `json:"token_type"`
	} `json:"session"`
}

// PDPIMembersResponse represents response from /api-members endpoint
type PDPIMembersResponse struct {
	Success    bool         `json:"success"`
	Data       []PDPIMember `json:"data"`
	Pagination struct {
		Page       int `json:"page"`
		Limit      int `json:"limit"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	} `json:"pagination"`
}

// PDPIErrorResponse represents error response from PDPI API
type PDPIErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
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

// doRequest performs HTTP request to PDPI API
func (s *PDPIService) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	url := s.config.BaseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-api-key", s.config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// Login authenticates a member via PDPI API
func (s *PDPIService) Login(email, password string) (*PDPILoginResponse, error) {
	reqBody := PDPILoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := s.doRequest("POST", "/api-login", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		var errResp PDPIErrorResponse
		if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
			return nil, fmt.Errorf("PDPI API error (status %d): %s", resp.StatusCode, string(bodyBytes))
		}
		return nil, fmt.Errorf("PDPI API error: %s - %s", errResp.Error, errResp.Message)
	}

	// Parse success response
	var loginResp PDPILoginResponse
	if err := json.Unmarshal(bodyBytes, &loginResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !loginResp.Success {
		return nil, fmt.Errorf("login failed: response indicates failure")
	}

	return &loginResp, nil
}

// GetMembers retrieves members list from PDPI API
func (s *PDPIService) GetMembers(filter MembersFilter) (*PDPIMembersResponse, error) {
	// Build query parameters
	endpoint := "/api-members?"

	if filter.Page > 0 {
		endpoint += fmt.Sprintf("page=%d&", filter.Page)
	} else {
		endpoint += "page=1&"
	}

	if filter.Limit > 0 {
		endpoint += fmt.Sprintf("limit=%d&", filter.Limit)
	} else {
		endpoint += "limit=100&"
	}

	if filter.Cabang != "" {
		endpoint += fmt.Sprintf("cabang=%s&", filter.Cabang)
	}

	if filter.Provinsi != "" {
		endpoint += fmt.Sprintf("provinsi=%s&", filter.Provinsi)
	}

	if filter.Status != "" {
		endpoint += fmt.Sprintf("status=%s&", filter.Status)
	}

	if filter.Search != "" {
		endpoint += fmt.Sprintf("search=%s&", filter.Search)
	}

	resp, err := s.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		var errResp PDPIErrorResponse
		if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
			return nil, fmt.Errorf("PDPI API error (status %d): %s", resp.StatusCode, string(bodyBytes))
		}
		return nil, fmt.Errorf("PDPI API error: %s - %s", errResp.Error, errResp.Message)
	}

	// Parse success response
	var membersResp PDPIMembersResponse
	if err := json.Unmarshal(bodyBytes, &membersResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &membersResp, nil
}

// GetMemberByEmail retrieves a specific member by email
func (s *PDPIService) GetMemberByEmail(email string) (*PDPIMember, error) {
	filter := MembersFilter{
		Page:   1,
		Limit:  1,
		Search: email,
	}

	resp, err := s.GetMembers(filter)
	if err != nil {
		return nil, err
	}

	if !resp.Success || len(resp.Data) == 0 {
		return nil, fmt.Errorf("member not found with email: %s", email)
	}

	// Return first match
	return &resp.Data[0], nil
}

// GetMemberByNPA retrieves a specific member by NPA
func (s *PDPIService) GetMemberByNPA(npa string) (*PDPIMember, error) {
	filter := MembersFilter{
		Page:   1,
		Limit:  1,
		Search: npa,
	}

	resp, err := s.GetMembers(filter)
	if err != nil {
		return nil, err
	}

	if !resp.Success || len(resp.Data) == 0 {
		return nil, fmt.Errorf("member not found with NPA: %s", npa)
	}

	// Find exact NPA match
	for _, member := range resp.Data {
		if member.NPA == npa {
			return &member, nil
		}
	}

	return nil, fmt.Errorf("member not found with NPA: %s", npa)
}
