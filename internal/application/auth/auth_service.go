package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/token"
	"github.com/yelaco/ai-document-backend/pkg/util"
)

type AuthService struct {
	userRepo         interfaces.UserRepository
	refreshTokenRepo interfaces.RefreshTokenRepository
	passwordHasher   util.PasswordHasher
	tokenMaker       token.Maker
}

func NewAuthService(
	userRepo interfaces.UserRepository,
	refreshTokenRepo interfaces.RefreshTokenRepository,
	passwordHasher util.PasswordHasher,
	tokenMaker token.Maker,
) interfaces.AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		passwordHasher:   passwordHasher,
		tokenMaker:       tokenMaker,
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

	accessToken, err := a.tokenMaker.CreateToken(user.ID.String(), email, string(user.Role))
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := createRefreshToken()
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to create refresh token: %w", err)
	}
	refreshTokenHash, err := a.passwordHasher.HashPassword(refreshToken)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to hash refresh token: %w", err)
	}
	newExpiresAt := time.Now().Add(RefreshTokenExpirationDays * 24 * time.Hour)
	err = a.refreshTokenRepo.StoreRefreshToken(ctx, user.ID, refreshTokenHash, newExpiresAt)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return entity.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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

func (a *AuthService) RefreshFlow(ctx context.Context, oldRefreshToken string) (entity.Auth, error) {
	userId, exist := UserIDFromContext(ctx)
	if !exist {
		return entity.Auth{}, AuthErrorInvalidCredentials
	}

	tokenHash, err := a.refreshTokenRepo.GetRefreshTokenHash(ctx, userId)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to get refresh token hash: %w", err)
	}

	err = a.passwordHasher.VerifyPassword(oldRefreshToken, tokenHash)
	if err != nil {
		return entity.Auth{}, AuthErrorInvalidCredentials
	}

	err = a.refreshTokenRepo.RevokeRefreshToken(ctx, userId)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	authClaims, exist := AuthClaimsFromContext(ctx)
	if !exist {
		return entity.Auth{}, AuthErrorInvalidCredentials
	}
	accessToken, err := a.tokenMaker.CreateToken(userId.String(), authClaims.Email, string(authClaims.Role))
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := createRefreshToken()
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to create refresh token: %w", err)
	}
	refreshTokenHash, err := a.passwordHasher.HashPassword(refreshToken)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("failed to hash refresh token: %w", err)
	}
	newExpiresAt := time.Now().Add(RefreshTokenExpirationDays * 24 * time.Hour)
	a.refreshTokenRepo.StoreRefreshToken(ctx, userId, refreshTokenHash, newExpiresAt)

	return entity.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *AuthService) GetPublicKey(ctx context.Context) string {
	return a.tokenMaker.GetPublicKey()
}
