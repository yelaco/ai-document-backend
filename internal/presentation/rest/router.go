package rest

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/token"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/dtos"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/handlers"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/middlewares"
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
	r.engine.Use(gin.Recovery())
	r.engine.Use(cors.New(cors.Config{
		AllowOrigins:           []string{"*"},
		AllowCredentials:       true,
		AllowBrowserExtensions: false,
		AllowFiles:             true,
	}))
	r.engine.Use(requestid.New())
	r.engine.MaxMultipartMemory = 32 << 20 // 32 MiB
	return &r
}

func (r *Router) SetupRoutes(
	tokenMaker token.Maker,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	documentHandler *handlers.DocumentHandler,
	sseHandler *handlers.SSEHandler,
) {
	r.engine.Use(middlewares.ZapLoggerMiddleware(r.logger))
	authMiddlware := middlewares.AuthMiddleware(tokenMaker)

	// health check
	r.engine.GET("/health-check", func(c *gin.Context) {
		data := dtos.BaseAPIResponse{
			Status: dtos.StatusSuccess,
			Data:   gin.H{"message": "Health check OK!"},
		}
		c.JSON(http.StatusOK, data)
	})

	apiRouter := r.engine.Group("/api")
	{
		userRouter := apiRouter.Group("/users")
		{
			userRouter.Use(authMiddlware)
			userRouter.GET("/:id", userHandler.GetUserByID)
		}
		authRouter := apiRouter.Group("/auth")
		{
			authRouter.POST("/register", authHandler.Register)
			authRouter.POST("/login", authHandler.Login)
			authRouter.POST("/logout", authMiddlware, authHandler.Logout)
			authRouter.POST("/refresh", authMiddlware, authHandler.RefreshAccessToken)
			authRouter.GET("/publicKey", authHandler.GetPublicKey)
		}
		documentRouter := apiRouter.Group("/documents")
		{
			documentRouter.Use(authMiddlware)
			documentRouter.POST("/", documentHandler.UploadDocument)
			documentRouter.GET("/", documentHandler.GetPaginatedDocuments)
			documentRouter.GET("/:id", documentHandler.GetDocumentByID)
			documentRouter.DELETE("/:id", documentHandler.DeleteDocument)
		}
		chatRouter := apiRouter.Group("/chat")
		{
			chatRouter.Use(authMiddlware)
		}
		sseRouter := apiRouter.Group("/stream")
		{
			sseRouter.Use(authMiddlware)
			sseRouter.GET("/chat/:id", sseHandler.StreamChat)
		}
	}
}

func (r *Router) GetHandler() http.Handler {
	return r.engine
}
