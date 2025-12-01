package gemini

import (
	"context"

	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
)

type GeminiAIGateway struct{}

func NewGeminiAIGateway() interfaces.AIGateway {
	return &GeminiAIGateway{}
}

// GenerateImage implements interfaces.AIGateway.
func (g *GeminiAIGateway) GenerateImage(ctx context.Context, description string) ([]byte, error) {
	panic("unimplemented")
}

// GenerateText implements interfaces.AIGateway.
func (g *GeminiAIGateway) GenerateText(ctx context.Context, prompt string) (string, error) {
	panic("unimplemented")
}
