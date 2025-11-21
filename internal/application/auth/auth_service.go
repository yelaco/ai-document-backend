package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
	"github.com/yelaco/ai-document-backend/pkg/util"
	"go.uber.org/zap"
)

type AuthService struct {
	logger           *zap.Logger
	userRepo         interfaces.UserRepository
	refreshTokenRepo interfaces.RefreshTokenRepository
	passwordHasher   util.PasswordHasher
}

func NewAuthService(logger *zap.Logger, userRepo interfaces.UserRepository, refreshTokenRepo interfaces.RefreshTokenRepository, passwordHasher util.PasswordHasher) interfaces.AuthService {
	return &AuthService{
		logger:           logger.With(zap.String("application", "AuthService")),
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		passwordHasher:   passwordHasher,
	}
}

// LoginUser implements interfaces.AuthService.
func (a *AuthService) LoginUser(ctx context.Context, email string, password string) (entity.Auth, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return entity.Auth{}, AuthErrorUserNotFound
	}
	if err := a.passwordHasher.VerifyPassword(user.PasswordHash, password); err != nil {
		return entity.Auth{}, AuthErrorInvalidCredentials
	}

	accessToken, err := createAccessToken(user.ID.String(), email, string(user.Role))
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := createRefreshToken()
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	newExpiresAt := time.Now().Add(RefreshTokenExpirationDays * 24 * time.Hour)
	err = a.refreshTokenRepo.StoreRefreshToken(refreshToken, user.ID, newExpiresAt)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return entity.Auth{
		AccessToken: accessToken, RefreshToken: refreshToken,
	}, nil
}

// RegisterUser implements interfaces.AuthService.
func (a *AuthService) RegisterUser(ctx context.Context, email string, fullName string, password string) (entity.User, error) {
	passwordHash, err := a.passwordHasher.HashPassword(password)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := entity.User{
		Email:        email,
		FullName:     fullName,
		PasswordHash: passwordHash,
		Role:         auth.RoleUser,
	}

	err = a.userRepo.CreateUser(ctx, &newUser)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}
