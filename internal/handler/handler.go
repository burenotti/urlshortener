package handler

import (
	"github.com/burenotti/urlshortener/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	services *service.Service
	basePath string
}

func NewHandler(services *service.Service, basePath string) *Handler {
	return &Handler{
		services: services,
		basePath: basePath,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/link", h.createLink)
		api.GET("/link/:link_id", h.getLinkInfo)
	}

	router.GET("/l/:link_id", h.redirect)

	return router
}
