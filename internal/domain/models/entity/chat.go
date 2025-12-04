package entity

import (
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type Chat struct {
	ID         types.ChatID
	Title      string
	DocumentID types.DocumentID
	UserID     types.UserID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
