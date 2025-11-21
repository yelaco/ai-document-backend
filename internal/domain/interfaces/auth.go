package interfaces

import (
	"context"

	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type AuthService interface {
	RegisterUser(ctx context.Context, email string, fullName string, password string) (entity.User, error)
	LoginUser(ctx context.Context, email string, password string) (entity.Auth, error)
}
