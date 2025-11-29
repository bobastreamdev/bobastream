package repositories

import (
	"bobastream/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

// Create creates a new video
func (r *VideoRepository) Create(video *models.Video) error {
	return r.db.Create(video).Error
}

// FindByID finds video by ID with wrapper link
func (r *VideoRepository) FindByID(id uuid.UUID) (*models.Video, error) {
	var video models.Video
	err := r.db.Preload("WrapperLink").First(&video, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// FindByWrapperToken finds video by wrapper token
func (r *VideoRepository) FindByWrapperToken(token string) (*models.Video, error) {
	var video models.Video
	err := r.db.Joins("JOIN wrapper_links ON wrapper_links.video_id = videos.id").
		Where("wrapper_links.wrapper_token = ?", token).
		Preload("WrapperLink").
		First(&video).Error
	if err != nil {
		return nil, err
	}
	return &video, nil
}

// Update updates video
func (r *VideoRepository) Update(video *models.Video) error {
	return r.db.Save(video).Error
}

// Delete soft deletes video
func (r *VideoRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Video{}, "id = ?", id).Error
}

// GetPublishedVideos gets published videos with scoring for feed
func (r *VideoRepository) GetPublishedVideos(page, limit int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	offset := (page - 1) * limit

	// Count total published videos
	if err := r.db.Model(&models.Video{}).Where("is_published = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get videos with scoring (simplified query, actual scoring done in service layer)
	err := r.db.Where("is_published = ?", true).
		Preload("WrapperLink").
		Order("published_at DESC, view_count DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// GetAllVideos gets all videos (admin)
func (r *VideoRepository) GetAllVideos(page, limit int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&models.Video{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("WrapperLink").
		Preload("PCloudCredential").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// GetRelatedVideos gets related videos based on category and tags
func (r *VideoRepository) GetRelatedVideos(videoID uuid.UUID, categoryID string, tags []string, limit int) ([]models.Video, error) {
	var videos []models.Video

	query := r.db.Where("id != ? AND is_published = ?", videoID, true)

	// Prioritize same category
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	err := query.Preload("WrapperLink").
		Preload("Category").
		Order("view_count DESC").
		Limit(limit).
		Find(&videos).Error

	return videos, err
}

// SearchVideos searches videos by title or description
func (r *VideoRepository) SearchVideos(keyword string, page, limit int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	offset := (page - 1) * limit

	query := r.db.Where("is_published = ? AND (title ILIKE ? OR description ILIKE ?)", 
		true, "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Model(&models.Video{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("WrapperLink").
		Preload("Category").
		Order("view_count DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// SearchVideosByCategory searches videos within a category
func (r *VideoRepository) SearchVideosByCategory(keyword string, categoryID uuid.UUID, page, limit int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	offset := (page - 1) * limit

	query := r.db.Where("is_published = ? AND category_id = ? AND (title ILIKE ? OR description ILIKE ?)", 
		true, categoryID, "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Model(&models.Video{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("WrapperLink").
		Preload("Category").
		Order("view_count DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// GetVideosByCategory gets videos by category
func (r *VideoRepository) GetVideosByCategory(categoryID uuid.UUID, page, limit int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	offset := (page - 1) * limit

	query := r.db.Where("is_published = ? AND category_id = ?", true, categoryID)

	if err := query.Model(&models.Video{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("WrapperLink").
		Preload("Category").
		Order("published_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// IncrementViewCount increments video view count
func (r *VideoRepository) IncrementViewCount(id uuid.UUID) error {
	return r.db.Model(&models.Video{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// IncrementLikeCount increments video like count
func (r *VideoRepository) IncrementLikeCount(id uuid.UUID) error {
	return r.db.Model(&models.Video{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

// DecrementLikeCount decrements video like count
func (r *VideoRepository) DecrementLikeCount(id uuid.UUID) error {
	return r.db.Model(&models.Video{}).Where("id = ?", id).
		UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)")).Error
}

// GetExpiredSourceURLVideos gets videos with expired source URLs
func (r *VideoRepository) GetExpiredSourceURLVideos() ([]models.Video, error) {
	var videos []models.Video
	err := r.db.Where("source_url_expires_at < ?", time.Now().Add(1*time.Hour)).
		Preload("PCloudCredential").
		Find(&videos).Error
	return videos, err
}

// UpdateSourceURL updates video source URL and expiry
func (r *VideoRepository) UpdateSourceURL(id uuid.UUID, sourceURL string, expiresAt time.Time) error {
	return r.db.Model(&models.Video{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"source_url":           sourceURL,
			"source_url_expires_at": expiresAt,
		}).Error
}

// GetTopVideos gets top videos by view count
func (r *VideoRepository) GetTopVideos(limit int, days int) ([]models.Video, error) {
	var videos []models.Video
	
	cutoffDate := time.Now().AddDate(0, 0, -days)
	
	err := r.db.Where("is_published = ? AND published_at >= ?", true, cutoffDate).
		Preload("WrapperLink").
		Order("view_count DESC, like_count DESC").
		Limit(limit).
		Find(&videos).Error

	return videos, err
}