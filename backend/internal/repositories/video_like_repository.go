package repositories

import (
	"bobastream/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoLikeRepository struct {
	db *gorm.DB
}

func NewVideoLikeRepository(db *gorm.DB) *VideoLikeRepository {
	return &VideoLikeRepository{db: db}
}

// Create creates a new like
func (r *VideoLikeRepository) Create(like *models.VideoLike) error {
	return r.db.Create(like).Error
}

// Delete deletes a like
func (r *VideoLikeRepository) Delete(videoID, userID uuid.UUID) error {
	return r.db.Where("video_id = ? AND user_id = ?", videoID, userID).
		Delete(&models.VideoLike{}).Error
}

// IsLiked checks if user has liked the video
func (r *VideoLikeRepository) IsLiked(videoID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.VideoLike{}).
		Where("video_id = ? AND user_id = ?", videoID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetUserLikedVideos gets videos liked by user
func (r *VideoLikeRepository) GetUserLikedVideos(userID uuid.UUID, page, limit int) ([]models.Video, int64, error) {
	var videos []models.Video
	var total int64

	offset := (page - 1) * limit

	// Count total
	if err := r.db.Model(&models.VideoLike{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get liked videos
	err := r.db.Table("videos").
		Select("videos.*").
		Joins("JOIN video_likes ON video_likes.video_id = videos.id").
		Where("video_likes.user_id = ? AND videos.is_published = ?", userID, true).
		Preload("WrapperLink").
		Order("video_likes.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error

	return videos, total, err
}

// GetLikeCount gets like count for a video
func (r *VideoLikeRepository) GetLikeCount(videoID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.VideoLike{}).Where("video_id = ?", videoID).Count(&count).Error
	return count, err
}