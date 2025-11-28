package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PCloudCredential struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AccountName     string         `gorm:"type:varchar(100);not null" json:"account_name"`
	APIToken        string         `gorm:"type:text;not null" json:"-"` // Hidden
	AccessToken     string         `gorm:"type:text" json:"-"`          // Hidden
	RefreshToken    string         `gorm:"type:text" json:"-"`          // Hidden
	TokenExpiresAt  *time.Time     `json:"token_expires_at,omitempty"`
	StorageUsedGB   float64        `gorm:"type:decimal(10,2);default:0" json:"storage_used_gb"`
	StorageLimitGB  float64        `gorm:"type:decimal(10,2);not null" json:"storage_limit_gb"`
	IsActive        bool           `gorm:"default:true;index" json:"is_active"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Videos []Video `gorm:"foreignKey:PCloudCredentialID" json:"-"`
}

func (PCloudCredential) TableName() string {
	return "pcloud_credentials"
}

func (p *PCloudCredential) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}