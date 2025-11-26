package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/dtos"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/middlewares"
	"go.uber.org/zap"
)

const RefreshTokenCookieExpiryInDays = 7

type AuthHandler struct {
	authService interfaces.AuthService
}

func NewAuthHandler(authService interfaces.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	l := middlewares.MustGetContextualLogger(c)
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		l.Error("AuthHandler.Login: failed to extract request body", zap.Error(err))
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
		l.Error("AuthHandler.Register: failed to register user", zap.Error(err))
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
	l := middlewares.MustGetContextualLogger(c)
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		l.Error("AuthHandler.Login: failed to extract request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusBadRequest,
				ErrorMessage: "invalid request payload",
			},
		})
	}
	auth, err := h.authService.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		l.Error("AuthHandler.Login: failed to login user", zap.Error(err))
		c.JSON(http.StatusUnauthorized, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "invalid email or password",
			},
		})
		return
	}

	// Set refresh token as HttpOnly cookie
	setRefreshTokenCookie(c, auth.RefreshToken, false)

	c.JSON(http.StatusOK, dtos.BaseAPIResponse{
		Status: dtos.StatusSuccess,
		Data: dtos.TokenResponse{
			AccessToken: auth.AccessToken,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	setRefreshTokenCookie(c, "", true)
	c.JSON(http.StatusOK, dtos.BaseAPIResponse{
		Status: dtos.StatusSuccess,
	})
}

func (h *AuthHandler) GetPublicKey(c *gin.Context) {
	publicKey := h.authService.GetPublicKey(c.Request.Context())

	c.JSON(http.StatusOK, dtos.BaseAPIResponse{
		Status: dtos.StatusSuccess,
		Data: gin.H{
			"publicKey": publicKey,
		},
	})
}

func (h *AuthHandler) RefreshAccessToken(c *gin.Context) {
	l := middlewares.MustGetContextualLogger(c)
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		l.Error("AuthHandler.RefreshAccessToken: failed to get refresh token from cookie", zap.Error(err))
		c.JSON(http.StatusUnauthorized, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "missing refresh token",
			},
		})
		return
	}

	auth, err := h.authService.RefreshFlow(c.Request.Context(), refreshToken)
	if err != nil {
		l.Error("AuthHandler.RefreshAccessToken: failed to refresh access token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, dtos.BaseErrorResponse{
			Status: dtos.StatusError,
			Error: dtos.ErrorResponse{
				ErrorCode:    http.StatusUnauthorized,
				ErrorMessage: "invalid refresh token",
			},
		})
		return
	}

	setRefreshTokenCookie(c, auth.RefreshToken, false)
	c.JSON(http.StatusOK, dtos.BaseAPIResponse{
		Status: dtos.StatusSuccess,
		Data: dtos.TokenResponse{
			AccessToken: auth.AccessToken,
		},
	})
}

func setRefreshTokenCookie(c *gin.Context, refreshToken string, expired bool) {
	if expired {
		c.SetCookie(
			"refresh_token",
			"",
			0,
			"/api/auth/refresh",
			"localhost",
			false,
			true,
		)
	} else {
		c.SetCookie(
			"refresh_token",
			refreshToken,
			RefreshTokenCookieExpiryInDays*24*60*60,
			"/api/auth/refresh",
			"localhost",
			false,
			true,
		)
	}
	c.SetSameSite(http.SameSiteStrictMode)
}
