package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
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

func (p *PostgresRefreshTokenRepository) DeleteRefreshToken(userID uuid.UUID) error {
	panic("unimplemented")
}

func (p *PostgresRefreshTokenRepository) GetRefreshTokenHash(userID uuid.UUID) (string, error) {
	panic("unimplemented")
}

func (p *PostgresRefreshTokenRepository) RevokeRefreshToken(userID uuid.UUID) error {
	panic("unimplemented")
}

func (p *PostgresRefreshTokenRepository) StoreRefreshToken(refreshToken string, userID uuid.UUID, expiredsAt time.Duration) error {
	panic("unimplemented")
}
