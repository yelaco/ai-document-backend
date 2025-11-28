package rag

import (
	"context"

	chroma "github.com/amikos-tech/chroma-go"
)

type Store interface {
	StoreDocumentEmbeddings(ctx context.Context, embeddings [][]float32) error
}

type ChromaStore struct {
	client *chroma.Client
}

func NewChromaStore(client *chroma.Client) Store {
	return &ChromaStore{
		client: client,
	}
}

func (s *ChromaStore) StoreDocumentEmbeddings(ctx context.Context, embeddings [][]float32) error {
	return nil
}
