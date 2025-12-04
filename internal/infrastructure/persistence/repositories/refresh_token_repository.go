package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/persistence/database/sqlc"
)

type PostgresRefreshTokenRepository struct {
	connPool *pgxpool.Pool
	queries  *sqlc.Queries
}

func NewRefreshTokenRepository(connPool *pgxpool.Pool) interfaces.RefreshTokenRepository {
	return &PostgresRefreshTokenRepository{
		connPool: connPool,
		queries:  sqlc.New(connPool),
	}
}

func (p *PostgresRefreshTokenRepository) DeleteRefreshToken(ctx context.Context, userID types.UserID) error {
	err := p.queries.DeleteRefreshTokenByUserID(ctx, userID.UUID())
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}

func (p *PostgresRefreshTokenRepository) GetRefreshTokenHash(ctx context.Context, userID types.UserID) (string, error) {
	row, err := p.queries.GetRefreshTokenHashByUserID(ctx, userID.UUID())
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token hash: %w", err)
	}
	return row.TokenHash, nil
}

func (p *PostgresRefreshTokenRepository) RevokeRefreshToken(ctx context.Context, userID types.UserID) error {
	err := p.queries.RevokeRefreshTokenByUserID(ctx, userID.UUID())
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

func (p *PostgresRefreshTokenRepository) StoreRefreshToken(ctx context.Context, userID types.UserID, refreshTokenHash string, expiredsAt time.Time) error {
	err := p.queries.InsertRefreshToken(ctx, sqlc.InsertRefreshTokenParams{
		UserID:    userID.UUID(),
		TokenHash: refreshTokenHash,
		ExpiresAt: expiredsAt,
	})
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	return nil
}
