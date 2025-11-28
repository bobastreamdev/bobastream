package repositories

import (
	"bobastream/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoViewRepository struct {
	db *gorm.DB
}

func NewVideoViewRepository(db *gorm.DB) *VideoViewRepository {
	return &VideoViewRepository{db: db}
}

// Create creates a new video view record
func (r *VideoViewRepository) Create(view *models.VideoView) error {
	return r.db.Create(view).Error
}

// Update updates video view record
func (r *VideoViewRepository) Update(view *models.VideoView) error {
	return r.db.Save(view).Error
}

// FindBySessionAndVideo finds view by session and video
func (r *VideoViewRepository) FindBySessionAndVideo(sessionID string, videoID uuid.UUID) (*models.VideoView, error) {
	var view models.VideoView
	err := r.db.Where("session_id = ? AND video_id = ?", sessionID, videoID).First(&view).Error
	if err != nil {
		return nil, err
	}
	return &view, nil
}

// GetValidViewCount gets count of valid views (watched >= 30%)
func (r *VideoViewRepository) GetValidViewCount(videoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.VideoView{}).
		Where("video_id = ? AND watched_percentage >= ?", videoID, 30).
		Distinct("session_id").
		Count(&count).Error
	return count, err
}

// GetViewsByDateRange gets views within date range
func (r *VideoViewRepository) GetViewsByDateRange(startDate, endDate time.Time) ([]models.VideoView, error) {
	var views []models.VideoView
	err := r.db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&views).Error
	return views, err
}

// GetDailyViewCount gets total views for a specific date
func (r *VideoViewRepository) GetDailyViewCount(date time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.VideoView{}).
		Where("date = ?", date).
		Count(&count).Error
	return count, err
}

// GetDailyUniqueViewers gets unique viewers for a specific date
func (r *VideoViewRepository) GetDailyUniqueViewers(date time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.VideoView{}).
		Where("date = ?", date).
		Distinct("COALESCE(user_id::text, viewer_ip)").
		Count(&count).Error
	return count, err
}

// GetAverageWatchTime gets average watch time for a date
func (r *VideoViewRepository) GetAverageWatchTime(date time.Time) (float64, error) {
	var avgTime float64
	err := r.db.Model(&models.VideoView{}).
		Select("COALESCE(AVG(watch_duration_seconds), 0)").
		Where("date = ?", date).
		Scan(&avgTime).Error
	return avgTime, err
}

// GetTopVideosByViews gets top videos by view count in date range
func (r *VideoViewRepository) GetTopVideosByViews(startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	err := r.db.Table("video_views").
		Select("video_id, COUNT(DISTINCT session_id) as view_count").
		Where("date BETWEEN ? AND ? AND watched_percentage >= ?", startDate, endDate, 30).
		Group("video_id").
		Order("view_count DESC").
		Limit(limit).
		Find(&results).Error

	return results, err
}