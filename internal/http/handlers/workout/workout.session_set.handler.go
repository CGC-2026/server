package workout

import (
	"encoding/json"
	"net/http"
	service "server/internal/core/service/workout"
	"server/internal/data"
	"server/internal/http/middleware"
	"server/internal/http/responses"
	"server/internal/log"
)

func RegisterWorkoutSessionSetRoutes(mux *http.ServeMux, logger log.Logger, store *data.Store) {

	// POST /api/workouts/sessions - Create a new workout session
	mux.HandleFunc("POST /sessions/{id}/sets", func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetDBUserIDFromClerkID(store, r)
		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"}, logger)
			return
		}

		sessionID := r.PathValue("id")
		if sessionID == "" {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid Session ID"}, logger)
			return
		}

		var dto service.CreateWorkoutSessionSetDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"}, logger)
			return
		}

		setRow, err := service.CreateWorkoutSessionSet(r.Context(), store, userID, sessionID, dto)
		if err != nil {
			logger.Error("failed to create workout session set: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not save set"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusCreated, setRow, logger)
	})

}
