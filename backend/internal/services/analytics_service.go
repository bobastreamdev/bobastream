package services

import (
	"bobastream/internal/models"
	"bobastream/internal/repositories"
	"time"
)

type AnalyticsService struct {
	analyticsRepo    *repositories.AnalyticsRepository
	videoViewRepo    *repositories.VideoViewRepository
	adImpressionRepo *repositories.AdImpressionRepository
}

func NewAnalyticsService(
	analyticsRepo *repositories.AnalyticsRepository,
	videoViewRepo *repositories.VideoViewRepository,
	adImpressionRepo *repositories.AdImpressionRepository,
) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo:    analyticsRepo,
		videoViewRepo:    videoViewRepo,
		adImpressionRepo: adImpressionRepo,
	}
}

// GetOverviewStats gets dashboard overview statistics
func (s *AnalyticsService) GetOverviewStats() (map[string]interface{}, error) {
	return s.analyticsRepo.GetOverviewStats()
}

// GetDailyStats gets stats for a specific date
func (s *AnalyticsService) GetDailyStats(date time.Time) (*models.DailyStat, error) {
	return s.analyticsRepo.GetDailyStat(date)
}

// GetStatsByDateRange gets stats within date range
func (s *AnalyticsService) GetStatsByDateRange(startDate, endDate time.Time) ([]models.DailyStat, error) {
	return s.analyticsRepo.GetStatsByDateRange(startDate, endDate)
}

// GetMonthlyStats gets aggregated monthly stats
func (s *AnalyticsService) GetMonthlyStats(year, month int) (map[string]interface{}, error) {
	return s.analyticsRepo.GetMonthlyStats(year, month)
}

// AggregateDailyStats aggregates stats for a specific date (called by cron)
func (s *AnalyticsService) AggregateDailyStats(date time.Time) error {
	// Get total views
	totalViews, err := s.videoViewRepo.GetDailyViewCount(date)
	if err != nil {
		return err
	}

	// Get unique viewers
	uniqueViewers, err := s.videoViewRepo.GetDailyUniqueViewers(date)
	if err != nil {
		return err
	}

	// Get average watch time
	avgWatchTime, err := s.videoViewRepo.GetAverageWatchTime(date)
	if err != nil {
		return err
	}

	// Get ad views
	totalAdViews, err := s.adImpressionRepo.GetDailyAdViews(date)
	if err != nil {
		return err
	}

	// Get ad clicks
	totalAdClicks, err := s.adImpressionRepo.GetDailyAdClicks(date)
	if err != nil {
		return err
	}

	// Create or update daily stat
	stat := &models.DailyStat{
		Date:                date,
		TotalViews:          int(totalViews),
		UniqueViewers:       int(uniqueViewers),
		TotalAdViews:        int(totalAdViews),
		TotalAdClicks:       int(totalAdClicks),
		AvgWatchTimeSeconds: int(avgWatchTime),
	}

	return s.analyticsRepo.UpsertDailyStat(stat)
}

// GetTopVideosByViews gets top videos by view count
func (s *AnalyticsService) GetTopVideosByViews(startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	return s.videoViewRepo.GetTopVideosByViews(startDate, endDate, limit)
}