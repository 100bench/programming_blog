package domain

import (
	"time"
)

// Blog represents a blog post.
type Blog struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content"`
	Photo       string    `json:"photo"`
	TimeCreated time.Time `json:"time_created"`
	TimeUpdate  time.Time `json:"time_update"`
	IsPublished bool      `json:"is_published"`
	CategoryID  uint      `json:"category_id"`
	Category    *Category `json:"category,omitempty"` // Omitempty for optional eager loading
}

// BlogRepository defines the interface for interacting with Blog data.
type BlogRepository interface {
	Create(blog *Blog) error
	FindByID(id uint) (*Blog, error)
	FindBySlug(slug string) (*Blog, error)
	FindAll(publishedOnly bool) ([]Blog, error)
	FindByCategoryID(categoryID uint, publishedOnly bool) ([]Blog, error)
	Update(blog *Blog) error
	Delete(id uint) error
}
