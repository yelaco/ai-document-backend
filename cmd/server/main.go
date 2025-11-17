package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yelaco/ai-document-backend/internal/config"
	"github.com/yelaco/ai-document-backend/internal/presentation/rest"
	"github.com/yelaco/ai-document-backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// load config
	cfg := config.MustLoadConfig("./configs")

	// setup logger
	logger := logger.NewLogger(cfg.AppEnv)
	defer logger.Sync()

	// setup router
	engine := gin.New()
	router := rest.NewRouter(logger, engine)
	router.SetupRoutes()

	// start server
	srv := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: router.GetHandler(),
	}

	go func() {
		// service connections
		logger.Info("Server starting", zap.String("addr", cfg.Host+":"+cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Listen error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 30 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown failed", zap.Error(err))
	}
	logger.Info("Server exited gracefully")
}
