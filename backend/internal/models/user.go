package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleViewer UserRole = "viewer"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Username     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Role         UserRole       `gorm:"type:user_role;default:'viewer'" json:"role"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	LastLogin    *time.Time     `json:"last_login,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Relationships
	VideoLikes []VideoLike `gorm:"foreignKey:UserID" json:"-"`
	VideoViews []VideoView `gorm:"foreignKey:UserID" json:"-"`
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate hook for UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}