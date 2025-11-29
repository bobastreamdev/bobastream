package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name         string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Slug         string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description  string         `gorm:"type:text" json:"description,omitempty"`
	Icon         string         `gorm:"type:varchar(50)" json:"icon,omitempty"`
	IsActive     bool           `gorm:"default:true;index" json:"is_active"`
	DisplayOrder int            `gorm:"default:0;index" json:"display_order"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Videos []Video `gorm:"foreignKey:CategoryID" json:"-"`
}

func (Category) TableName() string {
	return "categories"
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}