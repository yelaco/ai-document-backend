package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	http.Server
}

func NewServer(logger *zap.Logger, address string, router Router) *Server {
	return &Server{
		logger: logger,
		Server: http.Server{
			Addr:    address,
			Handler: router.GetHandler(),
		},
	}
}

func (s *Server) Start() {
	go func() {
		// service connections
		s.logger.Info("Server starting", zap.String("addr", s.Addr))
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("Listen error", zap.Error(err))
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
	s.logger.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := s.Shutdown(ctx); err != nil {
		s.logger.Fatal("Server shutdown failed", zap.Error(err))
	}
	s.logger.Info("Server exited gracefully")
}
