package context

import (
	"context"

	"github.com/google/uuid"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
)

type UserIDKey struct{}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, UserIDKey{}, id)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey{}).(string)
	if !ok {
		return uuid.Nil, false
	}
	userId, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, false
	}
	return userId, ok
}

func UserIDMustFromContext(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(UserIDKey{}).(string)
	if !ok {
		panic("user ID not found in context")
	}
	userId, err := uuid.Parse(id)
	if err != nil {
		panic("invalid user ID format in context")
	}
	return userId
}

func UserIDTryFromContext(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(UserIDKey{}).(string)
	if !ok {
		return uuid.Nil
	}
	userId, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil
	}
	return userId
}

type AuthClaimsKey struct{}

func WithAuthClaims(ctx context.Context, claims *auth.AuthClaims) context.Context {
	return context.WithValue(ctx, AuthClaimsKey{}, claims)
}

func AuthClaimsFromContext(ctx context.Context) (*auth.AuthClaims, bool) {
	claims, ok := ctx.Value(AuthClaimsKey{}).(*auth.AuthClaims)
	return claims, ok
}

func AuthClaimsMustFromContext(ctx context.Context) *auth.AuthClaims {
	claims, ok := ctx.Value(AuthClaimsKey{}).(*auth.AuthClaims)
	if !ok {
		panic("auth claims not found in context")
	}
	return claims
}
