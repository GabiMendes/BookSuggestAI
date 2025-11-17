package services

import (
	"book-suggest-ai/internal/config"
	"book-suggest-ai/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIService struct {
	config *config.Config
	client *genai.Client
}

func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		config: cfg,
	}
}

func (s *AIService) initClient(ctx context.Context) error {
	if s.client != nil {
		return nil
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(s.config.GeminiAPIKey))
	if err != nil {
		return fmt.Errorf("failed to create Gemini client: %w", err)
	}

	s.client = client
	return nil
}

func (s *AIService) GenerateRecommendations(ctx context.Context, books []models.Book) (*models.AIRecommendationResult, error) {
	fmt.Printf("🤖 AI: Iniciando geração de recomendações para %d livros\n", len(books))

	if err := s.initClient(ctx); err != nil {
		fmt.Printf("❌ AI: Erro ao inicializar cliente: %v\n", err)
		return nil, err
	}
	fmt.Printf("✅ AI: Cliente Gemini inicializado com sucesso\n")

	// Analyze user preferences
	analysis := s.analyzeReadingPatterns(books)
	fmt.Printf("📊 AI: Análise concluída - Gêneros: %v, Autores: %v\n",
		analysis["favorite_genres"], analysis["favorite_authors"])

	// Create structured prompt
	prompt := s.createRecommendationPrompt(books, analysis)
	fmt.Printf("📝 AI: Prompt criado com %d caracteres\n", len(prompt))

	// List available models first
	availableModel, err := s.getAvailableModel(ctx)
	if err != nil {
		fmt.Printf("❌ AI: Erro ao listar modelos: %v\n", err)
		return nil, fmt.Errorf("failed to get available models: %w", err)
	}
	fmt.Printf("🎯 AI: Usando modelo disponível: %s\n", availableModel)

	// Get recommendations from Gemini
	model := s.client.GenerativeModel(availableModel)
	fmt.Printf("🔄 AI: Enviando requisição para Gemini com modelo %s...\n", availableModel)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		fmt.Printf("❌ AI: Erro na chamada Gemini: %v\n", err)
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}
	fmt.Printf("✅ AI: Resposta recebida do Gemini\n")

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		fmt.Printf("❌ AI: Resposta vazia do Gemini\n")
		return nil, fmt.Errorf("no recommendations received from AI")
	}

	// Extract text response
	responseText := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			responseText += string(textPart)
		}
	}
	fmt.Printf("📄 AI: Texto extraído (%d chars): %s\n", len(responseText), responseText[:min(200, len(responseText))]+"...")

	// Limpar resposta: remover markdown code blocks se presentes
	cleanText := responseText
	if strings.HasPrefix(cleanText, "```json") && strings.HasSuffix(cleanText, "```") {
		cleanText = strings.TrimPrefix(cleanText, "```json")
		cleanText = strings.TrimSuffix(cleanText, "```")
		cleanText = strings.TrimSpace(cleanText)
		fmt.Printf("🧹 AI: Removido markdown code block JSON\n")
	} else if strings.HasPrefix(cleanText, "```") && strings.HasSuffix(cleanText, "```") {
		cleanText = strings.TrimPrefix(cleanText, "```")
		cleanText = strings.TrimSuffix(cleanText, "```") 
		cleanText = strings.TrimSpace(cleanText)
		fmt.Printf("🧹 AI: Removido code block genérico\n")
	}

	// Parse JSON response
	var result models.AIRecommendationResult
	if err := json.Unmarshal([]byte(cleanText), &result); err != nil {
		fmt.Printf("❌ AI: Erro ao fazer parse JSON: %v\nResposta completa: %s\n", err, cleanText)
		return nil, fmt.Errorf("failed to parse AI response: %w\nResponse: %s", err, cleanText)
	}

	fmt.Printf("🎉 AI: Recomendações geradas com sucesso!\n")
	return &result, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *AIService) analyzeReadingPatterns(books []models.Book) map[string]interface{} {
	if len(books) == 0 {
		return map[string]interface{}{
			"total_books":      0,
			"average_rating":   0,
			"favorite_genres":  []string{},
			"favorite_authors": []string{},
		}
	}

	// Count genres and authors
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

	// Find favorite genres (by frequency and rating)
	type genreScore struct {
		genre string
		score float64
	}

	var genreScores []genreScore
	for genre, count := range genreCounts {
		ratings := genreRatings[genre]
		avgRating := 0.0
		for _, r := range ratings {
			avgRating += float64(r)
		}
		avgRating /= float64(len(ratings))

		// Score = frequency weight + rating weight
		score := float64(count)*0.5 + avgRating*0.5
		genreScores = append(genreScores, genreScore{genre, score})
	}

	sort.Slice(genreScores, func(i, j int) bool {
		return genreScores[i].score > genreScores[j].score
	})

	// Get top 3 genres
	favoriteGenres := []string{}
	for i, gs := range genreScores {
		if i >= 3 {
			break
		}
		favoriteGenres = append(favoriteGenres, gs.genre)
	}

	// Find favorite authors (by frequency)
	type authorScore struct {
		author string
		count  int
	}

	var authorScores []authorScore
	for author, count := range authorCounts {
		authorScores = append(authorScores, authorScore{author, count})
	}

	sort.Slice(authorScores, func(i, j int) bool {
		return authorScores[i].count > authorScores[j].count
	})

	// Get top 3 authors
	favoriteAuthors := []string{}
	for i, as := range authorScores {
		if i >= 3 || as.count < 2 { // Only include authors with 2+ books
			break
		}
		favoriteAuthors = append(favoriteAuthors, as.author)
	}

	averageRating := float64(totalRating) / float64(len(books))

	return map[string]interface{}{
		"total_books":      len(books),
		"average_rating":   averageRating,
		"favorite_genres":  favoriteGenres,
		"favorite_authors": favoriteAuthors,
	}
}

