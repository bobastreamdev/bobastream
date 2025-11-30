package cron

import (
	"bobastream/internal/services"
	"log"
)

type RefreshLinksJob struct {
	pcloudService *services.PCloudService
	lock          *JobLock
}

func NewRefreshLinksJob(pcloudService *services.PCloudService) *RefreshLinksJob {
	return &RefreshLinksJob{
		pcloudService: pcloudService,
		lock:          NewJobLock(),
	}
}

// Run refreshes expired pCloud video links
func (j *RefreshLinksJob) Run() {
	// ✅ Prevent overlapping runs
	if !j.lock.TryLock() {
		log.Println("⏭️  [CRON] Refresh links already running, skipping...")
		return
	}
	defer j.lock.Unlock()

	log.Println("🔄 [CRON] Starting refresh expired pCloud links...")

	if err := j.pcloudService.RefreshExpiredLinks(); err != nil {
		log.Printf("❌ [CRON] Failed to refresh links: %v\n", err)
		return
	}

	log.Println("✅ [CRON] Successfully refreshed expired pCloud links")
}