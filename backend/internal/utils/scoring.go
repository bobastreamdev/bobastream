package utils

import (
	"math"
	"time"

	"bobastream/internal/models"
)

// ✅ FIXED: CalculateVideoScore with bounds check and edge case handling
func CalculateVideoScore(video *models.Video) float64 {
	now := time.Now()
	hoursSincePublish := 0.0
	
	// ✅ Check PublishedAt is not nil AND is before now (prevent future dates)
	if video.PublishedAt != nil && video.PublishedAt.Before(now) {
		hoursSincePublish = now.Sub(*video.PublishedAt).Hours()
		
		// ✅ CAP maximum decay at 30 days (720 hours)
		if hoursSincePublish > 720 {
			hoursSincePublish = 720
		}
	}

	// Faktor 1: Recency (newer = better, decay setelah 7 hari, max 30 hari)
	recencyScore := math.Exp(-hoursSincePublish/(7*24)) * 100

	// Faktor 2: Popularity (log scale prevents extreme values)
	popularityScore := math.Log10(float64(video.ViewCount+1)) * 20
	likeScore := math.Log10(float64(video.LikeCount+1)) * 30

	// Faktor 3: Engagement rate (safe division)
	engagementRate := 0.0
	if video.ViewCount > 0 {
		engagementRate = float64(video.LikeCount) / float64(video.ViewCount) * 100
		
		// ✅ CAP engagement at 100% (handle edge cases)
		if engagementRate > 100 {
			engagementRate = 100
		}
	}

	// Total score
	totalScore := recencyScore + popularityScore + likeScore + engagementRate

	// ✅ Ensure score is not NaN or Inf
	if math.IsNaN(totalScore) || math.IsInf(totalScore, 0) {
		return 0
	}

	return totalScore
}