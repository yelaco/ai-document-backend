package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yelaco/ai-document-backend/internal/domain/interfaces"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/dtos"
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
	var req dtos.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(fmt.Errorf("AuthHandler.Register: failed to bind request body: %w", err))
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
		_ = c.Error(fmt.Errorf("AuthHandler.Register: failed to register user: %w", err))
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
	var req dtos.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(fmt.Errorf("AuthHandler.Login: failed to bind request body: %w", err))
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
		_ = c.Error(fmt.Errorf("AuthHandler.Login: failed to login user: %w", err))
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
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		_ = c.Error(fmt.Errorf("AuthHandler.RefreshAccessToken: missing refresh token cookie: %w", err))
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
		_ = c.Error(fmt.Errorf("AuthHandler.RefreshAccessToken: failed to refresh access token: %w", err))
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
