package services

import (
	"bobastream/config"
	"bobastream/internal/models"
	"bobastream/internal/repositories"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type PCloudService struct {
	pcloudRepo *repositories.PCloudCredentialRepository
	videoRepo  *repositories.VideoRepository
}

func NewPCloudService(
	pcloudRepo *repositories.PCloudCredentialRepository,
	videoRepo *repositories.VideoRepository,
) *PCloudService {
	return &PCloudService{
		pcloudRepo: pcloudRepo,
		videoRepo:  videoRepo,
	}
}

// PCloudUploadResponse response from pCloud upload API
type PCloudUploadResponse struct {
	Result   int    `json:"result"`
	Metadata []struct {
		Name   string `json:"name"`
		FileID int64  `json:"fileid"`
		Size   int64  `json:"size"`
		Hash   string `json:"hash"`
	} `json:"metadata"`
	Error string `json:"error,omitempty"`
}

// PCloudLinkResponse response from pCloud getfilelink API
type PCloudLinkResponse struct {
	Result  int      `json:"result"`
	Hosts   []string `json:"hosts"`
	Path    string   `json:"path"`
	Expires string   `json:"expires"`
	Error   string   `json:"error,omitempty"`
}

// GetAvailableAccount gets pCloud account with most available space (auto-rotate)
func (s *PCloudService) GetAvailableAccount() (*models.PCloudCredential, error) {
	credentials, err := s.pcloudRepo.GetActive()
	if err != nil {
		return nil, err
	}

	if len(credentials) == 0 {
		return nil, errors.New("no active pCloud accounts available")
	}

	// Find account with most available space
	var bestAccount *models.PCloudCredential
	maxAvailable := 0.0

	for i := range credentials {
		available := credentials[i].StorageLimitGB - credentials[i].StorageUsedGB
		if available > maxAvailable {
			maxAvailable = available
			bestAccount = &credentials[i]
		}
	}

	if bestAccount == nil || maxAvailable < 0.1 { // Less than 100MB available
		return nil, errors.New("no pCloud account with sufficient storage available")
	}

	return bestAccount, nil
}

// UploadFile uploads file to pCloud with auto-rotate account
func (s *PCloudService) UploadFile(file multipart.File, filename string, fileSize int64) (*models.PCloudCredential, int64, string, error) {
	// ✅ CALCULATE FILE SIZE IN GB
	fileSizeMB := float64(fileSize) / (1024 * 1024)
	fileSizeGB := fileSizeMB / 1024

	// Get available account (auto-rotate based on storage)
	credential, err := s.GetAvailableAccount()
	if err != nil {
		return nil, 0, "", err
	}

	// ✅ CHECK IF FILE SIZE EXCEEDS AVAILABLE STORAGE
	availableGB := credential.StorageLimitGB - credential.StorageUsedGB
	if fileSizeGB > availableGB {
		return nil, 0, "", fmt.Errorf(
			"file size (%.2fGB) exceeds available storage (%.2fGB) in account '%s'",
			fileSizeGB,
			availableGB,
			credential.AccountName,
		)
	}

	// ✅ CHECK IF FILE SIZE EXCEEDS 10% OF TOTAL STORAGE (safety limit)
	maxAllowedGB := credential.StorageLimitGB * 0.1 // 10% of total
	if fileSizeGB > maxAllowedGB {
		return nil, 0, "", fmt.Errorf(
			"file size (%.2fGB) exceeds maximum allowed file size (%.2fGB) for account '%s'",
			fileSizeGB,
			maxAllowedGB,
			credential.AccountName,
		)
	}

	// Upload to pCloud
	fileID, hash, err := s.uploadToPCloud(credential.APIToken, file, filename)
	if err != nil {
		return nil, 0, "", fmt.Errorf("pCloud upload failed: %w", err)
	}

	// ✅ UPDATE STORAGE USED (WITH SAFETY MARGIN 5%)
	storageIncrease := fileSizeGB * 1.05 // Add 5% overhead for metadata
	credential.StorageUsedGB += storageIncrease

	if err := s.pcloudRepo.Update(credential); err != nil {
		// ⚠️ NON-CRITICAL ERROR: File uploaded but storage tracking failed
		// Log this but don't fail the upload
		fmt.Printf("⚠️  WARNING: Failed to update storage tracking for account '%s': %v\n",
			credential.AccountName, err)
	}

	return credential, fileID, hash, nil
}

// uploadToPCloud actual upload to pCloud API
func (s *PCloudService) uploadToPCloud(apiToken string, file multipart.File, filename string) (int64, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return 0, "", err
	}

	if _, err := io.Copy(part, file); err != nil {
		return 0, "", err
	}

	writer.Close()

	// Create request
	url := fmt.Sprintf("%s/uploadfile?auth=%s&folderid=0", config.GlobalConfig.PCloud.BaseURL, apiToken)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	// Parse response
	var uploadResp PCloudUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return 0, "", err
	}

	if uploadResp.Result != 0 {
		return 0, "", fmt.Errorf("pCloud upload failed: %s", uploadResp.Error)
	}

	if len(uploadResp.Metadata) == 0 {
		return 0, "", errors.New("no file metadata returned from pCloud")
	}

	return uploadResp.Metadata[0].FileID, uploadResp.Metadata[0].Hash, nil
}

