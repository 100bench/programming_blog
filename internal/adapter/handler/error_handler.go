package handler

import (
	"net/http"
	"programming_blog_go/internal/domain"
	"programming_blog_go/internal/usecase"

	"github.com/gin-gonic/gin"
)

// HandleError centralizes error handling and maps internal errors to HTTP responses.
func HandleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case domain.ErrAlreadyExists:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case domain.ErrInvalidInput:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case usecase.ErrInvalidCredentials:
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case usecase.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case usecase.ErrUserAlreadyExists:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	// Add more specific error mappings here
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
	}
	c.Abort() // Stop further processing of the request
}
