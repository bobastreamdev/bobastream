package cron

import (
	"bobastream/internal/services"
	"log"
	"time"
)

type AggregateStatsJob struct {
	analyticsService *services.AnalyticsService
	lock             *JobLock
}

func NewAggregateStatsJob(analyticsService *services.AnalyticsService) *AggregateStatsJob {
	return &AggregateStatsJob{
		analyticsService: analyticsService,
		lock:             NewJobLock(),
	}
}

// Run aggregates daily statistics
func (j *AggregateStatsJob) Run() {
	// ✅ Prevent overlapping runs
	if !j.lock.TryLock() {
		log.Println("⏭️  [CRON] Aggregate stats already running, skipping...")
		return
	}
	defer j.lock.Unlock()

	yesterday := time.Now().AddDate(0, 0, -1).Truncate(24 * time.Hour)
	
	log.Printf("📊 [CRON] Aggregating stats for %s...\n", yesterday.Format("2006-01-02"))

	if err := j.analyticsService.AggregateDailyStats(yesterday); err != nil {
		log.Printf("❌ [CRON] Failed to aggregate stats: %v\n", err)
		return
	}

	log.Println("✅ [CRON] Successfully aggregated daily stats")
}