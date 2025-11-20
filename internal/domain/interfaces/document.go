package interfaces

import (
	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type DocumentService interface {
	GenerateDocumentSummary(documentID string) (string, error)
}

type DocumentRepository interface {
	CreateDocument(title string, userID uuid.UUID) (entity.Document, error)
	GetDocumentByID(documentID string) (string, error)
	GetDocumentsByUserID(userID uuid.UUID) ([]entity.Document, error)
	DeleteDocumentByID(documentID string) error
}
