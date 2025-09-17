package app

import (
	"net/http"
	"server/internal/http/responses"
	"server/internal/log"
	"time"
)

// HealthResponse represents the response from the health endpoint
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// RegisterRoutes registers all routes for the application
func RegisterRoutes(mux *http.ServeMux, logger log.Logger) {
	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, r, logger)
	})

	logger.Info("Routes registered successfully")
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, _ *http.Request, logger log.Logger) {
	logger.Info("Health check requested")

	response := HealthResponse{
		Status:    "OK",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	responses.JSONResponse(w, http.StatusOK, response, logger)
}
