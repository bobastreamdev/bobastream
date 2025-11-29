package services

import (
	"bobastream/internal/models"
	"bobastream/internal/repositories"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(name, description, icon string) (*models.Category, error) {
	// Check if name exists
	exists, err := s.categoryRepo.NameExists(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category name already exists")
	}

	// Generate slug from name
	slug := s.generateSlug(name)

	// Check if slug exists
	slugExists, err := s.categoryRepo.SlugExists(slug)
	if err != nil {
		return nil, err
	}
	if slugExists {
		return nil, errors.New("category slug already exists")
	}

	category := &models.Category{
		Name:        name,
		Slug:        slug,
		Description: description,
		Icon:        icon,
		IsActive:    true,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory updates category
func (s *CategoryService) UpdateCategory(id uuid.UUID, name, description, icon string, displayOrder int) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// If name changed, regenerate slug and check uniqueness
	if name != "" && name != category.Name {
		slug := s.generateSlug(name)
		slugExists, err := s.categoryRepo.SlugExists(slug)
		if err != nil {
			return nil, err
		}
		if slugExists {
			return nil, errors.New("category slug already exists")
		}
		category.Name = name
		category.Slug = slug
	}

	if description != "" {
		category.Description = description
	}

	if icon != "" {
		category.Icon = icon
	}

	category.DisplayOrder = displayOrder

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory soft deletes category
func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	return s.categoryRepo.Delete(id)
}

// GetCategoryByID gets category by ID
func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

// GetCategoryBySlug gets category by slug
func (s *CategoryService) GetCategoryBySlug(slug string) (*models.Category, error) {
	return s.categoryRepo.FindBySlug(slug)
}

// GetAllCategories gets all categories
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

// GetActiveCategories gets active categories only
func (s *CategoryService) GetActiveCategories() ([]models.Category, error) {
	return s.categoryRepo.GetActive()
}

// ToggleActive toggles category active status
func (s *CategoryService) ToggleActive(id uuid.UUID, isActive bool) error {
	return s.categoryRepo.ToggleActive(id, isActive)
}

// UpdateDisplayOrder updates category display order
func (s *CategoryService) UpdateDisplayOrder(id uuid.UUID, order int) error {
	return s.categoryRepo.UpdateDisplayOrder(id, order)
}

// generateSlug generates URL-friendly slug from name
func (s *CategoryService) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters (basic sanitization)
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	return slug
}