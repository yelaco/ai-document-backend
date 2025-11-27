package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	authInfra "github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
	reqContext "github.com/yelaco/ai-document-backend/internal/infrastructure/context"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/token"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/dtos"
	"go.uber.org/zap"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		l, ok := c.Get("logger")
		if !ok {
			panic("logger not found in context")
		}
		logger, _ := l.(*zap.Logger)

		authHeader := c.Request.Header.Get("Authorization")
		accessToken, err := getBearerToken(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(401, dtos.BaseErrorResponse{
				Status: dtos.StatusError,
				Error: dtos.ErrorResponse{
					ErrorCode:    401,
					ErrorMessage: "unauthorized: invalid authorization header format",
				},
			})
			return
		}

		skipExpirationCheck := strings.HasPrefix(c.Request.URL.Path, "api/auth/refresh")
		claims, err := tokenMaker.VerifyToken(accessToken, skipExpirationCheck)
		if err != nil {
			logger.Error("Unauthorized access", zap.Error(err))
			c.AbortWithStatusJSON(401, dtos.BaseErrorResponse{
				Status: dtos.StatusError,
				Error: dtos.ErrorResponse{
					ErrorCode:    401,
					ErrorMessage: "unauthorized: invalid or expired token",
				},
			})
			return
		}

		role, err := authInfra.ParseRole(claims.Role)
		if err != nil {
			c.AbortWithStatusJSON(401, dtos.BaseErrorResponse{
				Status: dtos.StatusError,
				Error: dtos.ErrorResponse{
					ErrorCode:    401,
					ErrorMessage: "unauthorized: invalid role",
				},
			})
			return
		}

		ctx := c.Request.Context()
		ctx = reqContext.WithUserID(ctx, claims.Sub)
		ctx = reqContext.WithAuthClaims(ctx, &authInfra.AuthClaims{
			Email: claims.Email,
			Role:  role,
		})
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func getBearerToken(authorizationHeader string) (string, error) {
	if authorizationHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	parts := strings.SplitN(authorizationHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}