// GetFileLink gets streaming link from pCloud
func (s *PCloudService) GetFileLink(fileID int64, apiToken string) (string, time.Time, error) {
	url := fmt.Sprintf("%s/getfilelink?auth=%s&fileid=%d", config.GlobalConfig.PCloud.BaseURL, apiToken, fileID)

	resp, err := http.Get(url)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	var linkResp PCloudLinkResponse
	if err := json.NewDecoder(resp.Body).Decode(&linkResp); err != nil {
		return "", time.Time{}, err
	}

	if linkResp.Result != 0 {
		return "", time.Time{}, fmt.Errorf("pCloud getfilelink failed: %s", linkResp.Error)
	}

	// Construct full URL
	if len(linkResp.Hosts) == 0 {
		return "", time.Time{}, errors.New("no hosts returned from pCloud")
	}

	fullURL := fmt.Sprintf("https://%s%s", linkResp.Hosts[0], linkResp.Path)

	// Parse expiry (pCloud returns "Thu, 01 Jan 2026 00:00:00 +0000" format)
	expiresAt, _ := time.Parse(time.RFC1123, linkResp.Expires)
	if expiresAt.IsZero() {
		// Default to 24 hours from now if parse fails
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	return fullURL, expiresAt, nil
}

// RefreshExpiredLinks refreshes expired video source URLs (called by cron)
func (s *PCloudService) RefreshExpiredLinks() error {
	videos, err := s.videoRepo.GetExpiredSourceURLVideos()
	if err != nil {
		return err
	}

	for _, video := range videos {
		if video.PCloudCredential == nil {
			continue
		}

		// Get new link
		fileID := int64(0)
		fmt.Sscanf(video.PCloudFileID, "%d", &fileID)

		newURL, expiresAt, err := s.GetFileLink(fileID, video.PCloudCredential.APIToken)
		if err != nil {
			// Log error but continue with other videos
			continue
		}

		// Update video
		if err := s.videoRepo.UpdateSourceURL(video.ID, newURL, expiresAt); err != nil {
			continue
		}
	}

	return nil
}

// CreateCredential creates new pCloud credential (admin)
func (s *PCloudService) CreateCredential(credential *models.PCloudCredential) error {
	return s.pcloudRepo.Create(credential)
}

// UpdateCredential updates pCloud credential
func (s *PCloudService) UpdateCredential(credential *models.PCloudCredential) error {
	return s.pcloudRepo.Update(credential)
}

// DeleteCredential deletes pCloud credential
func (s *PCloudService) DeleteCredential(id uuid.UUID) error {
	return s.pcloudRepo.Delete(id)
}

// GetCredentialByID gets credential by ID
func (s *PCloudService) GetCredentialByID(id uuid.UUID) (*models.PCloudCredential, error) {
	return s.pcloudRepo.FindByID(id)
}

// GetAllCredentials gets all credentials
func (s *PCloudService) GetAllCredentials() ([]models.PCloudCredential, error) {
	return s.pcloudRepo.GetAll()
}

// GetActiveCredentials gets active credentials
func (s *PCloudService) GetActiveCredentials() ([]models.PCloudCredential, error) {
	return s.pcloudRepo.GetActive()
}

// ToggleActive toggles credential active status
func (s *PCloudService) ToggleActive(id uuid.UUID, isActive bool) error {
	return s.pcloudRepo.ToggleActive(id, isActive)
}