package entity

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID           uuid.UUID
	OriginalName string
	SavePath     string
	Status       string
	UserID       uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
