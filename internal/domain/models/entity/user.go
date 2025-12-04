package entity

import (
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
)

type User struct {
	ID           types.UserID
	Email        string
	FullName     string
	PasswordHash string
	Role         auth.Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
