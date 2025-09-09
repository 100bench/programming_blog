package utils

import (
	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext retrieves the user ID from the Gin context.
// It returns the user ID and a boolean indicating if it was found and is valid.
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	// Type assertion from interface{} to uint
	id, ok := userID.(uint)
	return id, ok
}

// GetUsernameFromContext retrieves the username from the Gin context.
// It returns the username and a boolean indicating if it was found and is valid.
func GetUsernameFromContext(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	// Type assertion from interface{} to string
	name, ok := username.(string)
	return name, ok
}
