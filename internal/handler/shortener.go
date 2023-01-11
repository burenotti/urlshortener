package handler

import (
	"errors"
	"github.com/burenotti/urlshortener/internal/storage"
	"github.com/gin-gonic/gin"
	"net/url"
)

func (h *Handler) createLink(c *gin.Context) {
	srcUrl := c.Query("url")
	if srcUrl == "" {
		AbortWithJSONError(c, 422, "url is required query parameter")
		return
	}

	_, err := url.Parse(srcUrl)

	if err != nil {
		AbortWithJSONError(c, 422, "url is incorrect")
		return
	}

	linkID, err := h.services.CreateShortLink(c, srcUrl)

	// Partial save is non-critical error
	if err != nil && !errors.Is(err, storage.ErrPartialSave) {
		AbortWithJSONError(c, 500, err.Error())
		return
	}

	c.JSON(201, gin.H{
		"link_id": linkID,
		"url":     h.basePath + "/" + linkID,
	})
}

func (h *Handler) getLinkInfo(c *gin.Context) {
	linkID := c.Param("link_id")
	if linkID == "" {
		AbortWithJSONError(c, 400, "link_id is required path parameter")
		return
	}
	srcUrl, err := h.services.GetSource(c, linkID)
	if errors.Is(err, storage.ErrNoSourceUrl) {
		AbortWithJSONError(c, 400, err.Error())
		return
	} else if err != nil {
		AbortWithJSONError(c, 500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"url": srcUrl,
	})
}

func (h *Handler) redirect(c *gin.Context) {
	linkID := c.Param("link_id")
	if linkID == "" {
		AbortWithJSONError(c, 400, "link_id is required path parameter")
		return
	}
	srcUrl, err := h.services.GetSource(c, linkID)
	if errors.Is(err, storage.ErrNoSourceUrl) {
		c.AbortWithStatus(404)
		return
	} else if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.Redirect(307, srcUrl)
}
