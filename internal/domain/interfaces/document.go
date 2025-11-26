package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type DocumentService interface {
	ProcessDocument(ctx context.Context, documentID string) (string, error)
	CreateDocument(ctx context.Context, title string, userID uuid.UUID) (entity.Document, error)
}

type DocumentRepository interface {
	CreateDocument(ctx context.Context, document *entity.Document) error
	GetDocumentByID(ctx context.Context, documentID string) (string, error)
	GetDocumentsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Document, error)
	DeleteDocumentByID(ctx context.Context, documentID string) error
}
