package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/token"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/dtos"
	"go.uber.org/zap"
)

func AuthMiddleware(logger *zap.Logger, tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		var err error
		if strings.HasPrefix(c.Request.URL.Path, "api/auth/refresh") {
			_, err = tokenMaker.VerifyToken(token, true)
		} else {
			_, err = tokenMaker.VerifyToken(token, false)
		}
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
		c.Next()
	}
}
