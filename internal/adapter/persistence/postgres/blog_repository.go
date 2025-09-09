package postgres

import (
	"errors"
	"programming_blog_go/internal/domain"

	"gorm.io/gorm"
)

// BlogRepository implements domain.BlogRepository for PostgreSQL.
type BlogRepository struct {
	DB *gorm.DB
}

// NewBlogRepository creates a new PostgreSQL blog repository.
func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{DB: db}
}

// Create creates a new blog post in the database.
func (r *BlogRepository) Create(blog *domain.Blog) error {
	return r.DB.Create(blog).Error
}

// FindByID finds a blog post by its ID.
func (r *BlogRepository) FindByID(id uint) (*domain.Blog, error) {
	var blog domain.Blog
	if err := r.DB.Preload("Category").First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &blog, nil
}

// FindBySlug finds a blog post by its slug.
func (r *BlogRepository) FindBySlug(slug string) (*domain.Blog, error) {
	var blog domain.Blog
	if err := r.DB.Preload("Category").Where("slug = ?", slug).First(&blog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &blog, nil
}

// FindAll retrieves all blog posts, optionally filtered by published status.
func (r *BlogRepository) FindAll(publishedOnly bool) ([]domain.Blog, error) {
	var blogs []domain.Blog
	query := r.DB.Preload("Category")
	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}
	if err := query.Order("time_created DESC").Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// FindByCategoryID retrieves blog posts by category ID, optionally filtered by published status.
func (r *BlogRepository) FindByCategoryID(categoryID uint, publishedOnly bool) ([]domain.Blog, error) {
	var blogs []domain.Blog
	query := r.DB.Preload("Category").Where("category_id = ?", categoryID)
	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}
	if err := query.Order("time_created DESC").Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// Update updates an existing blog post.
func (r *BlogRepository) Update(blog *domain.Blog) error {
	return r.DB.Save(blog).Error
}

// Delete deletes a blog post by its ID.
func (r *BlogRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.Blog{}, id).Error
}
