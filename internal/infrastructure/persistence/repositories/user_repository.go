package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/persistence/database/sqlc"
)

type PostgresUserRepository struct {
	connPool *pgxpool.Pool
	queries  *sqlc.Queries
}

func NewPostgresUserRepository(connPool *pgxpool.Pool) interfaces.UserRepository {
	return &PostgresUserRepository{
		connPool: connPool,
		queries:  sqlc.New(connPool),
	}
}

// CreateUser implements interfaces.UserRepository.
func (p *PostgresUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	panic("unimplemented")
}

// GetUserByEmail implements interfaces.UserRepository.
func (p *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	panic("unimplemented")
}

// GetUserByID implements interfaces.UserRepository.
func (p *PostgresUserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	panic("unimplemented")
}
