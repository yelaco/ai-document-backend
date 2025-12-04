package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/domain/models/dtos"
	"github.com/yelaco/ai-document-backend/internal/domain/models/entity"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
	reqContext "github.com/yelaco/ai-document-backend/internal/infrastructure/context"
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
func (a *AuthService) LoginUser(ctx context.Context, params dtos.LoginUserParams) (entity.Auth, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, params.Email)
	if err != nil {
		return entity.Auth{}, AuthErrorUserNotFound
	}
	if err := a.passwordHasher.VerifyPassword(user.PasswordHash, params.Password); err != nil {
		return entity.Auth{}, AuthErrorInvalidCredentials
	}

	accessToken, err := a.tokenMaker.CreateToken(user.ID.String(), params.Email, string(user.Role))
	if err != nil {
		return entity.Auth{}, fmt.Errorf("auth.AuthService.LoginUser: failed to create access token: %w", err)
	}

	refreshToken, err := createRefreshToken()
	if err != nil {
		return entity.Auth{}, fmt.Errorf("auth.AuthService.LoginUser: failed to create refresh token: %w", err)
	}
	refreshTokenHash, err := a.passwordHasher.HashPassword(refreshToken)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("auth.AuthService.LoginUser: failed to hash refresh token: %w", err)
	}
	newExpiresAt := time.Now().Add(RefreshTokenExpirationDays * 24 * time.Hour)
	err = a.refreshTokenRepo.StoreRefreshToken(ctx, user.ID, refreshTokenHash, newExpiresAt)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("auth.AuthService.LoginUser: failed to store refresh token: %w", err)
	}

	return entity.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RegisterUser implements interfaces.AuthService.
func (a *AuthService) RegisterUser(ctx context.Context, params dtos.RegisterUserParams) (entity.User, error) {
	passwordHash, err := a.passwordHasher.HashPassword(params.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("AuthService.RegisterUser: failed to hash password: %w", err)
	}

	newUser := entity.User{
		Email:        params.Email,
		FullName:     params.FullName,
		PasswordHash: passwordHash,
		Role:         auth.RoleUser,
	}

	err = a.userRepo.CreateUser(ctx, &newUser)
	if err != nil {
		return entity.User{}, fmt.Errorf("AuthService.RegisterUser: failed to create user: %w", err)
	}

	return newUser, nil
}

func (a *AuthService) RefreshFlow(ctx context.Context, params dtos.RefreshFlowParams) (entity.Auth, error) {
	userId := reqContext.UserIDMustFromContext(ctx)
	authClaims := reqContext.AuthClaimsMustFromContext(ctx)

	refreshTokenHash, err := a.refreshTokenRepo.GetRefreshTokenHash(ctx, userId)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("AuthService.RefreshFlow: failed to get refresh token hash: %w", err)
	}

	err = a.passwordHasher.VerifyPassword(refreshTokenHash, params.OldRefreshToken)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("AuthService.RefreshFlow: invalid refresh token: %w", err)
	}

	err = a.refreshTokenRepo.RevokeRefreshToken(ctx, userId)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("AuthService.RefreshFlow: failed to revoke refresh token: %w", err)
	}

	accessToken, err := a.tokenMaker.CreateToken(userId.String(), authClaims.Email, string(authClaims.Role))
	if err != nil {
		return entity.Auth{}, fmt.Errorf("AuthService.RefreshFlow: failed to create access token: %w", err)
	}

	newRefreshToken, err := createRefreshToken()
	if err != nil {
		return entity.Auth{}, fmt.Errorf("AuthService.RefreshFlow: failed to create refresh token: %w", err)
	}
	newRefreshTokenHash, err := a.passwordHasher.HashPassword(newRefreshToken)
	if err != nil {
		return entity.Auth{}, fmt.Errorf("AuthService.RefreshFlow: failed to hash refresh token: %w", err)
	}
	newExpiresAt := time.Now().Add(RefreshTokenExpirationDays * 24 * time.Hour)
	a.refreshTokenRepo.StoreRefreshToken(ctx, userId, newRefreshTokenHash, newExpiresAt)

	return entity.Auth{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *AuthService) GetPublicKey(ctx context.Context) string {
	return a.tokenMaker.GetPublicKey()
}
