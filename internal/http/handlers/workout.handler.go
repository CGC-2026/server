package handlers

import (
	"encoding/json"
	"net/http"
	"server/internal/core/service"
	"server/internal/data"
	"server/internal/http/middleware"
	"server/internal/http/responses"
	"server/internal/log"
	"strconv"
)

func RegisterWorkoutRoutes(mux *http.ServeMux, logger log.Logger, store *data.Store) {
	// GET /api/workouts/types - Retrieve all workout types
	mux.HandleFunc("GET /types", func(w http.ResponseWriter, r *http.Request) {
		types, err := service.GetWorkoutTypes(r.Context(), store)
		if err != nil {
			responses.JSONResponse(w, 500, err, logger)
			return
		}
		responses.JSONResponse(w, 200, types, logger)
	})

	// POST /api/workouts/sessions - Create a new workout session
	mux.HandleFunc("POST /sessions", func(w http.ResponseWriter, r *http.Request) {
		var dto service.CreateWorkoutSessionDto
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"}, logger)
			return
		}

		// Ensure user can only create workout session for themselves
		userId, ok := middleware.GetUserIdFromContext(r.Context())
		if !ok {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"}, logger)
			return
		}
		dto.UserId = userId

		session, err := service.CreateWorkoutSession(r.Context(), store, dto)
		if err != nil {
			logger.Error("Failed to create session: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not save session"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusCreated, session, logger)
	})

	// PATCH /api/workouts/sessions/{id} - Update an existing workout session
	mux.HandleFunc("PATCH /sessions/{id}", func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.PathValue("id")

		var dto service.UpdateWorkoutSessionDto
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"}, logger)
			return
		}

		session, err := service.UpdateWorkoutSession(r.Context(), store, sessionId, dto)
		if err != nil {
			logger.Error("Failed to update session: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not update session"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusOK, session, logger)
	})

	// POST /api/workouts/sessions/{id} - Save workout session data (sets, reps, weights, etc.)
	mux.HandleFunc("POST /sessions/{id}", func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.PathValue("id")

		var dto service.SaveWorkoutSetDto
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid JSON"}, logger)
			return
		}
		dto.WorkoutSessionID = sessionId

		err := service.SaveWorkoutSet(r.Context(), store, dto)
		if err != nil {
			logger.Error("Failed to save workout set: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not save workout set"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusCreated, map[string]string{"message": "Workout set saved successfully"}, logger)
	})

	// GET api/workouts/sessions/history - Retrieve workout session history
	mux.HandleFunc("GET /sessions/history", func(w http.ResponseWriter, r *http.Request) {
		userId, ok := middleware.GetUserIdFromContext(r.Context())
		if !ok {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"}, logger)
			return
		}

		limit := int32(20)

		// Ensure string is within 32-bit int range
		limitStr := r.URL.Query().Get("limit")
		if limitStr != "" {
			l, err := strconv.ParseInt(limitStr, 10, 32)
			if err != nil || l <= 0 {
				responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid limit parameter"}, logger)
				return
			}
			limit = int32(l)
		}

		history, err := service.GetUserWorkoutSessionHistory(r.Context(), store, userId, limit)
		if err != nil {
			logger.Error("History fetch error: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not fetch history"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusOK, history, logger)
	})

}
