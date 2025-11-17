package main

import (
	"book-suggest-ai/internal/config"
	"book-suggest-ai/internal/handlers"
	"book-suggest-ai/internal/middleware"
	"book-suggest-ai/internal/services"
	"book-suggest-ai/pkg/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg.DBPath)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	authService := services.NewAuthService(cfg)
	sheetsService := services.NewSheetsService(cfg)
	aiService := services.NewAIService(cfg)
	bookService := services.NewBookService(db, sheetsService, aiService)

	// Setup routes
	router := setupRouter(cfg, authService, bookService)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

func setupRouter(cfg *config.Config, authService *services.AuthService, bookService *services.BookService) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Static files
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	bookHandler := handlers.NewBookHandler(bookService)
	webHandler := handlers.NewWebHandler()

	// Web routes
	router.GET("/", webHandler.Index)
	router.GET("/dashboard", webHandler.Dashboard) // Removi RequireAuth temporariamente para debug

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.GET("/login", authHandler.Login)
		auth.GET("/google", authHandler.GoogleLogin)
		auth.GET("/google/callback", authHandler.GoogleCallback)
		auth.GET("/callback", authHandler.GoogleCallback) // Rota alternativa para compatibilidade
		auth.POST("/logout", authHandler.Logout)
	}

	// API routes
	api := router.Group("/api/v1")
	// Removido RequireAuth temporariamente para debug
	{
		api.GET("/spreadsheets", bookHandler.ListSpreadsheets)
		api.GET("/spreadsheets/:id/analyze", bookHandler.AnalyzeSpreadsheet)
		api.POST("/recommendations", bookHandler.GenerateRecommendations)
		api.GET("/recommendations/history", bookHandler.GetRecommendationsHistory)
	}

	return router
}
