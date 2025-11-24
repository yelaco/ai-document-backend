package repositories

import (
	"context"
	"fmt"

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
	row, err := p.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		Role:         string(user.Role),
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	user.ID = row.ID
	user.CreatedAt = row.CreatedAt
	user.UpdatedAt = row.UpdatedAt
	return nil
}

// GetUserByEmail implements interfaces.UserRepository.
func (p *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	panic("unimplemented")
}

// GetUserByID implements interfaces.UserRepository.
func (p *PostgresUserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	panic("unimplemented")
}
