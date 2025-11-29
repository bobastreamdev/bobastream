package repositories

import (
	"bobastream/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID finds category by ID
func (r *CategoryRepository) FindByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindBySlug finds category by slug
func (r *CategoryRepository) FindBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Update updates category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete soft deletes category
func (r *CategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Category{}, "id = ?", id).Error
}

// GetAll gets all categories
func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Order("display_order ASC, name ASC").Find(&categories).Error
	return categories, err
}

// GetActive gets all active categories
func (r *CategoryRepository) GetActive() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("is_active = ?", true).
		Order("display_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// ToggleActive toggles category active status
func (r *CategoryRepository) ToggleActive(id uuid.UUID, isActive bool) error {
	return r.db.Model(&models.Category{}).Where("id = ?", id).
		Update("is_active", isActive).Error
}

// SlugExists checks if slug exists
func (r *CategoryRepository) SlugExists(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

// NameExists checks if name exists
func (r *CategoryRepository) NameExists(name string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// UpdateDisplayOrder updates display order
func (r *CategoryRepository) UpdateDisplayOrder(id uuid.UUID, order int) error {
	return r.db.Model(&models.Category{}).Where("id = ?", id).
		Update("display_order", order).Error
}