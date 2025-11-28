package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdType string

const (
	AdTypePreroll AdType = "preroll"
	AdTypeBanner  AdType = "banner"
	AdTypePopup   AdType = "popup"
)

type Ad struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Title            string         `gorm:"type:varchar(255);not null" json:"title"`
	AdType           AdType         `gorm:"type:ad_type;not null;index" json:"ad_type"`
	ContentURL       string         `gorm:"type:text" json:"content_url"`
	RedirectURL      string         `gorm:"type:text" json:"redirect_url"`
	DurationSeconds  int            `json:"duration_seconds"`
	DisplayFrequency int            `gorm:"default:1" json:"display_frequency"`
	IsActive         bool           `gorm:"default:true;index" json:"is_active"`
	Priority         int            `gorm:"default:0;index" json:"priority"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	AdImpressions []AdImpression `gorm:"foreignKey:AdID" json:"-"`
}

func (Ad) TableName() string {
	return "ads"
}

func (a *Ad) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}