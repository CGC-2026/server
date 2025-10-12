package http

import (
	"context"
	"net/http"
	"time"

	"server/internal/data"
	"server/internal/http/middleware"
	"server/internal/log"
)

// Server represents the HTTP server
type Server struct {
	cfg    ServerConfig
	server *http.Server
	store  *data.Store
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
		Port:           "3000",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}

// NewServer creates a new server with the given config
func NewServer(config ServerConfig) *Server {
	logger := log.NewStandardLogger()

	// Initialize Clerk
	if err := middleware.InitClerk(); err != nil {
		logger.Info("Clerk not initialized %v", err)
	} else {
		logger.Info("Clerk initialized successfully")
	}

	// Initialize db store
	store, err := data.NewStore(context.Background())
	if err != nil {
		logger.Error("failed to init data store: %v", err)
		panic(err)
	}

	// Create router and register routes
	r := New(logger, store)
	r.RegisterRoutes()

	// Apply middleware to router
	handler := r.Handler()
	handler = middleware.JSONMiddleware(handler)

	srv := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	return &Server{
		cfg:    config,
		server: srv,
		store:  store,
		logger: logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting server at http://localhost%s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")
	_ = s.server.Shutdown(ctx)

	if s.store != nil {
		s.store.Close()
	}
	return nil
}
