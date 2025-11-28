package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoView struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VideoID              uuid.UUID  `gorm:"type:uuid;not null;index" json:"video_id"`
	UserID               *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	ViewerIP             string     `gorm:"type:varchar(45);not null" json:"viewer_ip"`
	UserAgent            string     `gorm:"type:text" json:"user_agent,omitempty"`
	WatchDurationSeconds int        `gorm:"default:0" json:"watch_duration_seconds"`
	WatchedPercentage    float64    `gorm:"type:decimal(5,2);default:0" json:"watched_percentage"`
	SessionID            string     `gorm:"type:varchar(255);not null;index" json:"session_id"`
	ViewedAt             time.Time  `gorm:"autoCreateTime" json:"viewed_at"`
	Date                 time.Time  `gorm:"type:date;default:CURRENT_DATE;index" json:"date"`

	// Relationships
	Video *Video `gorm:"foreignKey:VideoID" json:"-"`
	User  *User  `gorm:"foreignKey:UserID" json:"-"`
}

func (VideoView) TableName() string {
	return "video_views"
}

func (v *VideoView) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	if v.Date.IsZero() {
		v.Date = time.Now()
	}
	return nil
}