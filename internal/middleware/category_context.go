package middleware

import (
	"log"
	"programming_blog_go/internal/usecase"

	"github.com/gin-gonic/gin"
)

// CategoryContextMiddleware fetches all categories and adds them to the Gin context.
func CategoryContextMiddleware(getAllCategoriesUC *usecase.GetAllCategoriesUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("categories", []interface{}{}) // Initialize with empty slice to avoid nil pointer issues in templates

		categories, err := getAllCategoriesUC.Execute()
		if err != nil {
			log.Printf("Error fetching categories for context: %v", err)
			// Continue processing request even if categories can't be fetched
		} else {
			// Convert domain.Category to interface{} for template consumption
			interfaceCategories := make([]interface{}, len(categories))
			for i, cat := range categories {
				interfaceCategories[i] = cat
			}
			c.Set("categories", interfaceCategories)
		}
		c.Next()
	}
}
