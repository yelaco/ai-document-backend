package rag

import (
	"context"
	"fmt"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	chromaEmbeddings "github.com/amikos-tech/chroma-go/pkg/embeddings"
)

type Store interface {
	StoreDocumentEmbeddings(ctx context.Context, embedding [][]float32) error
}

type ChromaStore struct {
	client chroma.Client
}

func NewChromaStore(client chroma.Client) Store {
	return &ChromaStore{
		client: client,
	}
}

func (s *ChromaStore) StoreDocumentEmbeddings(ctx context.Context, embeddings [][]float32) error {
	collection, err := s.client.GetOrCreateCollection(ctx, "documents")
	if err != nil {
		return fmt.Errorf("ChromaStore.StoreDocumentEmbeddings: failed to get or create collection: %w", err)
	}

	ebds, err := chromaEmbeddings.NewEmbeddingsFromFloat32(embeddings)
	if err != nil {
		return fmt.Errorf("ChromaStore.StoreDocumentEmbeddings: failed to create chroma embeddings: %w", err)
	}
	err = collection.Add(ctx, chroma.WithIDs("1"), chroma.WithEmbeddings(ebds...))
	if err != nil {
		return fmt.Errorf("ChromaStore.StoreDocumentEmbeddings: failed to add embeddings: %w", err)
	}
	return nil
}
