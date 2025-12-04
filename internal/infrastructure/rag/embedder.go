package rag

import (
	"context"
	"fmt"

	"github.com/amikos-tech/chroma-go/pkg/embeddings"
	g "github.com/amikos-tech/chroma-go/pkg/embeddings/gemini"
	"github.com/samber/lo"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
)

const (
	geminiEmbeddingModel = "gemini-embedding-001"
)

type ChromaEmbedder struct {
	embedFunc *g.GeminiEmbeddingFunction
}

func NewChromaEmbedder(apiKey string) interfaces.RagEmbedder {
	geminiEmbedFunc, err := g.NewGeminiEmbeddingFunction(
		g.WithAPIKey(apiKey),
		g.WithDefaultModel(embeddings.EmbeddingModel(geminiEmbeddingModel)),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create Gemini embedding function: %v", err))
	}
	return &ChromaEmbedder{
		embedFunc: geminiEmbedFunc,
	}
}

func (e *ChromaEmbedder) GetDocumentEmbeddings(ctx context.Context, documentPath string) ([][]float32, error) {
	// This is a placeholder implementation - you'll need to implement document reading
	// and chunking based on the documentPath
	documents := []string{"placeholder document content"} // TODO: implement document reading

	result, err := e.embedFunc.EmbedDocuments(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("ChromaEmbedder.GetDocumentEmbeddings: failed to embed documents: %w", err)
	}
	return lo.Map(result, func(e embeddings.Embedding, _ int) []float32 {
		return e.ContentAsFloat32()
	}), nil
}
