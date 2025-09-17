package router

import (
	"net/http"
	handlers "server/internal/app/handlers"
	"server/internal/log"
)

// Router manages all HTTP routes for the application
type Router struct {
	mux    *http.ServeMux
	logger log.Logger
}

// New creates a new router instance
func New(logger log.Logger) *Router {
	return &Router{
		mux:    http.NewServeMux(),
		logger: logger,
	}
}

// Handler returns the http.Handler for the router
func (r *Router) Handler() http.Handler {
	return r.mux
}

// RegisterRoutes registers all application routes
func (r *Router) RegisterRoutes() {
	handlers.RegisterHealthRoutes(r.mux, r.logger)

	r.logger.Info("Routes registered successfully")
}
