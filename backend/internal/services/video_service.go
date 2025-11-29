package services

import (
	"bobastream/internal/models"
	"bobastream/internal/repositories"
	"bobastream/internal/utils"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoService struct {
	videoRepo       *repositories.VideoRepository
	wrapperRepo     *repositories.WrapperLinkRepository
	videoViewRepo   *repositories.VideoViewRepository
	videoLikeRepo   *repositories.VideoLikeRepository
}

func NewVideoService(
	videoRepo *repositories.VideoRepository,
	wrapperRepo *repositories.WrapperLinkRepository,
	videoViewRepo *repositories.VideoViewRepository,
	videoLikeRepo *repositories.VideoLikeRepository,
) *VideoService {
	return &VideoService{
		videoRepo:     videoRepo,
		wrapperRepo:   wrapperRepo,
		videoViewRepo: videoViewRepo,
		videoLikeRepo: videoLikeRepo,
	}
}

// GetFeedVideos gets feed videos with scoring
func (s *VideoService) GetFeedVideos(page, limit int) ([]models.Video, int64, error) {
	videos, total, err := s.videoRepo.GetPublishedVideos(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Apply scoring and sort
	for i := range videos {
		videos[i].ViewCount = int(utils.CalculateVideoScore(&videos[i]))
	}

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].ViewCount > videos[j].ViewCount
	})

	return videos, total, nil
}

// GetVideoByID gets video by ID
func (s *VideoService) GetVideoByID(id uuid.UUID) (*models.Video, error) {
	return s.videoRepo.FindByID(id)
}

// GetVideoByWrapperToken gets video by wrapper token
func (s *VideoService) GetVideoByWrapperToken(token string) (*models.Video, error) {
	return s.videoRepo.FindByWrapperToken(token)
}

// SearchVideos searches videos with optional category filter
func (s *VideoService) SearchVideos(keyword string, categoryID *uuid.UUID, page, limit int) ([]models.Video, int64, error) {
	// If category filter is set, search within category
	if categoryID != nil {
		return s.videoRepo.SearchVideosByCategory(keyword, *categoryID, page, limit)
	}
	return s.videoRepo.SearchVideos(keyword, page, limit)
}

// GetVideosByCategory gets videos by category
func (s *VideoService) GetVideosByCategory(categoryID uuid.UUID, page, limit int) ([]models.Video, int64, error) {
	return s.videoRepo.GetVideosByCategory(categoryID, page, limit)
}

// GetRelatedVideos gets related videos based on current video
func (s *VideoService) GetRelatedVideos(videoID uuid.UUID, limit int) ([]models.Video, error) {
	// Get current video
	video, err := s.videoRepo.FindByID(videoID)
	if err != nil {
		return nil, err
	}

	categoryID := ""
	if video.CategoryID != nil {
		categoryID = video.CategoryID.String()
	}

	tags := []string{}
	if video.Tags != nil {
		tags = video.Tags
	}

	return s.videoRepo.GetRelatedVideos(videoID, categoryID, tags, limit)
}

// TrackVideoView tracks video view with watch duration
func (s *VideoService) TrackVideoView(videoID uuid.UUID, userID *uuid.UUID, sessionID, viewerIP, userAgent string, watchDuration int, videoDuration int) error {
	// Calculate watched percentage
	watchedPercentage := 0.0
	if videoDuration > 0 {
		watchedPercentage = (float64(watchDuration) / float64(videoDuration)) * 100
		if watchedPercentage > 100 {
			watchedPercentage = 100
		}
	}

	// Check if view already exists
	existingView, err := s.videoViewRepo.FindBySessionAndVideo(sessionID, videoID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existingView != nil {
		// Update existing view
		existingView.WatchDurationSeconds = watchDuration
		existingView.WatchedPercentage = watchedPercentage
		return s.videoViewRepo.Update(existingView)
	}

	// Create new view
	view := &models.VideoView{
		VideoID:              videoID,
		UserID:               userID,
		ViewerIP:             viewerIP,
		UserAgent:            userAgent,
		WatchDurationSeconds: watchDuration,
		WatchedPercentage:    watchedPercentage,
		SessionID:            sessionID,
	}

	if err := s.videoViewRepo.Create(view); err != nil {
		return err
	}

	// Increment view count if watched >= 30%
	if watchedPercentage >= 30 {
		return s.videoRepo.IncrementViewCount(videoID)
	}

	return nil
}

// LikeVideo likes a video
func (s *VideoService) LikeVideo(videoID, userID uuid.UUID) error {
	// Check if already liked
	isLiked, err := s.videoLikeRepo.IsLiked(videoID, userID)
	if err != nil {
		return err
	}
	if isLiked {
		return errors.New("already liked")
	}

	// Create like
	like := &models.VideoLike{
		VideoID: videoID,
		UserID:  userID,
	}

	if err := s.videoLikeRepo.Create(like); err != nil {
		return err
	}

	// Increment like count
	return s.videoRepo.IncrementLikeCount(videoID)
}

// UnlikeVideo unlikes a video
func (s *VideoService) UnlikeVideo(videoID, userID uuid.UUID) error {
	// Check if liked
	isLiked, err := s.videoLikeRepo.IsLiked(videoID, userID)
	if err != nil {
		return err
	}
	if !isLiked {
		return errors.New("not liked")
	}

	// Delete like
	if err := s.videoLikeRepo.Delete(videoID, userID); err != nil {
		return err
	}

	// Decrement like count
	return s.videoRepo.DecrementLikeCount(videoID)
}

// IsVideoLiked checks if user has liked video
func (s *VideoService) IsVideoLiked(videoID, userID uuid.UUID) (bool, error) {
	return s.videoLikeRepo.IsLiked(videoID, userID)
}

// GetUserLikedVideos gets user's liked videos
func (s *VideoService) GetUserLikedVideos(userID uuid.UUID, page, limit int) ([]models.Video, int64, error) {
	return s.videoLikeRepo.GetUserLikedVideos(userID, page, limit)
}

// GetTopVideos gets top videos by view count
func (s *VideoService) GetTopVideos(limit, days int) ([]models.Video, error) {
	return s.videoRepo.GetTopVideos(limit, days)
}

// CreateVideo creates a new video (admin)
func (s *VideoService) CreateVideo(video *models.Video) error {
	return s.videoRepo.Create(video)
}

// UpdateVideo updates video (admin)
func (s *VideoService) UpdateVideo(video *models.Video) error {
	return s.videoRepo.Update(video)
}

// DeleteVideo deletes video (admin)
func (s *VideoService) DeleteVideo(id uuid.UUID) error {
	return s.videoRepo.Delete(id)
}

// CreateWrapperLink creates wrapper link for video
func (s *VideoService) CreateWrapperLink(link *models.WrapperLink) error {
	return s.wrapperRepo.Create(link)
}

// GetAllVideos gets all videos (admin)
func (s *VideoService) GetAllVideos(page, limit int) ([]models.Video, int64, error) {
	return s.videoRepo.GetAllVideos(page, limit)
}

// UpdateSourceURL updates video source URL
func (s *VideoService) UpdateSourceURL(id uuid.UUID, sourceURL string, expiresAt time.Time) error {
	return s.videoRepo.UpdateSourceURL(id, sourceURL, expiresAt)
}