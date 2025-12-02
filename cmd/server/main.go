package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	chromaLogger "github.com/amikos-tech/chroma-go/pkg/logger"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yelaco/ai-document-backend/internal/application/auth"
	"github.com/yelaco/ai-document-backend/internal/application/document"
	"github.com/yelaco/ai-document-backend/internal/config"
	authInfra "github.com/yelaco/ai-document-backend/internal/infrastructure/auth"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/persistence/repositories"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/rag"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/tasks"
	"github.com/yelaco/ai-document-backend/internal/infrastructure/token"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest/handlers"
	"github.com/yelaco/ai-document-backend/pkg/logger"
	"github.com/yelaco/ai-document-backend/pkg/server"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGINT,
	syscall.SIGTERM,
}

func main() {
	ctx := context.Background()

	// load config
	cfg := config.MustLoadConfig("./configs")

	// setup logger
	logger := logger.NewLogger(cfg.App.Env)
	defer func() {
		err := logger.Sync()
		if err != nil && !errors.Is(err, syscall.ENOTTY) {
			fmt.Println("failed to sync logger: " + err.Error())
		}
	}()

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

	// setup chroma client
	chromaClient, err := chroma.NewHTTPClient(
		chroma.WithBaseURL("http://localhost:8000"),
		chroma.WithLogger(chromaLogger.NewZapLogger(logger)),
	)
	if err != nil {
		logger.Fatal("failed to create chroma client:", zap.Error(err))
	}
	if err := chromaClient.Heartbeat(ctx); err != nil {
		logger.Fatal("failed to heartbeat chroma client:", zap.Error(err))
	}
	defer func() {
		if err := chromaClient.Close(); err != nil {
			logger.Error("failed to close chroma client:", zap.Error(err))
		}
	}()

	// setup dependencies
	passwordHasher := authInfra.NewArgon2PasswordHasher()
	userRepo := repositories.NewPostgresUserRepository(connPool)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(connPool)
	documentRepo := repositories.NewDocumentRepository(connPool)
	ragEmbedder := rag.NewChromaEmbedder(cfg.AI.GeminiAPIKey)
	ragStore := rag.NewChromaStore(chromaClient)

	// setup background task processor
	redisOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	}
	taskDistributor := tasks.NewAsynqTaskDistributor(redisOpt, logger)
	taskProcessor := tasks.NewAsynqProcessor(redisOpt, logger, documentRepo, ragEmbedder, ragStore)

	// inject dependencies
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, passwordHasher, tokenMaker)
	documentService := document.NewDocumentService(documentRepo)
	userHandler := handlers.NewUserHandler()
	authHandler := handlers.NewAuthHandler(authService)
	documentHandler := handlers.NewDocumentHandler(documentService, taskDistributor)

	// setup router and server
	router := rest.NewRouter(logger)
	router.SetupRoutes(tokenMaker, userHandler, authHandler, documentHandler)
	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	server := server.NewServer(logger, addr, router)

	// run the servers concurrently
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()
	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, taskProcessor)
	runHttpServer(ctx, waitGroup, server)

	err = waitGroup.Wait()
	if err != nil {
		logger.Fatal("Application exited with error:", zap.Error(err))
	}
}

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, taskProcessor tasks.TaskProcessor) {
	if err := taskProcessor.Start(); err != nil {
		panic(fmt.Errorf("failed to start task processor: %w", err))
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		taskProcessor.Shutdown()
		return nil
	})
}

func runHttpServer(ctx context.Context, waitGroup *errgroup.Group, server *server.Server) {
	waitGroup.Go(func() error {
		server.Start(ctx)
		return nil
	})
}
