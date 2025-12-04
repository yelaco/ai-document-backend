package interfaces

import (
	"context"
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type RefreshTokenRepository interface {
	StoreRefreshToken(ctx context.Context, userID types.UserID, refreshToken string, expiredsAt time.Time) error
	GetRefreshTokenHash(ctx context.Context, userID types.UserID) (string, error)
	RevokeRefreshToken(ctx context.Context, userID types.UserID) error
	DeleteRefreshToken(ctx context.Context, userID types.UserID) error
}
