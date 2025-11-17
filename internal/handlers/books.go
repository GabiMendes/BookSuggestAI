package handlers

import (
	"book-suggest-ai/internal/models"
	"book-suggest-ai/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type BookHandler struct {
	bookService *services.BookService
}

func NewBookHandler(bookService *services.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) ListSpreadsheets(c *gin.Context) {
	// Get access token from cookie
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not found"})
		return
	}

	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	spreadsheets, err := h.bookService.ListUserSpreadsheets(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list spreadsheets", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"spreadsheets": spreadsheets})
}

func (h *BookHandler) AnalyzeSpreadsheet(c *gin.Context) {
	spreadsheetID := c.Param("id")
	if spreadsheetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spreadsheet ID is required"})
		return
	}

	// Get user ID from session (simplified for now)
	_, err := c.Cookie("session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User session not found"})
		return
	}

	// For simplicity, using a fixed user ID (in production, you'd have a proper user table)
	userID := uint(1) // This should be retrieved from database based on session

	// Get access token
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not found"})
		return
	}

	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	analysis, err := h.bookService.AnalyzeSpreadsheet(c.Request.Context(), userID, token, spreadsheetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze spreadsheet", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analysis": analysis})
}

func (h *BookHandler) GenerateRecommendations(c *gin.Context) {
	fmt.Printf("🚀 Handler: Iniciando geração de recomendações\n")

	var request models.RecommendationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("❌ Handler: Erro no parse do request: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}
	fmt.Printf("✅ Handler: Request parseado - SpreadsheetID: %s\n", request.SpreadsheetID)

	// Get user ID from session
	_, err := c.Cookie("session")
	if err != nil {
		fmt.Printf("❌ Handler: Session não encontrada\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User session not found"})
		return
	}

	// Convert to uint (in production, you'd lookup the actual user ID)
	userID := uint(1) // Simplified
	fmt.Printf("👤 Handler: UserID definido como %d\n", userID)

	// Get access token
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		fmt.Printf("❌ Handler: Access token não encontrado\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not found"})
		return
	}
	fmt.Printf("🔑 Handler: Access token obtido\n")

	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	fmt.Printf("📞 Handler: Chamando bookService.GenerateRecommendations\n")
	recommendations, err := h.bookService.GenerateRecommendations(
		c.Request.Context(),
		userID,
		token,
		request.SpreadsheetID,
	)
	if err != nil {
		fmt.Printf("❌ Handler: Erro no bookService: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate recommendations", "details": err.Error()})
		return
	}

	fmt.Printf("🎉 Handler: Recomendações geradas com sucesso - Tipo: %T\n", recommendations)
	c.JSON(http.StatusOK, recommendations)
}

func (h *BookHandler) GetRecommendationsHistory(c *gin.Context) {
	// Get user ID from session
	_, err := c.Cookie("session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User session not found"})
		return
	}

	userID := uint(1) // Simplified

	recommendations, err := h.bookService.GetRecommendationsHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations history", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}
