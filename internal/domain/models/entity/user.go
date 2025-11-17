package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
)

type User struct {
	ID           uuid.UUID
	Email        string
	FullName     string
	PasswordHash string
	Role         auth.Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
