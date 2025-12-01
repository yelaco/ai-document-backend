package interfaces

import "context"

type AIGateway interface {
	GenerateText(ctx context.Context, prompt string) (string, error)
	GenerateImage(ctx context.Context, description string) ([]byte, error)
}
