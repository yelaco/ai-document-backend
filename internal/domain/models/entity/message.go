package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	UserID    uuid.UUID
	Content   string
	Role      string
	Timestamp time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
