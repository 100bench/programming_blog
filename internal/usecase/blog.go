package usecase

import (
	"errors"
	"programming_blog_go/internal/domain"
	"time"

	"gorm.io/gorm"
)

// GetBlogPostsUseCase retrieves all published blog posts or all posts if publishedOnly is false.
type GetBlogPostsUseCase struct {
	BlogRepository domain.BlogRepository
}

func (uc *GetBlogPostsUseCase) Execute(publishedOnly bool) ([]domain.Blog, error) {
	return uc.BlogRepository.FindAll(publishedOnly)
}

// GetBlogPostsByCategoryUseCase retrieves published blog posts by category slug.
type GetBlogPostsByCategoryUseCase struct {
	BlogRepository     domain.BlogRepository
	CategoryRepository domain.CategoryRepository
}

func (uc *GetBlogPostsByCategoryUseCase) Execute(categorySlug string, publishedOnly bool) ([]domain.Blog, error) {
	category, err := uc.CategoryRepository.FindBySlug(categorySlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	if category == nil {
		return []domain.Blog{}, nil // Return empty slice if category not found
	}
	return uc.BlogRepository.FindByCategoryID(category.ID, publishedOnly)
}

// GetBlogPostBySlugUseCase retrieves a single blog post by its slug.
type GetBlogPostBySlugUseCase struct {
	BlogRepository domain.BlogRepository
}

func (uc *GetBlogPostBySlugUseCase) Execute(postSlug string) (*domain.Blog, error) {
	post, err := uc.BlogRepository.FindBySlug(postSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return post, nil
}

// CreateBlogPostUseCase handles the creation of a new blog post.
type CreateBlogPostUseCase struct {
	BlogRepository     domain.BlogRepository
	CategoryRepository domain.CategoryRepository
}

type CreateBlogPostRequest struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Content     string `json:"content"`
	Photo       string `json:"photo"`
	IsPublished bool   `json:"is_published"`
	CategoryID  uint   `json:"category_id" binding:"required"`
}

func (uc *CreateBlogPostUseCase) Execute(req CreateBlogPostRequest) (*domain.Blog, error) {
	// Check if category exists
	_, err := uc.CategoryRepository.FindByID(req.CategoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound // Category not found
		}
		return nil, err // Other error
	}

	blog := &domain.Blog{
		Title:       req.Title,
		Slug:        req.Slug,
		Content:     req.Content,
		Photo:       req.Photo,
		TimeCreated: time.Now(),
		TimeUpdate:  time.Now(),
		IsPublished: req.IsPublished,
		CategoryID:  req.CategoryID,
	}

	err = uc.BlogRepository.Create(blog)
	if err != nil {
		return nil, err
	}
	return blog, nil
}
