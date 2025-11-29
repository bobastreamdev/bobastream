package services

import (
	"bobastream/internal/models"
	"bobastream/internal/repositories"

	"github.com/google/uuid"
)

type AdService struct {
	adRepo           *repositories.AdRepository
	adImpressionRepo *repositories.AdImpressionRepository
}

func NewAdService(
	adRepo *repositories.AdRepository,
	adImpressionRepo *repositories.AdImpressionRepository,
) *AdService {
	return &AdService{
		adRepo:           adRepo,
		adImpressionRepo: adImpressionRepo,
	}
}

// GetActiveAdByType gets active ad by type
func (s *AdService) GetActiveAdByType(adType models.AdType) (*models.Ad, error) {
	return s.adRepo.GetActiveAdByType(adType)
}

// TrackAdImpression tracks ad impression (view/click/skip)
func (s *AdService) TrackAdImpression(
	adID uuid.UUID,
	videoID *uuid.UUID,
	userID *uuid.UUID,
	viewerIP string,
	impressionType models.ImpressionType,
	watchedDuration int,
	sessionID string,
) error {
	impression := &models.AdImpression{
		AdID:            adID,
		VideoID:         videoID,
		UserID:          userID,
		ViewerIP:        viewerIP,
		ImpressionType:  impressionType,
		WatchedDuration: watchedDuration,
		SessionID:       sessionID,
	}

	return s.adImpressionRepo.Create(impression)
}

// CreateAd creates a new ad (admin)
func (s *AdService) CreateAd(ad *models.Ad) error {
	return s.adRepo.Create(ad)
}

// UpdateAd updates ad (admin)
func (s *AdService) UpdateAd(ad *models.Ad) error {
	return s.adRepo.Update(ad)
}

// DeleteAd deletes ad (admin)
func (s *AdService) DeleteAd(id uuid.UUID) error {
	return s.adRepo.Delete(id)
}

// GetAdByID gets ad by ID
func (s *AdService) GetAdByID(id uuid.UUID) (*models.Ad, error) {
	return s.adRepo.FindByID(id)
}

// GetAllAds gets all ads with pagination (admin)
func (s *AdService) GetAllAds(page, limit int) ([]models.Ad, int64, error) {
	return s.adRepo.GetAll(page, limit)
}

// GetActiveAds gets all active ads
func (s *AdService) GetActiveAds() ([]models.Ad, error) {
	return s.adRepo.GetActiveAds()
}

// GetAdsByType gets ads by type
func (s *AdService) GetAdsByType(adType models.AdType) ([]models.Ad, error) {
	return s.adRepo.GetAdsByType(adType)
}

// ToggleActive toggles ad active status
func (s *AdService) ToggleActive(id uuid.UUID, isActive bool) error {
	return s.adRepo.ToggleActive(id, isActive)
}