package rag

import (
	"context"
	"fmt"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/amikos-tech/chroma-go/pkg/embeddings"
	g "github.com/amikos-tech/chroma-go/pkg/embeddings/gemini"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
)

type ChromaRetriever struct {
	client    chroma.Client
	embedFunc *g.GeminiEmbeddingFunction
}

func NewChromaRetriever(client chroma.Client, apiKey string) interfaces.RagRetriever {
	geminiEmbedFunc, err := g.NewGeminiEmbeddingFunction(
		g.WithAPIKey(apiKey),
		g.WithDefaultModel(embeddings.EmbeddingModel(geminiEmbeddingModel)),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create Gemini embedding function: %v", err))
	}
	return &ChromaRetriever{
		client:    client,
		embedFunc: geminiEmbedFunc,
	}
}

func (c *ChromaRetriever) RetrieveDocumentContent(ctx context.Context, documentID uuid.UUID, query string, topK int) ([]string, error) {
	collection, err := c.client.GetOrCreateCollection(ctx, "documents",
		chroma.WithEmbeddingFunctionCreate(c.embedFunc),
	)
	if err != nil {
		return nil, fmt.Errorf("ChromaStore.StoreDocumentEmbeddings: failed to get or create collection: %w", err)
	}

	result, err := collection.Query(
		ctx,
		chroma.WithQueryTexts(query),
		chroma.WithNResults(topK),
		chroma.WithWhereQuery(chroma.EqString("document_id", documentID.String())),
	)
	if err != nil {
		return nil, fmt.Errorf("ChromaRetriever.RetrieveDocumentContent: failed to query collection: %w", err)
	}
	return lo.Map(result.ToRecordsGroups()[0], func(doc chroma.Record, _ int) string {
		return doc.Document().ContentString()
	}), nil
}
