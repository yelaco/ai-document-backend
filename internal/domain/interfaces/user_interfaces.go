package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type UserService interface {
	GetUser(id int) string
}

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (entity.User, error)
}
