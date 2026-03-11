package workout

import (
	"encoding/json"
	"net/http"
	service "server/internal/core/service/workout"
	"server/internal/data"
	"server/internal/http/middleware"
	"server/internal/http/responses"
	"server/internal/log"
	"strconv"
)

func RegisterWorkoutSessionRoutes(mux *http.ServeMux, logger log.Logger, store *data.Store) {
	// POST /api/workouts/sessions - Create a new workout session
	mux.HandleFunc("POST /sessions", func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetDBUserIDFromClerkID(store, r)
		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"}, logger)
			return
		}

		var dto service.CreateWorkoutSessionDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"}, logger)
			return
		}

		session, err := service.CreateWorkoutSession(r.Context(), store, userID, dto)
		if err != nil {
			logger.Error("Failed to create session: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not save session"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusCreated, session, logger)
	})

	// GET /api/workouts/sessions/{id} - Retreive a session with its sets and reps
	mux.HandleFunc("GET /sessions/{id}", func(w http.ResponseWriter, r *http.Request) {
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

		session, err := service.GetWorkoutSessionByID(r.Context(), store, userID, sessionID)
		if err != nil {
			logger.Error("Failed to fetch session details: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not fetch session"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusOK, session, logger)
	})

	// PATCH /api/workouts/sessions/{id} - Update an existing workout session
	mux.HandleFunc("PATCH /sessions/{id}", func(w http.ResponseWriter, r *http.Request) {
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

		var dto service.UpdateWorkoutSessionDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"}, logger)
			return
		}

		session, err := service.UpdateWorkoutSession(r.Context(), store, userID, sessionID, dto)
		if err != nil {
			logger.Error("Failed to update session: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not update session"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusOK, session, logger)
	})

	// GET api/workouts/sessions/history - Retrieve workout session history
	mux.HandleFunc("GET /sessions/history", func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetDBUserIDFromClerkID(store, r)

		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"}, logger)
			return
		}

		limit := int32(20)
		if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil {
			limit = int32(l)
		}

		history, err := service.GetUserWorkoutSessionHistory(r.Context(), store, userID, limit)
		if err != nil {
			logger.Error("History fetch error: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not fetch history"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusOK, history, logger)
	})

}
