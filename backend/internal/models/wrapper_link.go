package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WrapperLink struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VideoID      uuid.UUID `gorm:"type:uuid;not null;index" json:"video_id"`
	WrapperToken string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"wrapper_token"`
	Slug         string    `gorm:"type:varchar(500);uniqueIndex" json:"slug,omitempty"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Video *Video `gorm:"foreignKey:VideoID" json:"-"`
}

func (WrapperLink) TableName() string {
	return "wrapper_links"
}

func (w *WrapperLink) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	if w.WrapperToken == "" {
		w.WrapperToken = uuid.New().String()
	}
	return nil
}