package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"server/internal/core/service"
	"server/internal/data"
	"server/internal/http/middleware"
	"server/internal/http/responses"
	"server/internal/log"

	"github.com/jackc/pgx/v5"
)

func RegisterUsersRoutes(mux *http.ServeMux, logger log.Logger, store *data.Store) {
	// Get current user profile (protected)
	mux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetDBUserIDFromClerkID(store, r)
		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, logger)
			return
		}

		user, err := service.GetCurrentUserProfile(r.Context(), store, userID)
		if err != nil {
			logger.Error("GetUserByID: %v", err)
			responses.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "user not found"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusOK, user, logger)
	})

	// GET /api/users/me/calibration - get current user's calibration data
	mux.HandleFunc("GET /me/calibration", func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetDBUserIDFromClerkID(store, r)
		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, logger)
			return
		}

		calibration, err := service.GetCurrentUserCalibration(r.Context(), store, userID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				responses.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "no calibration data found for user"}, logger)
				return
			}
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user calibration"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusOK, calibration, logger)
	})

	// POST /api/users/me/calibration - upsert current user's calibration data
	mux.HandleFunc("POST /me/calibration", func(w http.ResponseWriter, r *http.Request) {
		userId, err := middleware.GetDBUserIDFromClerkID(store, r)
		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, logger)
			return
		}

		var request service.UpsertUserCalibrationDTO
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"}, logger)
			return
		}

		calibration, err := service.UpsertUserCalibration(r.Context(), store, userId, request)
		if err != nil {
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to upsert user calibration"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusOK, calibration, logger)
	})

	// DELETE /api/users/me/calibration - delete current user's calibration data
	mux.HandleFunc("DELETE /me/calibration", func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetDBUserIDFromClerkID(store, r)
		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, logger)
			return
		}

		if err := service.DeleteUserCalibration(r.Context(), store, userID); err != nil {
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete user calibration"}, logger)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
