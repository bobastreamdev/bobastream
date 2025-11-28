package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DailyStat struct {
	ID                   uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Date                 time.Time `gorm:"type:date;uniqueIndex;not null" json:"date"`
	TotalViews           int       `gorm:"default:0" json:"total_views"`
	UniqueViewers        int       `gorm:"default:0" json:"unique_viewers"`
	TotalAdViews         int       `gorm:"default:0" json:"total_ad_views"`
	TotalAdClicks        int       `gorm:"default:0" json:"total_ad_clicks"`
	AvgWatchTimeSeconds  int       `gorm:"default:0" json:"avg_watch_time_seconds"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (DailyStat) TableName() string {
	return "daily_stats"
}

func (d *DailyStat) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}