package rag

import (
	"context"
	"fmt"
	"time"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	chromaEmbeddings "github.com/amikos-tech/chroma-go/pkg/embeddings"
	"github.com/samber/lo"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type ChromaStore struct {
	client chroma.Client
}

func NewChromaStore(client chroma.Client) interfaces.RagStore {
	return &ChromaStore{
		client: client,
	}
}

func (s *ChromaStore) StoreDocumentEmbeddings(ctx context.Context, documentID types.DocumentID, embeddings [][]float32) error {
	collection, err := s.client.GetOrCreateCollection(ctx, "documents",
		// We won't use this embedding function for actual embedding,
		chroma.WithEmbeddingFunctionCreate(&chromaEmbeddings.ConsistentHashEmbeddingFunction{}),
	)
	if err != nil {
		return fmt.Errorf("ChromaStore.StoreDocumentEmbeddings: failed to get or create collection: %w", err)
	}

	ebds, err := chromaEmbeddings.NewEmbeddingsFromFloat32(embeddings)
	if err != nil {
		return fmt.Errorf("ChromaStore.StoreDocumentEmbeddings: failed to create chroma embeddings: %w", err)
	}
	err = collection.Upsert(
		ctx,
		chroma.WithIDs(lo.Map(ebds, func(_ chromaEmbeddings.Embedding, i int) chroma.DocumentID {
			return chroma.DocumentID(fmt.Sprintf("%s_%d_%d", documentID.String(), time.Now().UnixNano(), i))
		})...),
		chroma.WithEmbeddings(ebds...),
		chroma.WithMetadatas(
			lo.Map(nil, func(_ struct{}, _ int) chroma.DocumentMetadata {
				return chroma.NewDocumentMetadata(
					chroma.NewStringAttribute("document_id", documentID.String()),
					chroma.NewIntAttribute("created_at", time.Now().Unix()),
				)
			})...,
		),
	)
	if err != nil {
		return fmt.Errorf("ChromaStore.StoreDocumentEmbeddings: failed to add embeddings: %w", err)
	}
	return nil
}
