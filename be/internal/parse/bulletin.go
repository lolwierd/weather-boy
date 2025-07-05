package parse

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/lolwierd/weatherboy/be/internal/config"
	"github.com/sashabaranov/go-openai"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// ParseBulletinPDF extracts text from a bulletin PDF and uses the Gemini API
// to get a structured forecast for a given city.
func ParseBulletinPDF(ctx context.Context, pdfPath string, city string) (string, error) {
	f, err := os.Open(pdfPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return "", err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	var text string
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			return "", err
		}

		ex, err := extractor.New(page)
		if err != nil {
			return "", err
		}

		pageText, err := ex.ExtractText()
		if err != nil {
			return "", err
		}
		text += pageText
	}

	config.LoadEnv()
	// Extract a relevant snippet from the full text to reduce token usage and improve focus.
	snippet := extractForecastSnippet(text, city)

	client := openai.NewClient(config.OpenAIAPIKey)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: snippet, // Use the extracted snippet here
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// extractForecastSnippet attempts to find the forecast for a specific city within the given text.
// It returns a snippet of text around the city name.
func extractForecastSnippet(fullText, city string) string {
	// Define a window size for the snippet (characters before and after the city name)
	const snippetWindowSize = 400

	cityIdx := -1
	// Search for the city name case-insensitively
	lowerFullText := strings.ToLower(fullText)
	lowerCity := strings.ToLower(city)

	cityIdx = strings.Index(lowerFullText, lowerCity)

	if cityIdx == -1 {
		// If city not found, return the first few lines or a default snippet
		lines := strings.Split(fullText, "\n")
		if len(lines) > 5 {
			return strings.Join(lines[:5], "\n") + "\n..." // Return first 5 lines as a fallback
		}
		return fullText // Fallback to full text if very short
	}

	start := cityIdx - snippetWindowSize
	if start < 0 {
		start = 0
	}

	end := cityIdx + len(city) + snippetWindowSize
	if end > len(fullText) {
		end = len(fullText)
	}

	// Adjust start to the beginning of a line and end to the end of a line if possible
	// to avoid cutting words in half.
	for start > 0 && fullText[start-1] != '\n' && fullText[start-1] != '\r' {
		start--
	}
	for end < len(fullText) && fullText[end] != '\n' && fullText[end] != '\r' {
		end++
	}

	snippet := fullText[start:end]

	// Add context to the prompt to guide the LLM
	return fmt.Sprintf("Here is a snippet from a weather bulletin. Extract the forecast specifically for %s from this text:\n\n%s", city, snippet)
}
