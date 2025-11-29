package cron

import (
	"bobastream/internal/services"
	"log"
)

type RefreshLinksJob struct {
	pcloudService *services.PCloudService
}

func NewRefreshLinksJob(pcloudService *services.PCloudService) *RefreshLinksJob {
	return &RefreshLinksJob{pcloudService: pcloudService}
}

// Run refreshes expired pCloud video links
func (j *RefreshLinksJob) Run() {
	log.Println("🔄 [CRON] Starting refresh expired pCloud links...")

	if err := j.pcloudService.RefreshExpiredLinks(); err != nil {
		log.Printf("❌ [CRON] Failed to refresh links: %v\n", err)
		return
	}

	log.Println("✅ [CRON] Successfully refreshed expired pCloud links")
}