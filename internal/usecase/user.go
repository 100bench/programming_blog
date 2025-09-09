package usecase

import (
	"errors"
	"programming_blog_go/internal/domain"
	"programming_blog_go/internal/utils"
	"time"

	"gorm.io/gorm"
)

// Use specific domain errors for consistency
var (
	ErrUserAlreadyExists  = domain.ErrAlreadyExists           // Use domain error
	ErrUserNotFound       = domain.ErrNotFound                // Use domain error
	ErrInvalidCredentials = errors.New("invalid credentials") // Specific to auth
)

// RegisterUserUseCase handles new user registration.
type RegisterUserUseCase struct {
	UserRepository domain.UserRepository
}

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (uc *RegisterUserUseCase) Execute(req RegisterUserRequest) (*domain.User, error) {
	// Check if user already exists by username or email
	_, err := uc.UserRepository.FindByUsername(req.Username)
	if err == nil { // User found
		return nil, ErrUserAlreadyExists
	}
	_, err = uc.UserRepository.FindByEmail(req.Email)
	if err == nil { // User found
		return nil, ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = uc.UserRepository.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUserUseCase handles user login and authentication.
type AuthenticateUserUseCase struct {
	UserRepository domain.UserRepository
}

type AuthenticateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (uc *AuthenticateUserUseCase) Execute(req AuthenticateUserRequest) (*domain.User, error) {
	user, err := uc.UserRepository.FindByUsername(req.Username)
	if err != nil {
		// Assuming FindByUsername returns an error like gorm.ErrRecordNotFound if not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Compare the provided password with the stored hashed password
	err = utils.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
