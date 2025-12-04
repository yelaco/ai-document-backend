package interfaces

import (
	"context"

	"github.com/yelaco/ai-document-backend/internal/domain/models/dtos"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
)

type AuthService interface {
	RegisterUser(ctx context.Context, params dtos.RegisterUserParams) (entity.User, error)
	LoginUser(ctx context.Context, params dtos.LoginUserParams) (entity.Auth, error)
	RefreshFlow(ctx context.Context, params dtos.RefreshFlowParams) (entity.Auth, error)
	GetPublicKey(ctx context.Context) string
}
