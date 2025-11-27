package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
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
		Title:  document.Title,
		UserID: document.UserID,
	})
	if err != nil {
		return fmt.Errorf("PostgresDocumentRepository.CreateDocument: failed to create document: %w", err)
	}
	document.ID = row.ID
	document.CreatedAt = row.CreatedAt
	document.UpdatedAt = row.UpdatedAt
	return nil
}

// GetDocumentByID implements interfaces.DocumentRepository.
func (p *PostgresDocumentRepository) GetDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) (entity.Document, error) {
	row, err := p.queries.GetDocument(ctx, sqlc.GetDocumentParams{
		ID:     documentID,
		UserID: userID,
	})
	if err != nil {
		return entity.Document{}, fmt.Errorf("PostgresDocumentRepository.GetDocument: failed to get document: %w", err)
	}
	return entity.Document{
		ID:        row.ID,
		Title:     row.Title,
		UserID:    row.UserID,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// GetDocumentsByUserID implements interfaces.DocumentRepository.
func (p *PostgresDocumentRepository) GetPaginatedDocuments(ctx context.Context, userID uuid.UUID, skip int64, limit int64) ([]entity.Document, int64, error) {
	rows, err := p.queries.GetPaginatedDocuments(ctx, sqlc.GetPaginatedDocumentsParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(skip),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("PostgresDocumentRepository.GetPaginatedDocuments: failed to get documents: %w", err)
	}
	count, err := p.queries.CountUserDocuments(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("PostgresDocumentRepository.GetPaginatedDocuments: failed to count documents: %w", err)
	}
	documents := lo.Map(rows, func(row sqlc.GetPaginatedDocumentsRow, _ int) entity.Document {
		return entity.Document{
			ID:        row.ID,
			Title:     row.Title,
			Status:    row.Status,
			UserID:    row.UserID,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}
	})

	return documents, count, nil
}

// DeleteDocumentByID implements interfaces.DocumentRepository.
func (p *PostgresDocumentRepository) DeleteDocument(ctx context.Context, documentID uuid.UUID, userID uuid.UUID) error {
	err := p.queries.DeleteDocument(ctx, sqlc.DeleteDocumentParams{
		ID:     documentID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("PostgresDocumentRepository.DeleteDocument: failed to delete document: %w", err)
	}
	return nil
}
