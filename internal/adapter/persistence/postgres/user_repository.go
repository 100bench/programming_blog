package postgres

import (
	"errors"
	"programming_blog_go/internal/domain"

	"gorm.io/gorm"
)

// UserRepository implements domain.UserRepository for PostgreSQL.
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new PostgreSQL user repository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create creates a new user in the database.
func (r *UserRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

// FindByID finds a user by its ID.
func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by their username.
func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by their email.
func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user.
func (r *UserRepository) Update(user *domain.User) error {
	return r.DB.Save(user).Error
}

// Delete deletes a user by their ID.
func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.User{}, id).Error
}
