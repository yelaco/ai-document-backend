package entity

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID         uuid.UUID
	Title      string
	DocumentID uuid.UUID
	UserID     uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
