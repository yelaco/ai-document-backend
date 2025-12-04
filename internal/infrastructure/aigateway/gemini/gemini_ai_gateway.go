package gemini

import (
	"context"
	"fmt"

	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"google.golang.org/genai"
)

const geminiGenerativeModel = "gemini-2.5-flash"

type GeminiAIGateway struct {
	client *genai.Client
}

func NewGeminiAIGateway(client *genai.Client) interfaces.AIGateway {
	return &GeminiAIGateway{
		client: client,
	}
}

// GenerateImage implements interfaces.AIGateway.
func (g *GeminiAIGateway) GenerateImage(ctx context.Context, description string) ([]byte, error) {
	panic("unimplemented")
}

// GenerateText implements interfaces.AIGateway.
func (g *GeminiAIGateway) GenerateText(ctx context.Context, prompt string) (string, error) {
	panic("unimplemented")
}

func (g *GeminiAIGateway) GenerateTextStream(ctx context.Context, prompt string) (string, error) {
	parts := []*genai.Part{
		{Text: "What's this image about?"},
	}
	iter := g.client.Models.GenerateContentStream(ctx, geminiGenerativeModel, []*genai.Content{{Parts: parts}}, nil)
	for resp, err := range iter {
		if err != nil {
			continue
		}
		fmt.Println(resp.Text())
	}
	return "", nil
}
