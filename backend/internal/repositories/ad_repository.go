package repositories

import (
	"bobastream/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdRepository struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) *AdRepository {
	return &AdRepository{db: db}
}

// Create creates a new ad
func (r *AdRepository) Create(ad *models.Ad) error {
	return r.db.Create(ad).Error
}

// FindByID finds ad by ID
func (r *AdRepository) FindByID(id uuid.UUID) (*models.Ad, error) {
	var ad models.Ad
	err := r.db.First(&ad, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &ad, nil
}

// Update updates ad
func (r *AdRepository) Update(ad *models.Ad) error {
	return r.db.Save(ad).Error
}

// Delete soft deletes ad
func (r *AdRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Ad{}, "id = ?", id).Error
}

// GetAll gets all ads with pagination
func (r *AdRepository) GetAll(page, limit int) ([]models.Ad, int64, error) {
	var ads []models.Ad
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&models.Ad{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Order("priority DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&ads).Error

	return ads, total, err
}

// GetActiveAdByType gets active ad by type with highest priority
func (r *AdRepository) GetActiveAdByType(adType models.AdType) (*models.Ad, error) {
	var ad models.Ad
	err := r.db.Where("ad_type = ? AND is_active = ?", adType, true).
		Order("priority DESC").
		First(&ad).Error
	if err != nil {
		return nil, err
	}
	return &ad, nil
}

// GetActiveAds gets all active ads
func (r *AdRepository) GetActiveAds() ([]models.Ad, error) {
	var ads []models.Ad
	err := r.db.Where("is_active = ?", true).
		Order("ad_type, priority DESC").
		Find(&ads).Error
	return ads, err
}

// GetAdsByType gets ads by type
func (r *AdRepository) GetAdsByType(adType models.AdType) ([]models.Ad, error) {
	var ads []models.Ad
	err := r.db.Where("ad_type = ?", adType).
		Order("is_active DESC, priority DESC").
		Find(&ads).Error
	return ads, err
}

// ToggleActive toggles ad active status
func (r *AdRepository) ToggleActive(id uuid.UUID, isActive bool) error {
	return r.db.Model(&models.Ad{}).Where("id = ?", id).
		Update("is_active", isActive).Error
}