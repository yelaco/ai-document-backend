package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	logger *zap.Logger
}

func NewUserHandler(logger *zap.Logger) *UserHandler {
	return &UserHandler{
		logger: logger.With(zap.String("handler", "UserHandler")),
	}
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	h.logger.Info("GetUserByID called")
	// Implementation goes here
}
