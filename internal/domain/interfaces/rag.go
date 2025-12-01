package interfaces

import "context"

type RagStore interface {
	StoreDocumentEmbeddings(ctx context.Context, embeddings [][]float32) error
}

type RagEmbedder interface {
	Embed(ctx context.Context, documents []string) ([][]float32, error)
}

type RagRetriever interface {
	Retrieve(ctx context.Context, query string, topK int) ([]string, error)
}
