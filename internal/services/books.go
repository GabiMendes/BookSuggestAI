package services

import (
	"book-suggest-ai/internal/models"
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type BookService struct {
	db            *gorm.DB
	sheetsService *SheetsService
	aiService     *AIService
}

func NewBookService(db *gorm.DB, sheetsService *SheetsService, aiService *AIService) *BookService {
	return &BookService{
		db:            db,
		sheetsService: sheetsService,
		aiService:     aiService,
	}
}

func (s *BookService) ListUserSpreadsheets(ctx context.Context, token *oauth2.Token) ([]models.GoogleSpreadsheet, error) {
	return s.sheetsService.ListSpreadsheets(ctx, token)
}

func (s *BookService) AnalyzeSpreadsheet(ctx context.Context, userID uint, token *oauth2.Token, spreadsheetID string) (*models.SpreadsheetAnalysis, error) {
	// Validate spreadsheet structure
	if err := s.sheetsService.ValidateSpreadsheetStructure(ctx, token, spreadsheetID); err != nil {
		return nil, fmt.Errorf("invalid spreadsheet structure: %w", err)
	}

	// Read books from spreadsheet
	books, err := s.sheetsService.ReadBooksFromSpreadsheet(ctx, token, spreadsheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to read books from spreadsheet: %w", err)
	}

	if len(books) == 0 {
		return nil, fmt.Errorf("no books found in spreadsheet")
	}

	// Analyze reading patterns
	patterns := s.analyzeUserPatterns(books)

	// Create analysis record
	analysis := &models.SpreadsheetAnalysis{
		UserID:          userID,
		SpreadsheetID:   spreadsheetID,
		TotalBooks:      len(books),
		FavoriteGenres:  patterns.FavoriteGenres,
		FavoriteAuthors: patterns.FavoriteAuthors,
		AverageRating:   patterns.AverageRating,
		ReadingPatterns: patterns.Summary,
	}

	// Save to database
	if err := s.db.Create(analysis).Error; err != nil {
		return nil, fmt.Errorf("failed to save analysis: %w", err)
	}

	return analysis, nil
}

func (s *BookService) GenerateRecommendations(ctx context.Context, userID uint, token *oauth2.Token, spreadsheetID string) (*models.RecommendationResponse, error) {
	// Read books from spreadsheet
	books, err := s.sheetsService.ReadBooksFromSpreadsheet(ctx, token, spreadsheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to read books: %w", err)
	}

	if len(books) == 0 {
		return nil, fmt.Errorf("no books found in spreadsheet")
	}

	// Generate AI recommendations
	aiResult, err := s.aiService.GenerateRecommendations(ctx, books)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI recommendations: %w", err)
	}

	// Create or get analysis
	var analysis models.SpreadsheetAnalysis
	err = s.db.Where("user_id = ? AND spreadsheet_id = ?", userID, spreadsheetID).First(&analysis).Error
	if err != nil {
		// Create new analysis if not found
		analysisPtr, err := s.AnalyzeSpreadsheet(ctx, userID, token, spreadsheetID)
		if err != nil {
			return nil, fmt.Errorf("failed to create analysis: %w", err)
		}
		analysis = *analysisPtr
	}

	// Convert AI result to recommendation models
	recommendations := []models.Recommendation{
		{
			UserID:        userID,
			AnalysisID:    analysis.ID,
			Title:         aiResult.Livro1.Titulo,
			Author:        aiResult.Livro1.Autor,
			Justification: aiResult.Livro1.Justificativa,
			Position:      1,
		},
		{
			UserID:        userID,
			AnalysisID:    analysis.ID,
			Title:         aiResult.Livro2.Titulo,
			Author:        aiResult.Livro2.Autor,
			Justification: aiResult.Livro2.Justificativa,
			Position:      2,
		},
		{
			UserID:        userID,
			AnalysisID:    analysis.ID,
			Title:         aiResult.Livro3.Titulo,
			Author:        aiResult.Livro3.Autor,
			Justification: aiResult.Livro3.Justificativa,
			Position:      3,
		},
	}

	// Save recommendations to database
	for i := range recommendations {
		if err := s.db.Create(&recommendations[i]).Error; err != nil {
			return nil, fmt.Errorf("failed to save recommendation %d: %w", i+1, err)
		}
	}

	response := &models.RecommendationResponse{
		Books:           books,
		Analysis:        analysis,
		Recommendations: recommendations,
	}

	return response, nil
}

func (s *BookService) GetRecommendationsHistory(ctx context.Context, userID uint) ([]models.Recommendation, error) {
	var recommendations []models.Recommendation

	err := s.db.Where("user_id = ?", userID).
		Preload("Analysis").
		Order("created_at DESC").
		Limit(50). // Limit to last 50 recommendations
		Find(&recommendations).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations history: %w", err)
	}

	return recommendations, nil
}

type UserPatterns struct {
	FavoriteGenres  []string
	FavoriteAuthors []string
	AverageRating   float64
	Summary         string
}

func (s *BookService) analyzeUserPatterns(books []models.Book) UserPatterns {
	if len(books) == 0 {
		return UserPatterns{}
	}

	// Count occurrences
	genreCounts := make(map[string]int)
	authorCounts := make(map[string]int)
	genreRatings := make(map[string][]int)
	totalRating := 0

	for _, book := range books {
		genreCounts[book.Genre]++
		authorCounts[book.Author]++
		genreRatings[book.Genre] = append(genreRatings[book.Genre], book.Rating)
		totalRating += book.Rating
	}

	// Find top genres by frequency and rating
	var favoriteGenres []string
	maxGenreCount := 0
	for genre, count := range genreCounts {
		if count > maxGenreCount {
			maxGenreCount = count
			favoriteGenres = []string{genre}
		} else if count == maxGenreCount {
			favoriteGenres = append(favoriteGenres, genre)
		}
	}

	// Limit to top 3 genres
	if len(favoriteGenres) > 3 {
		favoriteGenres = favoriteGenres[:3]
	}

	// Find top authors with multiple books
	var favoriteAuthors []string
	for author, count := range authorCounts {
		if count >= 2 { // Author must have at least 2 books
			favoriteAuthors = append(favoriteAuthors, author)
		}
	}

	// Limit to top 3 authors
	if len(favoriteAuthors) > 3 {
		favoriteAuthors = favoriteAuthors[:3]
	}

	averageRating := float64(totalRating) / float64(len(books))

	// Generate summary
	summary := fmt.Sprintf("Leitor com %d livros lidos, preferência por %s, nota média %.1f/5",
		len(books), favoriteGenres, averageRating)

	return UserPatterns{
		FavoriteGenres:  favoriteGenres,
		FavoriteAuthors: favoriteAuthors,
		AverageRating:   averageRating,
		Summary:         summary,
	}
}
