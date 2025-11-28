package ai

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

const generalPurposeModel = "gemini-2.5-flash"

type AIService struct {
	client *genai.Client
}

func NewAIService(client *genai.Client) *AIService {
	return &AIService{
		client: client,
	}
}

func (a *AIService) GenerateText(ctx context.Context, prompt string) (string, error) {
	contents := []*genai.Content{
		{
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	}
	result, err := a.client.Models.GenerateContent(ctx, generalPurposeModel, contents, nil)
	if err != nil {
		return "", fmt.Errorf("ai.AIService.GenerateText: failed to generate text: %w", err)
	}
	return result.ModelVersion, nil
}
