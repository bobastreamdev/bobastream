package utils

import (
	"math"
	"time"

	"bobastream/internal/models"
)

// CalculateVideoScore calculates feed score for a video
func CalculateVideoScore(video *models.Video) float64 {
	now := time.Now()
	hoursSincePublish := 0.0
	
	if video.PublishedAt != nil {
		hoursSincePublish = now.Sub(*video.PublishedAt).Hours()
	}

	// Faktor 1: Recency (newer = better, decay setelah 7 hari)
	recencyScore := math.Exp(-hoursSincePublish/(7*24)) * 100

	// Faktor 2: Popularity (views + likes)
	popularityScore := math.Log10(float64(video.ViewCount+1)) * 20
	likeScore := math.Log10(float64(video.LikeCount+1)) * 30

	// Faktor 3: Engagement rate (likes/views)
	engagementRate := 0.0
	if video.ViewCount > 0 {
		engagementRate = float64(video.LikeCount) / float64(video.ViewCount) * 100
	}

	// Total score
	totalScore := recencyScore + popularityScore + likeScore + engagementRate

	return totalScore
}