package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebHandler struct{}

func NewWebHandler() *WebHandler {
	return &WebHandler{}
}

func (h *WebHandler) Index(c *gin.Context) {
	// Check if user is already logged in
	if _, err := c.Cookie("session"); err == nil {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "BookSuggest AI - Recomendações Inteligentes de Livros",
	})
}

func (h *WebHandler) Dashboard(c *gin.Context) {
	// Get user info from session
	sessionID, err := c.Cookie("session")
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":   "Dashboard - BookSuggest AI",
		"user_id": sessionID,
	})
}
