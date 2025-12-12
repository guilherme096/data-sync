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

func (g *GeminiClient) SendMessageWithHistory(message string, history []ChatMessage) (string, error) {
	// Convert history to Gemini format
	var contents []*genai.Content

	for _, msg := range history {
		// Map frontend roles to Gemini roles
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}

		// Create content for this message with the correct role
		msgContents := genai.Text(msg.Content)
		if len(msgContents) > 0 {
			msgContents[0].Role = role
			contents = append(contents, msgContents...)
		}
	}

	// Add current user message
	userContents := genai.Text(message)
	if len(userContents) > 0 {
		userContents[0].Role = "user"
		contents = append(contents, userContents...)
	}

	// Generate response with conversation context
	res, err := g.client.Models.GenerateContent(g.ctx, "gemini-2.5-flash", contents, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate content with history: %w", err)
	}

	return res.Text(), nil
}
