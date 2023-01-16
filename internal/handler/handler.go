package handler

import (
	_ "github.com/burenotti/urlshortener/docs"
	"github.com/burenotti/urlshortener/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Handler struct {
	services *service.Service
	basePath string
	router   *gin.Engine
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.router.ServeHTTP(writer, request)
}

func NewHandler(services *service.Service, store sessions.Store, basePath string) *Handler {

	r := gin.Default()

	r.Use(sessions.Sessions("session", store))

	handler := Handler{
		services: services,
		basePath: basePath,
		router:   r,
	}

	handler.initRoutes()

	return &handler
}

func (h *Handler) initRoutes() {

	api := h.router.Group("/api")
	{
		api.POST("/link", h.createLink)
		api.GET("/link/:link_id", h.getLinkInfo)
	}

	h.router.GET("/l/:link_id", h.redirect)

	h.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
