package responses

import (
	"encoding/json"
	"net/http"

	"server/internal/log"
)

// JSONResponse handles sending JSON responses with appropriate error handling
func JSONResponse(w http.ResponseWriter, statusCode int, data any, logger log.Logger) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
