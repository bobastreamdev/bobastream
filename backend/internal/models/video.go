package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Video struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Title               string         `gorm:"type:varchar(500);not null" json:"title"`
	Description         string         `gorm:"type:text" json:"description"`
	ThumbnailURL        string         `gorm:"type:text" json:"thumbnail_url"`
	SourceURL           string         `gorm:"type:text;not null" json:"-"` // Hidden from API
	SourceURLExpiresAt  *time.Time     `json:"-"`
	DurationSeconds     int            `json:"duration_seconds"`
	FileSizeMB          float64        `gorm:"type:decimal(10,2)" json:"file_size_mb"`
	PCloudFileID        string         `gorm:"type:varchar(255)" json:"-"`
	PCloudCredentialID  uuid.UUID      `gorm:"type:uuid;not null" json:"-"`
	Genre               string         `gorm:"type:varchar(100)" json:"genre"`
	Tags                pq.StringArray `gorm:"type:text[]" json:"tags"`
	ViewCount           int            `gorm:"default:0" json:"view_count"`
	LikeCount           int            `gorm:"default:0" json:"like_count"`
	IsPublished         bool           `gorm:"default:true" json:"is_published"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	PublishedAt         *time.Time     `json:"published_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	PCloudCredential *PCloudCredential `gorm:"foreignKey:PCloudCredentialID" json:"-"`
	WrapperLink      *WrapperLink      `gorm:"foreignKey:VideoID" json:"wrapper_link,omitempty"`
	VideoLikes       []VideoLike       `gorm:"foreignKey:VideoID" json:"-"`
	VideoViews       []VideoView       `gorm:"foreignKey:VideoID" json:"-"`
}

func (Video) TableName() string {
	return "videos"
}

func (v *Video) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	if v.IsPublished && v.PublishedAt == nil {
		now := time.Now()
		v.PublishedAt = &now
	}
	return nil
}