package usecase

import (
	"programming_blog_go/internal/domain"
)

// GetAllCategoriesUseCase retrieves all categories.
type GetAllCategoriesUseCase struct {
	CategoryRepository domain.CategoryRepository
}

// Execute retrieves all categories from the repository.
func (uc *GetAllCategoriesUseCase) Execute() ([]domain.Category, error) {
	return uc.CategoryRepository.FindAll()
}
