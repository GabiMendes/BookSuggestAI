package services

import (
	"book-suggest-ai/internal/config"
	"book-suggest-ai/internal/models"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsService struct {
	config *config.Config
}

func NewSheetsService(cfg *config.Config) *SheetsService {
	return &SheetsService{
		config: cfg,
	}
}

func (s *SheetsService) ListSpreadsheets(ctx context.Context, token *oauth2.Token) ([]models.GoogleSpreadsheet, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Drive service: %w", err)
	}

	query := "mimeType='application/vnd.google-apps.spreadsheet'"
	fileList, err := driveService.Files.List().Q(query).Fields("files(id,name,webViewLink)").Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list spreadsheets: %w", err)
	}

	var spreadsheets []models.GoogleSpreadsheet
	for _, file := range fileList.Files {
		spreadsheets = append(spreadsheets, models.GoogleSpreadsheet{
			ID:   file.Id,
			Name: file.Name,
			URL:  file.WebViewLink,
		})
	}

	return spreadsheets, nil
}

func (s *SheetsService) ReadBooksFromSpreadsheet(ctx context.Context, token *oauth2.Token, spreadsheetID string) ([]models.Book, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	sheetsService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Sheets service: %w", err)
	}

	// Read the first sheet, starting from row 2 (assuming row 1 has headers)
	// Expected columns: A=Título, B=Autor, C=Gênero, D=Data Leitura, E=Avaliação, F=Observações
	readRange := "A2:F1000" // Read up to 1000 rows

	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to read spreadsheet: %w", err)
	}

	var books []models.Book
	for i, row := range resp.Values {
		// Skip empty rows or rows with less than 5 columns (Title, Author, Genre, Date, Rating are required)
		if len(row) < 5 {
			continue
		}

		book := models.Book{}

		// Title (Column A)
		if len(row) > 0 {
			book.Title = strings.TrimSpace(fmt.Sprintf("%v", row[0]))
		}

		// Author (Column B)
		if len(row) > 1 {
			book.Author = strings.TrimSpace(fmt.Sprintf("%v", row[1]))
		}

		// Genre (Column C)
		if len(row) > 2 {
			book.Genre = strings.TrimSpace(fmt.Sprintf("%v", row[2]))
		}

		// Date (Column D)
		if len(row) > 3 {
			dateStr := strings.TrimSpace(fmt.Sprintf("%v", row[3]))
			if parsedDate, err := s.parseDate(dateStr); err == nil {
				book.ReadDate = parsedDate
			} else {
				// If date parsing fails, use current time as fallback
				book.ReadDate = time.Now()
			}
		}

		// Rating (Column E)
		if len(row) > 4 {
			ratingStr := strings.TrimSpace(fmt.Sprintf("%v", row[4]))
			if rating, err := strconv.Atoi(ratingStr); err == nil && rating >= 1 && rating <= 5 {
				book.Rating = rating
			} else {
				// Skip books without valid ratings
				continue
			}
		}

		// Observations (Column F) - Optional
		if len(row) > 5 {
			book.Observations = strings.TrimSpace(fmt.Sprintf("%v", row[5]))
		}

		// Validate required fields
		if book.Title == "" || book.Author == "" || book.Genre == "" {
			fmt.Printf("Skipping row %d: missing required fields (title, author, or genre)\n", i+2)
			continue
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *SheetsService) parseDate(dateStr string) (time.Time, error) {
	// Try different date formats
	formats := []string{
		"02/01/2006", // DD/MM/YYYY
		"2006-01-02", // YYYY-MM-DD
		"01/02/2006", // MM/DD/YYYY
		"02-01-2006", // DD-MM-YYYY
		"2006/01/02", // YYYY/MM/DD
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func (s *SheetsService) ValidateSpreadsheetStructure(ctx context.Context, token *oauth2.Token, spreadsheetID string) error {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	sheetsService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("failed to create Sheets service: %w", err)
	}

	// Read the header row
	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, "A1:F1").Do()
	if err != nil {
		return fmt.Errorf("failed to read spreadsheet headers: %w", err)
	}

	if len(resp.Values) == 0 || len(resp.Values[0]) < 5 {
		return fmt.Errorf("spreadsheet must have at least 5 columns: Título, Autor, Gênero, Data Leitura, Avaliação")
	}

	// Optional: Check if headers match expected format
	// This is flexible as users might have different column names

	return nil
}
