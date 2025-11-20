package repositories

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
func (p *PostgresDocumentRepository) CreateDocument(title string, userID uuid.UUID) (entity.Document, error) {
	panic("unimplemented")
}

// DeleteDocumentByID implements interfaces.DocumentRepository.
func (p *PostgresDocumentRepository) DeleteDocumentByID(documentID string) error {
	panic("unimplemented")
}

// GetDocumentByID implements interfaces.DocumentRepository.
func (p *PostgresDocumentRepository) GetDocumentByID(documentID string) (string, error) {
	panic("unimplemented")
}

// GetDocumentsByUserID implements interfaces.DocumentRepository.
func (p *PostgresDocumentRepository) GetDocumentsByUserID(userID uuid.UUID) ([]entity.Document, error) {
	panic("unimplemented")
}
