package handlers

import (
	"net/http"
	"time"

	"server/internal/http/responses"
	"server/internal/log"
)

// Response represents the response from the health endpoint
type Response struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// RegisterRoutes registers health check routes
func RegisterHealthRoutes(mux *http.ServeMux, logger log.Logger) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Health check requested")

		response := Response{
			Status:      "OK",
			Description: "The server is running",
			Timestamp:   time.Now().Format(time.RFC3339),
		}

		responses.JSONResponse(w, http.StatusOK, response, logger)
	})
}
