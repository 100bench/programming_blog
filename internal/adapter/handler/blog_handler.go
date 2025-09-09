package handler

import (
	"net/http"
	"strconv"

	"programming_blog_go/internal/domain"
	"programming_blog_go/internal/usecase"

	"github.com/gin-gonic/gin"
)

// BlogHandler handles HTTP requests related to blog posts and categories.
type BlogHandler struct {
	GetBlogPostsUseCase           *usecase.GetBlogPostsUseCase
	GetBlogPostsByCategoryUseCase *usecase.GetBlogPostsByCategoryUseCase
	GetBlogPostBySlugUseCase      *usecase.GetBlogPostBySlugUseCase
	CreateBlogPostUseCase         *usecase.CreateBlogPostUseCase
}

// NewBlogHandler creates a new BlogHandler.
func NewBlogHandler(
	getBlogPostsUC *usecase.GetBlogPostsUseCase,
	getBlogPostsByCategoryUC *usecase.GetBlogPostsByCategoryUseCase,
	getBlogPostBySlugUC *usecase.GetBlogPostBySlugUseCase,
	createBlogPostUC *usecase.CreateBlogPostUseCase,
) *BlogHandler {
	return &BlogHandler{
		GetBlogPostsUseCase:           getBlogPostsUC,
		GetBlogPostsByCategoryUseCase: getBlogPostsByCategoryUC,
		GetBlogPostBySlugUseCase:      getBlogPostBySlugUC,
		CreateBlogPostUseCase:         createBlogPostUC,
	}
}

// GetBlogPosts handles the request to get all blog posts.
func (h *BlogHandler) GetBlogPosts(c *gin.Context) {
	posts, err := h.GetBlogPostsUseCase.Execute(true)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"posts": posts, "title": "Главная страница"})
}

// GetBlogPostsByCategory handles the request to get blog posts by category slug.
func (h *BlogHandler) GetBlogPostsByCategory(c *gin.Context) {
	categorySlug := c.Param("cat_slug")
	posts, err := h.GetBlogPostsByCategoryUseCase.Execute(categorySlug, true)
	if err != nil {
		HandleError(c, err)
		return
	}

	// TODO: Fetch category name for the title (will need a category use case)
	c.HTML(http.StatusOK, "index.html", gin.H{"posts": posts, "title": "Категория - " + categorySlug})
}

// GetBlogPost handles the request to get a single blog post by slug.
func (h *BlogHandler) GetBlogPost(c *gin.Context) {
	postSlug := c.Param("post_slug")
	post, err := h.GetBlogPostBySlugUseCase.Execute(postSlug)
	if err != nil {
		HandleError(c, err)
		return
	}
	// The use case now returns domain.ErrNotFound if not found, so this check is redundant.
	// if post == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
	// 	return
	// }
	c.HTML(http.StatusOK, "post.html", gin.H{"post": post, "title": post.Title})
}

// CreateBlogPost handles the request to create a new blog post.
func (h *BlogHandler) CreateBlogPost(c *gin.Context) {
	var req usecase.CreateBlogPostRequest
	if err := c.ShouldBind(&req); err != nil { // Changed from ShouldBindJSON to ShouldBind
		HandleError(c, domain.ErrInvalidInput) // Map binding errors to InvalidInput
		return
	}

	blog, err := h.CreateBlogPostUseCase.Execute(req)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, blog)
}

// AddPostPage renders the form for adding a new post.
func (h *BlogHandler) AddPostPage(c *gin.Context) {
	c.HTML(http.StatusOK, "addpage.html", gin.H{"title": "Добавление статьи"})
}

// Temporary solution to get categories for layout
// In a real application, this might be handled by a middleware or a more generic context provider
func (h *BlogHandler) GetCategoriesForLayout(c *gin.Context) {
	// This is a placeholder. In a proper Gin setup, you might have a middleware
	// that fetches common data like categories and injects them into the context
	// or directly into the template data.
	// For now, we'll just return a dummy list.
	c.Set("categories", []string{"Category 1", "Category 2"}) // Example
}

// GetCategoryIDFromName fetches category ID from category name. This is a placeholder for now.
func (h *BlogHandler) GetCategoryIDFromName(c *gin.Context) {
	categoryName := c.Param("name")
	// This function would typically interact with a Category Use Case
	// For now, let's just return a dummy ID
	id, _ := strconv.Atoi(categoryName) // In a real scenario, this would be a lookup
	c.JSON(http.StatusOK, gin.H{"id": id})
}
