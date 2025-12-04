package interfaces

import (
	"context"

	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type RagStore interface {
	StoreDocumentEmbeddings(ctx context.Context, documentID types.DocumentID, embeddings [][]float32) error
}

type RagEmbedder interface {
	GetDocumentEmbeddings(ctx context.Context, documentPath string) ([][]float32, error)
}

type RagRetriever interface {
	RetrieveDocumentContent(ctx context.Context, documentID types.DocumentID, query string, topK int) ([]string, error)
}
