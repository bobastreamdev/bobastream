package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImpressionType string

const (
	ImpressionView  ImpressionType = "view"
	ImpressionClick ImpressionType = "click"
	ImpressionSkip  ImpressionType = "skip"
)

type AdImpression struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AdID            uuid.UUID      `gorm:"type:uuid;not null;index" json:"ad_id"`
	VideoID         *uuid.UUID     `gorm:"type:uuid;index" json:"video_id,omitempty"`
	UserID          *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	ViewerIP        string         `gorm:"type:varchar(45);not null" json:"viewer_ip"`
	ImpressionType  ImpressionType `gorm:"type:impression_type;not null;index" json:"impression_type"`
	WatchedDuration int            `gorm:"default:0" json:"watched_duration"`
	SessionID       string         `gorm:"type:varchar(255)" json:"session_id"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	Date            time.Time      `gorm:"type:date;default:CURRENT_DATE;index" json:"date"`

	// Relationships
	Ad    *Ad    `gorm:"foreignKey:AdID" json:"-"`
	Video *Video `gorm:"foreignKey:VideoID" json:"-"`
	User  *User  `gorm:"foreignKey:UserID" json:"-"`
}

func (AdImpression) TableName() string {
	return "ad_impressions"
}

func (a *AdImpression) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Date.IsZero() {
		a.Date = time.Now()
	}
	return nil
}