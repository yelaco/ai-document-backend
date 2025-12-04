package document

import (
	"context"
	"fmt"

	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/dtos"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
	reqContext "github.com/yelaco/ai-document-backend/internal/infrastructure/context"
)

type DocumentService struct {
	documentRepo interfaces.DocumentRepository
}

func NewDocumentService(documentRepo interfaces.DocumentRepository) interfaces.DocumentService {
	return &DocumentService{
		documentRepo: documentRepo,
	}
}

// CreateDocument implements interfaces.DocumentService.
func (d *DocumentService) CreateDocument(ctx context.Context, params dtos.CreateDocumentParams) (entity.Document, error) {
	userID := reqContext.UserIDMustFromContext(ctx)
	document := entity.Document{
		UserID:       userID,
		OriginalName: params.OriginalName,
		SavePath:     params.SavePath,
	}
	if err := d.documentRepo.CreateDocument(ctx, &document); err != nil {
		return entity.Document{}, fmt.Errorf("DocumentService.CreateDocument: failed to create document: %w", err)
	}
	return document, nil
}

// ListDocumentsByUserID implements interfaces.DocumentService.
func (d *DocumentService) GetPaginatedDocuments(ctx context.Context, page int64, pageSize int64) ([]entity.Document, int64, error) {
	userID := reqContext.UserIDMustFromContext(ctx)
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	documents, count, err := d.documentRepo.GetPaginatedDocuments(ctx, userID, skip, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("DocumentService.GetPaginatedDocuments: failed to get documents: %w", err)
	}
	return documents, count, nil
}

// GetDocumentByID implements interfaces.DocumentService.
func (d *DocumentService) GetDocumentByID(ctx context.Context, documentID types.DocumentID) (entity.Document, error) {
	userID := reqContext.UserIDMustFromContext(ctx)
	document, err := d.documentRepo.GetDocument(ctx, documentID, userID)
	if err != nil {
		return entity.Document{}, fmt.Errorf("DocumentService.GetDocumentByID: failed to get document: %w", err)
	}
	return document, nil
}

// DeleteDocument implements interfaces.DocumentService.
func (d *DocumentService) DeleteDocument(ctx context.Context, documentID types.DocumentID) error {
	userID := reqContext.UserIDMustFromContext(ctx)
	err := d.documentRepo.DeleteDocument(ctx, documentID, userID)
	if err != nil {
		return fmt.Errorf("DocumentService.DeleteDocument: failed to delete document: %w", err)
	}
	return nil
}
