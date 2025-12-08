package server

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"showheaders/internal/config"
	"showheaders/internal/handlers"
	"showheaders/internal/middleware"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

// New creates a new server instance
func New(cfg *config.Config, logger *zap.Logger) *Server {
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/", handlers.ShowHeaders)
	mux.HandleFunc("/health", handlers.Health)

	// Paths to exclude from logging
	excludePaths := map[string]bool{
		"/health": true,
	}

	// Apply middleware chain
	handler := middleware.NoCacheMiddleware(
		middleware.LoggingMiddleware(logger, excludePaths)(mux),
	)

	httpServer := &http.Server{
		Addr:         cfg.Address(),
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("starting server",
		zap.String("address", s.httpServer.Addr),
	)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")
	return s.httpServer.Shutdown(ctx)
}
