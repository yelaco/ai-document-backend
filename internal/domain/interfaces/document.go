package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type DocumentService interface {
	CreateDocument(ctx context.Context, title string, savePath string) (entity.Document, error)
	GetPaginatedDocuments(ctx context.Context, page int64, pageSize int64) ([]entity.Document, int64, error)
	GetDocumentByID(ctx context.Context, documentID uuid.UUID) (entity.Document, error)
	DeleteDocument(ctx context.Context, documentID uuid.UUID) error
}

type DocumentRepository interface {
	CreateDocument(ctx context.Context, document *entity.Document) error
	GetDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) (entity.Document, error)
	GetPaginatedDocuments(ctx context.Context, userID uuid.UUID, skip int64, limit int64) ([]entity.Document, int64, error)
	DeleteDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error
}
