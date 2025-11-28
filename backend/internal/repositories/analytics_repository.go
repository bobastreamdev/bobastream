package repositories

import (
	"bobastream/internal/models"
	"time"

	"gorm.io/gorm"
)

type AnalyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// CreateDailyStat creates daily stats record
func (r *AnalyticsRepository) CreateDailyStat(stat *models.DailyStat) error {
	return r.db.Create(stat).Error
}

// GetDailyStat gets daily stat by date
func (r *AnalyticsRepository) GetDailyStat(date time.Time) (*models.DailyStat, error) {
	var stat models.DailyStat
	err := r.db.Where("date = ?", date).First(&stat).Error
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

// GetStatsByDateRange gets stats within date range
func (r *AnalyticsRepository) GetStatsByDateRange(startDate, endDate time.Time) ([]models.DailyStat, error) {
	var stats []models.DailyStat
	err := r.db.Where("date BETWEEN ? AND ?", startDate, endDate).
		Order("date DESC").
		Find(&stats).Error
	return stats, err
}

// UpsertDailyStat creates or updates daily stat
func (r *AnalyticsRepository) UpsertDailyStat(stat *models.DailyStat) error {
	return r.db.Where("date = ?", stat.Date).
		Assign(stat).
		FirstOrCreate(stat).Error
}

// GetOverviewStats gets overview statistics
func (r *AnalyticsRepository) GetOverviewStats() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Total videos
	var totalVideos int64
	if err := r.db.Model(&models.Video{}).Where("is_published = ?", true).Count(&totalVideos).Error; err != nil {
		return nil, err
	}
	result["total_videos"] = totalVideos

	// Total users
	var totalUsers int64
	if err := r.db.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}
	result["total_users"] = totalUsers

	// Total views (all time)
	var totalViews int64
	if err := r.db.Model(&models.VideoView{}).Count(&totalViews).Error; err != nil {
		return nil, err
	}
	result["total_views"] = totalViews

	// Today's stats
	today := time.Now().Truncate(24 * time.Hour)
	todayStat, err := r.GetDailyStat(today)
	if err == nil {
		result["today_views"] = todayStat.TotalViews
		result["today_unique_viewers"] = todayStat.UniqueViewers
		result["today_ad_views"] = todayStat.TotalAdViews
	} else {
		result["today_views"] = 0
		result["today_unique_viewers"] = 0
		result["today_ad_views"] = 0
	}

	return result, nil
}

// GetMonthlyStats gets aggregated monthly statistics
func (r *AnalyticsRepository) GetMonthlyStats(year int, month int) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	stats, err := r.GetStatsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var totalViews, totalUniqueViewers, totalAdViews, totalAdClicks int64
	var totalWatchTime int64

	for _, stat := range stats {
		totalViews += int64(stat.TotalViews)
		totalUniqueViewers += int64(stat.UniqueViewers)
		totalAdViews += int64(stat.TotalAdViews)
		totalAdClicks += int64(stat.TotalAdClicks)
		totalWatchTime += int64(stat.AvgWatchTimeSeconds * stat.TotalViews)
	}

	avgWatchTime := 0
	if totalViews > 0 {
		avgWatchTime = int(totalWatchTime / totalViews)
	}

	result["total_views"] = totalViews
	result["unique_viewers"] = totalUniqueViewers
	result["total_ad_views"] = totalAdViews
	result["total_ad_clicks"] = totalAdClicks
	result["avg_watch_time_seconds"] = avgWatchTime
	result["days_count"] = len(stats)

	return result, nil
}