func (s *AIService) createRecommendationPrompt(books []models.Book, analysis map[string]interface{}) string {
	// Build book list for context
	bookList := strings.Builder{}
	for i, book := range books {
		if i >= 20 { // Limit to recent 20 books to avoid token limit
			break
		}
		bookList.WriteString(fmt.Sprintf("- \"%s\" por %s (Gênero: %s, Avaliação: %d/5)\n",
			book.Title, book.Author, book.Genre, book.Rating))
	}

	favoriteGenres := analysis["favorite_genres"].([]string)
	favoriteAuthors := analysis["favorite_authors"].([]string)
	averageRating := analysis["average_rating"].(float64)
	totalBooks := analysis["total_books"].(int)

	prompt := fmt.Sprintf(`Como especialista em literatura, analise o histórico de leitura a seguir e sugira exatamente 3 livros que este leitor provavelmente amaria.

HISTÓRICO DE LEITURA (%d livros, nota média: %.1f/5):
%s

PADRÕES IDENTIFICADOS:
- Gêneros preferidos: %s
- Autores favoritos: %s
- Avaliação média: %.1f/5

INSTRUÇÕES:
1. Analise os padrões de preferência considerando gêneros bem avaliados, autores recorrentes e temas
2. Sugira 3 livros diferentes que se alinhem com o perfil do leitor
3. Evite sugerir livros já lidos pelo usuário
4. Para cada sugestão, explique por que o leitor provavelmente gostaria
5. Retorne APENAS um JSON válido no formato exato:

{
  "livro1": {
    "titulo": "Nome do Livro 1",
    "autor": "Nome do Autor 1", 
    "justificativa": "Explicação detalhada de por que este livro combina com o perfil do leitor"
  },
  "livro2": {
    "titulo": "Nome do Livro 2",
    "autor": "Nome do Autor 2",
    "justificativa": "Explicação detalhada de por que este livro combina com o perfil do leitor"
  },
  "livro3": {
    "titulo": "Nome do Livro 3", 
    "autor": "Nome do Autor 3",
    "justificativa": "Explicação detalhada de por que este livro combina com o perfil do leitor"
  }
}`,
		totalBooks,
		averageRating,
		bookList.String(),
		strings.Join(favoriteGenres, ", "),
		strings.Join(favoriteAuthors, ", "),
		averageRating,
	)

	return prompt
}

func (s *AIService) getAvailableModel(ctx context.Context) (string, error) {
	fmt.Printf("🔍 AI: Listando modelos disponíveis...\n")
	
	models := s.client.ListModels(ctx)
	
	// Preferred models in order of preference
	preferredModels := []string{
		"gemini-1.5-pro",
		"gemini-1.5-flash", 
		"gemini-pro",
		"gemini-1.0-pro",
		"gemini-pro-latest",
	}
	
	// Get all available models
	var availableModels []string
	for {
		model, err := models.Next()
		if err != nil {
			break
		}
		fmt.Printf("📋 AI: Modelo encontrado: %s (suporta generateContent: %v)\n", 
			model.Name, containsMethod(model.SupportedGenerationMethods, "generateContent"))
		
		// Only consider models that support generateContent
		if containsMethod(model.SupportedGenerationMethods, "generateContent") {
			availableModels = append(availableModels, model.Name)
		}
	}
	
	if len(availableModels) == 0 {
		return "", fmt.Errorf("no models available that support generateContent")
	}
	
	fmt.Printf("✅ AI: Modelos disponíveis que suportam generateContent: %v\n", availableModels)
	
	// Try to find preferred model
	for _, preferred := range preferredModels {
		for _, available := range availableModels {
			if strings.Contains(available, preferred) || available == preferred {
				fmt.Printf("🎯 AI: Selecionado modelo preferido: %s\n", available)
				return available, nil
			}
		}
	}
	
	// If no preferred model found, use the first available
	selected := availableModels[0]
	fmt.Printf("🎯 AI: Usando primeiro modelo disponível: %s\n", selected)
	return selected, nil
}

func containsMethod(methods []string, method string) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

func (s *AIService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
