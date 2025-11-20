package interfaces

import (
	"time"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	StoreRefreshToken(refreshToken string, userID uuid.UUID, expiredsAt time.Duration) error
	GetRefreshTokenHash(userID uuid.UUID) (string, error)
	RevokeRefreshToken(userID uuid.UUID) error
	DeleteRefreshToken(userID uuid.UUID) error
}
