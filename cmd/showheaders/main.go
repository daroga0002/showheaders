package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"showheaders/internal/config"
	"showheaders/internal/server"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	// Parse configuration
	cfg := config.Parse()

	// Create and start server
	srv := server.New(cfg, logger)

	// Channel to listen for errors from server
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		serverErrors <- srv.Start()
	}()

	// Channel to listen for OS signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	case sig := <-shutdown:
		logger.Info("shutdown signal received", zap.String("signal", sig.String()))

		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("graceful shutdown failed", zap.Error(err))
			os.Exit(1)
		}

		logger.Info("server shutdown complete")
	}
}
