package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/dtos"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/handlers"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/middleware"
	"go.uber.org/zap"
)

type Router struct {
	logger *zap.Logger
	engine *gin.Engine
}

func NewRouter(logger *zap.Logger) *Router {
	r := Router{
		logger: logger,
		engine: gin.New(),
	}
	r.engine.Use(middleware.ZapLogger(logger))
	r.engine.Use(gin.Recovery())
	return &r
}

func (r *Router) SetupRoutes(userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler) {
	// health check
	r.engine.GET("/health-check", func(c *gin.Context) {
		data := dtos.BaseAPIResponse{
			Status: "success",
			Data:   gin.H{"message": "Health check OK!"},
		}
		c.JSON(http.StatusOK, data)
	})

	apiRouter := r.engine.Group("/api")
	{
		userRouter := apiRouter.Group("/users")
		{
			userRouter.GET("/:id", userHandler.GetUserByID)
		}
		authRouter := apiRouter.Group("/auth")
		{
			authRouter.POST("/register", authHandler.Register)
			authRouter.POST("/login", authHandler.Login)
			authRouter.POST("/logout", authHandler.Logout)
			authRouter.POST("/refresh", authHandler.RefreshAccessToken)
		}
	}
}

func (r *Router) GetHandler() http.Handler {
	return r.engine
}
