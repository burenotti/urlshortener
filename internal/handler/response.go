package handler

import "github.com/gin-gonic/gin"

type JSONError struct {
	Error string `json:"error"`
}

func AbortWithJSONError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, JSONError{message})
}
