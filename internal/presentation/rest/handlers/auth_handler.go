package handlers

import (
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/dtos"
	"go.uber.org/zap"
)

type AuthHandler struct {
	logger      *zap.Logger
	authService interfaces.AuthService
}

func NewAuthHandler(logger *zap.Logger, authService interfaces.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger.With(zap.String("presentation", "AuthHandler")),
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	h.logger.Info(
		"Register called",
		zap.String("request_id", requestid.Get(c)),
	)

	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "invalid request payload",
			},
		})
	}

	user, err := h.authService.RegisterUser(c.Request.Context(), req.Email, req.FullName, req.Password)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusInternalServerError,
				ErrorMessage: "failed to register user",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, dtos.BaseAPIResponse{
		Status: dtos.StatusSuccess,
		Data: gin.H{
			"user": dtos.UserResponseFromEntity(&user),
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	h.logger.Info(
		"Login called",
		zap.String("request_id", requestid.Get(c)),
	)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.logger.Info(
		"Logout called",
		zap.String("request_id", requestid.Get(c)),
	)
}

func (h *AuthHandler) RefreshAccessToken(c *gin.Context) {
	h.logger.Info(
		"RefreshAccessToken called",
		zap.String("request_id", requestid.Get(c)),
	)
}
