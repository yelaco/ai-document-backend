package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yelaco/ai-document-backend/internal/application/auth"
	"github.com/yelaco/ai-document-backend/internal/config"
	authInfra "github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/persistence/repositories"
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

	// inject dependencies
	userRepo := repositories.NewPostgresUserRepository(connPool)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(connPool)
	passwordHasher := authInfra.NewArgon2PasswordHasher()
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, passwordHasher)
	userHandler := handlers.NewUserHandler(logger)
	authHandler := handlers.NewAuthHandler(logger, authService)

	// setup router
	router := rest.NewRouter(logger)
	router.SetupRoutes(userHandler, authHandler)

	// start server
	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	server := server.NewServer(logger, addr, router)
	server.Start()
}
