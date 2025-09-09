package usecase

import (
	"errors"
	"programming_blog_go/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBlogRepository is a mock implementation of domain.BlogRepository
type MockBlogRepository struct {
	mock.Mock
}

func (m *MockBlogRepository) Create(blog *domain.Blog) error {
	args := m.Called(blog)
	return args.Error(0)
}

func (m *MockBlogRepository) FindByID(id uint) (*domain.Blog, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Blog), args.Error(1)
}

func (m *MockBlogRepository) FindBySlug(slug string) (*domain.Blog, error) {
	args := m.Called(slug)
	return args.Get(0).(*domain.Blog), args.Error(1)
}

func (m *MockBlogRepository) FindAll(publishedOnly bool) ([]domain.Blog, error) {
	args := m.Called(publishedOnly)
	return args.Get(0).([]domain.Blog), args.Error(1)
}

func (m *MockBlogRepository) FindByCategoryID(categoryID uint, publishedOnly bool) ([]domain.Blog, error) {
	args := m.Called(categoryID, publishedOnly)
	return args.Get(0).([]domain.Blog), args.Error(1)
}

func (m *MockBlogRepository) Update(blog *domain.Blog) error {
	args := m.Called(blog)
	return args.Error(0)
}

func (m *MockBlogRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockCategoryRepository is a mock implementation of domain.CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) FindByID(id uint) (*domain.Category, error) {
	args := m.Called(id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindBySlug(slug string) (*domain.Category, error) {
	args := m.Called(slug)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindAll() ([]domain.Category, error) {
	args := m.Called()
	return args.Get(0).([]domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetBlogPostsUseCase_Execute(t *testing.T) {
	mockRepo := new(MockBlogRepository)
	usecase := &GetBlogPostsUseCase{BlogRepository: mockRepo}

	expectedBlogs := []domain.Blog{
		{ID: 1, Title: "Test Post 1", IsPublished: true},
		{ID: 2, Title: "Test Post 2", IsPublished: true},
	}
	mockRepo.On("FindAll", true).Return(expectedBlogs, nil)

	blogs, err := usecase.Execute(true)

	assert.NoError(t, err)
	assert.NotNil(t, blogs)
	assert.Len(t, blogs, 2)
	assert.Equal(t, expectedBlogs, blogs)
	mockRepo.AssertExpectations(t)
}

func TestGetBlogPostsByCategoryUseCase_Execute(t *testing.T) {
	mockBlogRepo := new(MockBlogRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	usecase := &GetBlogPostsByCategoryUseCase{
		BlogRepository:     mockBlogRepo,
		CategoryRepository: mockCategoryRepo,
	}

	categorySlug := "go-lang"
	categoryID := uint(1)
	expectedCategory := &domain.Category{ID: categoryID, Name: "Go Lang", Slug: categorySlug}
	expectedBlogs := []domain.Blog{
		{ID: 1, Title: "Go Post 1", CategoryID: categoryID, IsPublished: true},
	}

	// Test case: Category found, posts found
	mockCategoryRepo.On("FindBySlug", categorySlug).Return(expectedCategory, nil).Once()
	mockBlogRepo.On("FindByCategoryID", categoryID, true).Return(expectedBlogs, nil).Once()

	blogs, err := usecase.Execute(categorySlug, true)
	assert.NoError(t, err)
	assert.Len(t, blogs, 1)
	assert.Equal(t, expectedBlogs, blogs)

	// Test case: Category not found
	mockCategoryRepo.On("FindBySlug", "non-existent").Return(nil, domain.ErrNotFound).Once()

	blogs, err = usecase.Execute("non-existent", true)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, blogs)

	// Test case: Error fetching category
	mockCategoryRepo.On("FindBySlug", "error-slug").Return(nil, errors.New("db error")).Once()

	blogs, err = usecase.Execute("error-slug", true)
	assert.Error(t, err)
	assert.Nil(t, blogs)

	mockBlogRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}

func TestGetBlogPostBySlugUseCase_Execute(t *testing.T) {
	mockRepo := new(MockBlogRepository)
	usecase := &GetBlogPostBySlugUseCase{BlogRepository: mockRepo}

	slug := "test-post"
	expectedBlog := &domain.Blog{ID: 1, Title: "Test Post", Slug: slug, IsPublished: true}

	// Test case: Post found
	mockRepo.On("FindBySlug", slug).Return(expectedBlog, nil).Once()

	blog, err := usecase.Execute(slug)
	assert.NoError(t, err)
	assert.NotNil(t, blog)
	assert.Equal(t, expectedBlog, blog)

	// Test case: Post not found
	mockRepo.On("FindBySlug", "non-existent").Return(nil, domain.ErrNotFound).Once()

	blog, err = usecase.Execute("non-existent")
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, blog)

	// Test case: Error fetching post
	mockRepo.On("FindBySlug", "error-slug").Return(nil, errors.New("db error")).Once()

	blog, err = usecase.Execute("error-slug")
	assert.Error(t, err)
	assert.Nil(t, blog)

	mockRepo.AssertExpectations(t)
}

func TestCreateBlogPostUseCase_Execute(t *testing.T) {
	mockBlogRepo := new(MockBlogRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	usecase := &CreateBlogPostUseCase{
		BlogRepository:     mockBlogRepo,
		CategoryRepository: mockCategoryRepo,
	}

	categoryID := uint(1)
	existingCategory := &domain.Category{ID: categoryID, Name: "Go Lang"}

	request := CreateBlogPostRequest{
		Title:       "New Post",
		Slug:        "new-post",
		Content:     "Some content",
		Photo:       "photo.jpg",
		IsPublished: true,
		CategoryID:  categoryID,
	}

	// Test case: Successful creation
	mockCategoryRepo.On("FindByID", categoryID).Return(existingCategory, nil).Once()
	mockBlogRepo.On("Create", mock.AnythingOfType("*domain.Blog")).Return(nil).Once()

	blog, err := usecase.Execute(request)
	assert.NoError(t, err)
	assert.NotNil(t, blog)
	assert.Equal(t, request.Title, blog.Title)
	assert.Equal(t, request.Slug, blog.Slug)
	assert.Equal(t, request.CategoryID, blog.CategoryID)
	mockBlogRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)

	// Test case: Category not found
	mockCategoryRepo.On("FindByID", uint(999)).Return(nil, domain.ErrNotFound).Once()
	request.CategoryID = uint(999)

	blog, err = usecase.Execute(request)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, blog)
	mockCategoryRepo.AssertExpectations(t)

	// Test case: Error creating blog post
	mockCategoryRepo.On("FindByID", categoryID).Return(existingCategory, nil).Once()
	mockBlogRepo.On("Create", mock.AnythingOfType("*domain.Blog")).Return(errors.New("db error")).Once()
	request.CategoryID = categoryID // Reset for this test

	blog, err = usecase.Execute(request)
	assert.Error(t, err)
	assert.Nil(t, blog)
	mockBlogRepo.AssertExpectations(t)
	mockCategoryRepo.AssertExpectations(t)
}
