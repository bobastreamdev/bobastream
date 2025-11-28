package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoLike struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	VideoID   uuid.UUID `gorm:"type:uuid;not null;index" json:"video_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`

	// Relationships
	Video *Video `gorm:"foreignKey:VideoID" json:"-"`
	User  *User  `gorm:"foreignKey:UserID" json:"-"`
}

func (VideoLike) TableName() string {
	return "video_likes"
}

func (v *VideoLike) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}