package services

import (
	"bobastream/internal/cache"
	"bobastream/internal/models"
	"bobastream/internal/repositories"
	"bobastream/internal/utils"
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoService struct {
	videoRepo     *repositories.VideoRepository
	wrapperRepo   *repositories.WrapperLinkRepository
	videoViewRepo *repositories.VideoViewRepository
	videoLikeRepo *repositories.VideoLikeRepository
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

// ScoredVideo is a helper struct for sorting videos by score
type ScoredVideo struct {
	Video models.Video
	Score float64
}

// ✅ CACHED: GetFeedVideos with Redis caching
func (s *VideoService) GetFeedVideos(page, limit int) ([]models.Video, int64, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("feed:page:%d:limit:%d", page, limit)
	
	// ✅ Try cache first
	var cachedResult struct {
		Videos []models.Video
		Total  int64
	}
	
	if err := cache.Get(ctx, cacheKey, &cachedResult); err == nil {
		return cachedResult.Videos, cachedResult.Total, nil
	}
	
	// ✅ Cache miss - fetch from DB
	videos, total, err := s.videoRepo.GetPublishedVideos(page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Score videos
	scoredVideos := make([]ScoredVideo, len(videos))
	for i := range videos {
		scoredVideos[i] = ScoredVideo{
			Video: videos[i],
			Score: utils.CalculateVideoScore(&videos[i]),
		}
	}

	// Sort by score (highest first)
	sort.Slice(scoredVideos, func(i, j int) bool {
		return scoredVideos[i].Score > scoredVideos[j].Score
	})

	// Extract sorted videos
	result := make([]models.Video, len(scoredVideos))
	for i, sv := range scoredVideos {
		result[i] = sv.Video
	}

	// ✅ Cache result for 5 minutes
	cache.Set(ctx, cacheKey, map[string]interface{}{
		"Videos": result,
		"Total":  total,
	}, 5*time.Minute)

	return result, total, nil
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

// TrackVideoView tracks video view with proper view count increment logic
func (s *VideoService) TrackVideoView(videoID uuid.UUID, userID *uuid.UUID, sessionID, viewerIP, userAgent string, watchDuration int, videoDuration int) error {
	// Calculate watched percentage
	watchedPercentage := 0.0
	if videoDuration > 0 {
		watchedPercentage = (float64(watchDuration) / float64(videoDuration)) * 100
		if watchedPercentage > 100 {
			watchedPercentage = 100
		}
	}

	// Check if view already exists for this session
	existingView, err := s.videoViewRepo.FindBySessionAndVideo(sessionID, videoID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	shouldIncrementView := false

	if existingView != nil {
		// View exists - check if crossing 30% threshold
		wasInvalid := existingView.WatchedPercentage < 30
		isNowValid := watchedPercentage >= 30

		// Update existing view record
		existingView.WatchDurationSeconds = watchDuration
		existingView.WatchedPercentage = watchedPercentage

		if err := s.videoViewRepo.Update(existingView); err != nil {
			return err
		}

		// Increment ONLY when crossing 30% threshold (first time)
		if wasInvalid && isNowValid {
			shouldIncrementView = true
		}
	} else {
		// New view - create record
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

		// Increment if first view is already >= 30%
		if watchedPercentage >= 30 {
			shouldIncrementView = true
		}
	}

	// Increment view count if conditions met
	if shouldIncrementView {
		// ✅ Invalidate feed cache when view count changes
		ctx := context.Background()
		cache.Delete(ctx, "feed:*") // Delete all feed cache keys
		
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

	// ✅ Invalidate feed cache
	ctx := context.Background()
	cache.Delete(ctx, "feed:*")

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

	// ✅ Invalidate feed cache
	ctx := context.Background()
	cache.Delete(ctx, "feed:*")

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
	// ✅ Invalidate feed cache
	ctx := context.Background()
	cache.Delete(ctx, "feed:*")
	
	return s.videoRepo.Create(video)
}

// UpdateVideo updates video (admin)
func (s *VideoService) UpdateVideo(video *models.Video) error {
	// ✅ Invalidate feed cache
	ctx := context.Background()
	cache.Delete(ctx, "feed:*")
	
	return s.videoRepo.Update(video)
}

// DeleteVideo deletes video (admin)
func (s *VideoService) DeleteVideo(id uuid.UUID) error {
	// ✅ Invalidate feed cache
	ctx := context.Background()
	cache.Delete(ctx, "feed:*")
	
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