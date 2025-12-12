package chatbot

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
	ctx    context.Context
}

func NewGeminiClient() (*GeminiClient, error) {
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}
	return &GeminiClient{client: client, ctx: ctx}, nil
}

func (g *GeminiClient) SendMessage(message string) (string, error) {
	res, err := g.client.Models.GenerateContent(g.ctx, "gemini-2.5-flash", genai.Text(message), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}
	return res.Text(), nil
}
