package entity

import (
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type Message struct {
	ID        types.MessageID
	ChatID    types.ChatID
	UserID    types.UserID
	Content   string
	Role      string
	Timestamp time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
