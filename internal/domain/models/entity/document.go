package entity

import (
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/models/types"
)

type Document struct {
	ID           types.DocumentID
	OriginalName string
	SavePath     string
	Status       string
	UserID       types.UserID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
