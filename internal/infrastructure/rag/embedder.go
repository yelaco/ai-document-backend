package rag

import (
	"context"
	"fmt"

	"github.com/amikos-tech/chroma-go/pkg/embeddings"
	g "github.com/amikos-tech/chroma-go/pkg/embeddings/gemini"
	"github.com/samber/lo"
)

const (
	geminiEmbeddingModel = "gemini-embedding-001"
)

type Embedder interface {
	Embed(ctx context.Context, documents []string) ([][]float32, error)
}

type ChromaEmbedder struct {
	embedder *g.GeminiEmbeddingFunction
}

func NewChromaEmbedder(apiKey string) Embedder {
	geminiEmbedder, err := g.NewGeminiEmbeddingFunction(
		g.WithAPIKey(apiKey),
		g.WithDefaultModel(embeddings.EmbeddingModel(geminiEmbeddingModel)),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create Gemini embedding function: %v", err))
	}
	return &ChromaEmbedder{
		embedder: geminiEmbedder,
	}
}

func (e *ChromaEmbedder) Embed(ctx context.Context, documents []string) ([][]float32, error) {
	result, err := e.embedder.EmbedDocuments(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("ChromaGeminiEmbedder.Embed: failed to embed documents: %w", err)
	}
	return lo.Map(result, func(e embeddings.Embedding, _ int) []float32 {
		return e.ContentAsFloat32()
	}), nil
}
