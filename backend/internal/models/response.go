package models

import "github.com/gin-gonic/gin"

// Respond sends a consistent JSON response across all endpoints.
func Respond(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{
		"message": message,
		"data":    data,
	})
}
