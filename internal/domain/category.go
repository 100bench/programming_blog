package domain

import "time"

// Category represents a blog post category.
type Category struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CategoryRepository defines the interface for interacting with Category data.
type CategoryRepository interface {
	Create(category *Category) error
	FindByID(id uint) (*Category, error)
	FindBySlug(slug string) (*Category, error)
	FindAll() ([]Category, error)
	Update(category *Category) error
	Delete(id uint) error
}
