package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GoogleID    string    `gorm:"uniqueIndex" json:"google_id"`
	Email       string    `gorm:"uniqueIndex" json:"email"`
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar"`
	AccessToken string    `json:"-"` // Not included in JSON
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Book represents a book from the spreadsheet
type Book struct {
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Genre        string    `json:"genre"`
	ReadDate     time.Time `json:"read_date"`
	Rating       int       `json:"rating"` // 1-5 stars
	Observations string    `json:"observations"`
}

// SpreadsheetAnalysis represents the analysis of a user's reading data
type SpreadsheetAnalysis struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `json:"user_id"`
	SpreadsheetID   string    `json:"spreadsheet_id"`
	TotalBooks      int       `json:"total_books"`
	FavoriteGenres  []string  `gorm:"serializer:json" json:"favorite_genres"`
	FavoriteAuthors []string  `gorm:"serializer:json" json:"favorite_authors"`
	AverageRating   float64   `json:"average_rating"`
	ReadingPatterns string    `gorm:"type:text" json:"reading_patterns"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Recommendation represents a book recommendation
type Recommendation struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `json:"user_id"`
	AnalysisID    uint      `json:"analysis_id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Justification string    `gorm:"type:text" json:"justification"`
	Position      int       `json:"position"`      // 1, 2, or 3
	UserFeedback  *int      `json:"user_feedback"` // 1-5 rating or null
	CreatedAt     time.Time `json:"created_at"`

	User     User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Analysis SpreadsheetAnalysis `gorm:"foreignKey:AnalysisID" json:"analysis,omitempty"`
}

// RecommendationRequest represents the request for generating recommendations
type RecommendationRequest struct {
	SpreadsheetID string `json:"spreadsheet_id" binding:"required"`
}

// RecommendationResponse represents the response with recommendations
type RecommendationResponse struct {
	Books           []Book              `json:"books"`
	Analysis        SpreadsheetAnalysis `json:"analysis"`
	Recommendations []Recommendation    `json:"recommendations"`
}

// GoogleSpreadsheet represents a Google Sheets file
type GoogleSpreadsheet struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// AIRecommendationData represents the data sent to AI for analysis
type AIRecommendationData struct {
	Books    []Book `json:"books"`
	UserInfo struct {
		TotalBooks      int      `json:"total_books"`
		AverageRating   float64  `json:"average_rating"`
		FavoriteGenres  []string `json:"favorite_genres"`
		FavoriteAuthors []string `json:"favorite_authors"`
	} `json:"user_info"`
}

// AIRecommendationResult represents the AI response
type AIRecommendationResult struct {
	Livro1 struct {
		Titulo        string `json:"titulo"`
		Autor         string `json:"autor"`
		Justificativa string `json:"justificativa"`
	} `json:"livro1"`
	Livro2 struct {
		Titulo        string `json:"titulo"`
		Autor         string `json:"autor"`
		Justificativa string `json:"justificativa"`
	} `json:"livro2"`
	Livro3 struct {
		Titulo        string `json:"titulo"`
		Autor         string `json:"autor"`
		Justificativa string `json:"justificativa"`
	} `json:"livro3"`
}
