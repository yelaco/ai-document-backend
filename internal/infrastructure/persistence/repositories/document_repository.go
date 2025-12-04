package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/persistence/database/sqlc"
)

type PostgresDocumentRepository struct {
	connPool *pgxpool.Pool
	queries  *sqlc.Queries
}

func NewDocumentRepository(connPool *pgxpool.Pool) interfaces.DocumentRepository {
	return &PostgresDocumentRepository{
		connPool: connPool,
		queries:  sqlc.New(connPool),
	}
}

// CreateDocument implements interfaces.DocumentRepository.
func (p *PostgresDocumentRepository) CreateDocument(ctx context.Context, document *entity.Document) error {
	row, err := p.queries.CreateDocument(ctx, sqlc.CreateDocumentParams{
		OriginalName: document.OriginalName,
		SavePath:     document.SavePath,
		UserID:       document.UserID.UUID(),
	})
	if err != nil {
		return fmt.Errorf("PostgresDocumentRepository.CreateDocument: failed to create document: %w", err)
	}
	document.ID = types.NewDocumentID(row.ID)
	document.CreatedAt = row.CreatedAt
	document.UpdatedAt = row.UpdatedAt
	return nil
}

// GetDocumentByID implements interfaces.DocumentRepository.
func (p *PostgresDocumentRepository) GetDocument(ctx context.Context, documentID types.DocumentID, userID types.UserID) (entity.Document, error) {
	row, err := p.queries.GetDocument(ctx, sqlc.GetDocumentParams{
		ID:     documentID.UUID(),
		UserID: userID.UUID(),
	})
	if err != nil {
		return entity.Document{}, fmt.Errorf("PostgresDocumentRepository.GetDocument: failed to get document: %w", err)
	}
	return entity.Document{
		ID:           types.NewDocumentID(row.ID),
		OriginalName: row.OriginalName,
		SavePath:     row.SavePath,
		UserID:       types.NewUserID(row.UserID),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}

// GetDocumentsByUserID implements interfaces.DocumentRepository.
func (p *PostgresDocumentRepository) GetPaginatedDocuments(ctx context.Context, userID types.UserID, skip int64, limit int64) ([]entity.Document, int64, error) {
	rows, err := p.queries.GetPaginatedDocuments(ctx, sqlc.GetPaginatedDocumentsParams{
		UserID: userID.UUID(),
		Limit:  int32(limit),
		Offset: int32(skip),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("PostgresDocumentRepository.GetPaginatedDocuments: failed to get documents: %w", err)
	}
	count, err := p.queries.CountUserDocuments(ctx, userID.UUID())
	if err != nil {
		return nil, 0, fmt.Errorf("PostgresDocumentRepository.GetPaginatedDocuments: failed to count documents: %w", err)
	}
	documents := lo.Map(rows, func(row sqlc.GetPaginatedDocumentsRow, _ int) entity.Document {
		return entity.Document{
			ID:           types.NewDocumentID(row.ID),
			OriginalName: row.OriginalName,
			SavePath:     row.SavePath,
			Status:       row.Status,
			UserID:       types.NewUserID(row.UserID),
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
		}
	})

	return documents, count, nil
}

// DeleteDocumentByID implements interfaces.DocumentRepository.
func (p *PostgresDocumentRepository) DeleteDocument(ctx context.Context, documentID types.DocumentID, userID types.UserID) error {
	err := p.queries.DeleteDocument(ctx, sqlc.DeleteDocumentParams{
		ID:     documentID.UUID(),
		UserID: userID.UUID(),
	})
	if err != nil {
		return fmt.Errorf("PostgresDocumentRepository.DeleteDocument: failed to delete document: %w", err)
	}
	return nil
}
