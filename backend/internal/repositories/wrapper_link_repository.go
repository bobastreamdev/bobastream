package repositories

import (
	"bobastream/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WrapperLinkRepository struct {
	db *gorm.DB
}

func NewWrapperLinkRepository(db *gorm.DB) *WrapperLinkRepository {
	return &WrapperLinkRepository{db: db}
}

// Create creates a new wrapper link
func (r *WrapperLinkRepository) Create(link *models.WrapperLink) error {
	return r.db.Create(link).Error
}

// FindByToken finds wrapper link by token
func (r *WrapperLinkRepository) FindByToken(token string) (*models.WrapperLink, error) {
	var link models.WrapperLink
	err := r.db.Where("wrapper_token = ?", token).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// FindByVideoID finds wrapper link by video ID
func (r *WrapperLinkRepository) FindByVideoID(videoID uuid.UUID) (*models.WrapperLink, error) {
	var link models.WrapperLink
	err := r.db.Where("video_id = ?", videoID).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// FindBySlug finds wrapper link by slug
func (r *WrapperLinkRepository) FindBySlug(slug string) (*models.WrapperLink, error) {
	var link models.WrapperLink
	err := r.db.Where("slug = ?", slug).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// Update updates wrapper link
func (r *WrapperLinkRepository) Update(link *models.WrapperLink) error {
	return r.db.Save(link).Error
}

// Delete deletes wrapper link
func (r *WrapperLinkRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.WrapperLink{}, "id = ?", id).Error
}

// TokenExists checks if wrapper token exists
func (r *WrapperLinkRepository) TokenExists(token string) (bool, error) {
	var count int64
	err := r.db.Model(&models.WrapperLink{}).Where("wrapper_token = ?", token).Count(&count).Error
	return count > 0, err
}

// SlugExists checks if slug exists
func (r *WrapperLinkRepository) SlugExists(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.WrapperLink{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}