package auth

import (
	"context"

	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/pkg/util"
)

type AuthService struct {
	userRepo         interfaces.UserRepository
	refreshTokenRepo interfaces.RefreshTokenRepository
	passwordHasher   util.PasswordHasher
}

func NewAuthService(userRepo interfaces.UserRepository, refreshTokenRepo interfaces.RefreshTokenRepository, passwordHasher util.PasswordHasher) interfaces.AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		passwordHasher:   passwordHasher,
	}
}

// LoginUser implements interfaces.AuthService.
func (a *AuthService) LoginUser(ctx context.Context, email string, password string) (string, string, error) {
	panic("unimplemented")
}

// RegisterUser implements interfaces.AuthService.
func (a *AuthService) RegisterUser(ctx context.Context, email string, fullName string, password string) error {
	panic("unimplemented")
}
