package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger      *zap.Logger
	authService interfaces.AuthService
}

func NewAuthHandler(logger *zap.Logger, authService interfaces.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	h.logger.Info("Register called")
	// Implementation goes here
}

func (h *AuthHandler) Login(c *gin.Context) {
	h.logger.Info("Login called")
	// Implementation goes here
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.logger.Info("Logout called")
	// Implementation goes here
}

func (h *AuthHandler) RefreshAccessToken(c *gin.Context) {
	h.logger.Info("RefreshAccessToken called")
	// Implementation goes here
}
