package repositories

import (
	"bobastream/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PCloudCredentialRepository struct {
	db *gorm.DB
}

func NewPCloudCredentialRepository(db *gorm.DB) *PCloudCredentialRepository {
	return &PCloudCredentialRepository{db: db}
}

// Create creates a new pCloud credential
func (r *PCloudCredentialRepository) Create(credential *models.PCloudCredential) error {
	return r.db.Create(credential).Error
}

// FindByID finds pCloud credential by ID
func (r *PCloudCredentialRepository) FindByID(id uuid.UUID) (*models.PCloudCredential, error) {
	var credential models.PCloudCredential
	err := r.db.First(&credential, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

// Update updates pCloud credential
func (r *PCloudCredentialRepository) Update(credential *models.PCloudCredential) error {
	return r.db.Save(credential).Error
}

// Delete soft deletes pCloud credential
func (r *PCloudCredentialRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.PCloudCredential{}, "id = ?", id).Error
}

// GetAll gets all pCloud credentials
func (r *PCloudCredentialRepository) GetAll() ([]models.PCloudCredential, error) {
	var credentials []models.PCloudCredential
	err := r.db.Order("created_at DESC").Find(&credentials).Error
	return credentials, err
}

// GetActive gets all active pCloud credentials
func (r *PCloudCredentialRepository) GetActive() ([]models.PCloudCredential, error) {
	var credentials []models.PCloudCredential
	err := r.db.Where("is_active = ?", true).
		Order("storage_used_gb ASC"). // Order by least used first
		Find(&credentials).Error
	return credentials, err
}

// ToggleActive toggles credential active status
func (r *PCloudCredentialRepository) ToggleActive(id uuid.UUID, isActive bool) error {
	return r.db.Model(&models.PCloudCredential{}).Where("id = ?", id).
		Update("is_active", isActive).Error
}

// GetByStorageAvailable gets credentials ordered by available storage (most to least)
func (r *PCloudCredentialRepository) GetByStorageAvailable() ([]models.PCloudCredential, error) {
	var credentials []models.PCloudCredential
	err := r.db.Where("is_active = ?", true).
		Order("(storage_limit_gb - storage_used_gb) DESC").
		Find(&credentials).Error
	return credentials, err
}

// UpdateStorageUsed updates storage used for a credential
func (r *PCloudCredentialRepository) UpdateStorageUsed(id uuid.UUID, storageUsedGB float64) error {
	return r.db.Model(&models.PCloudCredential{}).Where("id = ?", id).
		Update("storage_used_gb", storageUsedGB).Error
}