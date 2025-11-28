package repositories

import (
	"bobastream/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdImpressionRepository struct {
	db *gorm.DB
}

func NewAdImpressionRepository(db *gorm.DB) *AdImpressionRepository {
	return &AdImpressionRepository{db: db}
}

// Create creates a new ad impression
func (r *AdImpressionRepository) Create(impression *models.AdImpression) error {
	return r.db.Create(impression).Error
}

// GetDailyAdViews gets total ad views for a date
func (r *AdImpressionRepository) GetDailyAdViews(date time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.AdImpression{}).
		Where("date = ? AND impression_type = ?", date, models.ImpressionView).
		Count(&count).Error
	return count, err
}

// GetDailyAdClicks gets total ad clicks for a date
func (r *AdImpressionRepository) GetDailyAdClicks(date time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.AdImpression{}).
		Where("date = ? AND impression_type = ?", date, models.ImpressionClick).
		Count(&count).Error
	return count, err
}

// GetAdPerformance gets ad performance by date range
func (r *AdImpressionRepository) GetAdPerformance(adID uuid.UUID, startDate, endDate time.Time) (map[string]int64, error) {
	result := make(map[string]int64)

	// Get views
	var views int64
	if err := r.db.Model(&models.AdImpression{}).
		Where("ad_id = ? AND date BETWEEN ? AND ? AND impression_type = ?", 
			adID, startDate, endDate, models.ImpressionView).
		Count(&views).Error; err != nil {
		return nil, err
	}
	result["views"] = views

	// Get clicks
	var clicks int64
	if err := r.db.Model(&models.AdImpression{}).
		Where("ad_id = ? AND date BETWEEN ? AND ? AND impression_type = ?", 
			adID, startDate, endDate, models.ImpressionClick).
		Count(&clicks).Error; err != nil {
		return nil, err
	}
	result["clicks"] = clicks

	// Get skips
	var skips int64
	if err := r.db.Model(&models.AdImpression{}).
		Where("ad_id = ? AND date BETWEEN ? AND ? AND impression_type = ?", 
			adID, startDate, endDate, models.ImpressionSkip).
		Count(&skips).Error; err != nil {
		return nil, err
	}
	result["skips"] = skips

	return result, nil
}

// GetTopPerformingAds gets top performing ads by clicks
func (r *AdImpressionRepository) GetTopPerformingAds(startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.db.Table("ad_impressions").
		Select("ad_id, COUNT(*) FILTER (WHERE impression_type = 'view') as views, COUNT(*) FILTER (WHERE impression_type = 'click') as clicks").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Group("ad_id").
		Order("clicks DESC").
		Limit(limit).
		Find(&results).Error

	return results, err
}

// GetImpressionsByDateRange gets impressions within date range
func (r *AdImpressionRepository) GetImpressionsByDateRange(startDate, endDate time.Time) ([]models.AdImpression, error) {
	var impressions []models.AdImpression
	err := r.db.Where("date BETWEEN ? AND ?", startDate, endDate).
		Find(&impressions).Error
	return impressions, err
}