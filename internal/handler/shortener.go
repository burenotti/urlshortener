package handler

import (
	"errors"
	"github.com/burenotti/urlshortener/internal/storage"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/url"
)

type createLinkResponse struct {
	LinkID string `json:"link_id,"`
	Url    string `json:"url,"`
}

// @Summary Creates a short link
// @Tags link
// @Accept json
// @Produce json
// @Param        url    query     string  true  "url to shorten"  Format(url)
// @Success 201 {object} handler.createLinkResponse
// @Failure 422 {object} handler.JSONError
// @Failure 400 {object} handler.JSONError
// @Failure 500 {object} handler.JSONError
// @Router /api/link/ [post]
func (h *Handler) createLink(c *gin.Context) {
	srcUrl := c.Query("url")
	if srcUrl == "" {
		AbortWithJSONError(c, 422, "url is required query parameter")
		return
	}

	u, err := url.Parse(srcUrl)

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	if err != nil || u.Host == "" {
		AbortWithJSONError(c, 422, "url is incorrect")
		return
	}

	linkID, err := h.services.CreateShortLink(c, srcUrl)

	// Partial save is non-critical error
	if err != nil && !errors.Is(err, storage.ErrPartialSave) {
		AbortWithJSONError(c, 500, err.Error())
		return
	}

	c.JSON(201, createLinkResponse{
		LinkID: linkID,
		Url:    h.basePath + "/" + linkID,
	})
}

type getLinkInfoResponse struct {
	Url string `json:"url"`
}

// @Summary Get information about shortened link
// @Tags link
// @Produce json
// @Param        link_id    path     string  true  "Link ID"  Format(url)
// @Success 200 {object} handler.getLinkInfoResponse
// @Failure 422 {object} handler.JSONError
// @Failure 400 {object} handler.JSONError
// @Failure 500 {object} handler.JSONError
// @Router /api/link/{link_id} [get]
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

	c.JSON(200, getLinkInfoResponse{srcUrl})
}

// @Summary Redirect to source url
// @Tags link
// @Param link_id path string true "link ID"
// @Success 308
// @Failure 404 {string} Not Found
// @Failure 500 {string} Internal Server Error
// @Router /l/{link_id} [get]
func (h *Handler) redirect(c *gin.Context) {
	linkID := c.Param("link_id")
	if linkID == "" {
		AbortWithJSONError(c, 400, "link_id is required path parameter")
		return
	}
	session := sessions.Default(c)
	sessionId, ok := session.Get("session_id").(string)
	if !ok {
		id, _ := uuid.NewRandom()
		sessionId = id.String()
		session.Set("session_id", sessionId)
		err := session.Save()
		if err != nil {
			logrus.WithField("error", err).Errorf("Can't save session_id")
		}
	}

	srcUrl, err := h.services.GetSourceForRedirect(c, sessionId, linkID)
	if errors.Is(err, storage.ErrNoSourceUrl) {
		c.AbortWithStatus(404)
		return
	} else if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.Redirect(307, srcUrl)
}
