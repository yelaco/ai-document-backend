package document

import (
	"context"

	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type DocumentService struct {
	documentRepo interfaces.DocumentRepository
}

// CreateDocument implements interfaces.DocumentService.
func (d *DocumentService) CreateDocument(ctx context.Context, title string, userID uuid.UUID) (entity.Document, error) {
	panic("unimplemented")
}

// ProcessDocument implements interfaces.DocumentService.
func (d *DocumentService) ProcessDocument(ctx context.Context, documentID string) (string, error) {
	panic("unimplemented")
}

func NewDocumentService(documentRepo interfaces.DocumentRepository) interfaces.DocumentService {
	return &DocumentService{
		documentRepo: documentRepo,
	}
}
