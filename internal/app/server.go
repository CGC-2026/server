package app

import (
	"context"
	"net/http"
	"time"

	"server/internal/http/middleware"
	"server/internal/log"
)

// Server represents the HTTP server
type Server struct {
	server *http.Server
	logger log.Logger
}

// ServerConfig contains configuration for the server
type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

// DefaultServerConfig returns the default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port:           "8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}

// NewServer creates a new server with the given config
func NewServer(config ServerConfig) *Server {
	logger := log.NewStandardLogger()

	mux := http.NewServeMux()

	RegisterRoutes(mux, logger)

	// Apply JSON middleware
	handler := middleware.JSONMiddleware(mux)

	srv := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	return &Server{
		server: srv,
		logger: logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting server on port %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")
	return s.server.Shutdown(ctx)
}
