package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
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
		return fmt.Errorf("PostgresUserRepository.CreateUser: failed to create user: %w", err)
	}
	user.ID = types.NewUserID(row.ID)
	user.CreatedAt = row.CreatedAt
	user.UpdatedAt = row.UpdatedAt
	return nil
}

// GetUserByEmail implements interfaces.UserRepository.
func (p *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	row, err := p.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.CreateUser: failed to get user by email: %w", err)
	}
	return &entity.User{
		ID:           types.NewUserID(row.ID),
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		FullName:     row.FullName,
		Role:         auth.Role(row.Role),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}

// GetUserByID implements interfaces.UserRepository.
func (p *PostgresUserRepository) GetUserByID(ctx context.Context, id types.UserID) (*entity.User, error) {
	row, err := p.queries.GetUserByID(ctx, id.UUID())
	if err != nil {
		return nil, fmt.Errorf("PostgresUserRepository.CreateUser: failed to get user by ID: %w", err)
	}
	return &entity.User{
		ID:           types.NewUserID(row.ID),
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		FullName:     row.FullName,
		Role:         auth.Role(row.Role),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}
