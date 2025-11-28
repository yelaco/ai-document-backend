package dtos

import (
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type PaginatedDocumentsParams struct {
	PaginationParams
}

type DocumentResponse struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	OriginalName string `json:"title"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func DocumentResponseFromEntity(doc entity.Document) DocumentResponse {
	return DocumentResponse{
		ID:           doc.ID.String(),
		UserID:       doc.UserID.String(),
		OriginalName: doc.OriginalName,
		CreatedAt:    doc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    doc.CreatedAt.Format(time.RFC3339),
	}
}

func DocumentResponsesFromEntities(docs []entity.Document) []DocumentResponse {
	responses := make([]DocumentResponse, len(docs))
	for i, doc := range docs {
		responses[i] = DocumentResponseFromEntity(doc)
	}
	return responses
}
