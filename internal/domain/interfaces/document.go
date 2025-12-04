package interfaces

import (
	"context"

	"github.com/yelaco/ai-document-backend/internal/domain/models/dtos"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type DocumentService interface {
	CreateDocument(ctx context.Context, params dtos.CreateDocumentParams) (entity.Document, error)
	GetPaginatedDocuments(ctx context.Context, page int64, pageSize int64) ([]entity.Document, int64, error)
	GetDocumentByID(ctx context.Context, documentID types.DocumentID) (entity.Document, error)
	DeleteDocument(ctx context.Context, documentID types.DocumentID) error
}

type DocumentRepository interface {
	CreateDocument(ctx context.Context, document *entity.Document) error
	GetDocument(ctx context.Context, documentID types.DocumentID, userID types.UserID) (entity.Document, error)
	GetPaginatedDocuments(ctx context.Context, userID types.UserID, skip int64, limit int64) ([]entity.Document, int64, error)
	DeleteDocument(ctx context.Context, documentID types.DocumentID, userID types.UserID) error
}
