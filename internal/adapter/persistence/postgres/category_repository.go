package postgres

import (
	"errors"
	"programming_blog_go/internal/domain"

	"gorm.io/gorm"
)

// CategoryRepository implements domain.CategoryRepository for PostgreSQL.
type CategoryRepository struct {
	DB *gorm.DB
}

// NewCategoryRepository creates a new PostgreSQL category repository.
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

// Create creates a new category in the database.
func (r *CategoryRepository) Create(category *domain.Category) error {
	return r.DB.Create(category).Error
}

// FindByID finds a category by its ID.
func (r *CategoryRepository) FindByID(id uint) (*domain.Category, error) {
	var category domain.Category
	if err := r.DB.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if not found
		}
		return nil, err
	}
	return &category, nil
}

// FindBySlug finds a category by its slug.
func (r *CategoryRepository) FindBySlug(slug string) (*domain.Category, error) {
	var category domain.Category
	if err := r.DB.Where("slug = ?", slug).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if not found
		}
		return nil, err
	}
	return &category, nil
}

// FindAll retrieves all categories.
func (r *CategoryRepository) FindAll() ([]domain.Category, error) {
	var categories []domain.Category
	if err := r.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Update updates an existing category.
func (r *CategoryRepository) Update(category *domain.Category) error {
	return r.DB.Save(category).Error
}

// Delete deletes a category by its ID.
func (r *CategoryRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.Category{}, id).Error
}
