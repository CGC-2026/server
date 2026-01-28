package http

import (
	"net/http"
	"server/internal/data"
	"server/internal/http/handlers"
	"server/internal/http/middleware"
	"server/internal/log"
)

// Router manages all HTTP routes for the application
type Router struct {
	mux    *http.ServeMux
	logger log.Logger
	store  *data.Store
}

// New creates a new router instance
func New(logger log.Logger, store *data.Store) *Router {
	return &Router{
		mux:    http.NewServeMux(),
		logger: logger,
		store:  store,
	}
}

// Handler returns the http.Handler for the router
func (r *Router) Handler() http.Handler {
	return r.mux
}

// RegisterRoutes registers all application routes
func (r *Router) RegisterRoutes() {
	// Public routes
	handlers.RegisterHealthRoutes(r.mux, r.logger)
	handlers.RegisterAuthRoutes(r.mux, r.logger, r.store) // Clerk webhooks

	// Protected API routes - API route structure
	//
	// Public endpoints:     /health, /auth/webhook, etc.
	// Protected endpoints:  /api/users, /api/users/me, etc.
	//
	// All routes under /api/* are automatically protected by auth middleware
	// Each resource gets its own submux registered at /api/{resource}

	// When adding a new resource API:
	// 1. Create a new mux for the resource
	// 2. Register routes to that mux (with paths relative to the resource root)
	// 3. Mount it at /api/{resource} with auth middleware

	// Users API
	usersMux := http.NewServeMux()
	handlers.RegisterUsersRoutes(usersMux, r.logger, r.store)
	r.mux.Handle("/api/users/", http.StripPrefix("/api/users", middleware.AuthMiddleware(r.logger)(usersMux)))

	r.logger.Info("Routes registered successfully")
}
