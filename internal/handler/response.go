package handler

import "github.com/gin-gonic/gin"

func AbortWithJSONError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error": message,
	})
}
