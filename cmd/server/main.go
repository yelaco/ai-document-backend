package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yelaco/ai-document-backend/internal/application/auth"
	"github.com/yelaco/ai-document-backend/internal/application/document"
	"github.com/yelaco/ai-document-backend/internal/config"
	authInfra "github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/persistence/repositories"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/token"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/handlers"
	"github.com/yelaco/ai-document-backend/pkg/logger"
	"github.com/yelaco/ai-document-backend/pkg/server"
	"go.uber.org/zap"
)

func main() {
	// load config
	cfg := config.MustLoadConfig("./configs")

	// setup logger
	logger := logger.NewLogger(cfg.App.Env)
	defer logger.Sync()

	ctx := context.Background()

	// setup database connection
	connPool, err := pgxpool.New(ctx, fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	))
	if err != nil {
		logger.Fatal("failed to create connection pool:", zap.Error(err))
	}
	if err := connPool.Ping(ctx); err != nil {
		logger.Fatal("failed to ping database:", zap.Error(err))
	}

	// setup paseto token maker
	tokenMaker, err := token.NewPasetoV4Maker()
	if err != nil {
		logger.Fatal("failed to create token maker:", zap.Error(err))
	}

	// inject dependencies
	passwordHasher := authInfra.NewArgon2PasswordHasher()
	userRepo := repositories.NewPostgresUserRepository(connPool)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(connPool)
	documentRepo := repositories.NewDocumentRepository(connPool)
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, passwordHasher, tokenMaker)
	documentService := document.NewDocumentService(documentRepo)
	userHandler := handlers.NewUserHandler()
	authHandler := handlers.NewAuthHandler(authService)
	documentHandler := handlers.NewDocumentHandler(documentService)

	// setup router
	router := rest.NewRouter(logger)
	router.SetupRoutes(tokenMaker, userHandler, authHandler, documentHandler)

	// start server
	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	server := server.NewServer(logger, addr, router)
	server.Start()
}
