package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type RagStore interface {
	StoreDocumentEmbeddings(ctx context.Context, documentID uuid.UUID, embeddings [][]float32) error
}

type RagEmbedder interface {
	Embed(ctx context.Context, documents []string) ([][]float32, error)
}

type RagRetriever interface {
	RetrieveDocumentContent(ctx context.Context, documentID uuid.UUID, query string, topK int) ([]string, error)
}
