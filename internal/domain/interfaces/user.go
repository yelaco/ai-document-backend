package interfaces

import (
	"context"

	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type UserService interface {
	GetUser(ctx context.Context, id int) string
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
}
