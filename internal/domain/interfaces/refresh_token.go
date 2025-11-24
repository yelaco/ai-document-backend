package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string, expiredsAt time.Time) error
	GetRefreshTokenHash(ctx context.Context, userID uuid.UUID) (string, error)
	RevokeRefreshToken(ctx context.Context, userID uuid.UUID) error
	DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error
}
