package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"server/internal/data"
	"server/internal/db"
	"server/internal/http/middleware"
	"server/internal/http/responses"
	"server/internal/log"

	"github.com/jackc/pgx/v5"
)

type (
	userDTO struct {
		ID        string `json:"id"`
		ClerkID   string `json:"clerk_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		ImageURL  string `json:"image_url,omitempty"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)

type (
	calibrationResponseDTO struct {
		UserID             string  `json:"user_id"`
		StandingYawAngle   float64 `json:"standing_yaw_angle"`
		StandingPitchAngle float64 `json:"standing_pitch_angle"`
		StandingRollAngle  float64 `json:"standing_roll_angle"`
		UpdatedAt          string  `json:"updated_at"`
	}
)

type (
	calibrationRequestDTO struct {
		StandingYawAngle   float64 `json:"standing_yaw_angle"`
		StandingPitchAngle float64 `json:"standing_pitch_angle"`
		StandingRollAngle  float64 `json:"standing_roll_angle"`
	}
)

func dtoFromCalibrationRow(row db.UserCalibration) calibrationResponseDTO {
	updated := ""
	if row.UpdatedAt.Valid {
		updated = row.UpdatedAt.Time.Format(time.RFC3339)
	}
	return calibrationResponseDTO{
		UserID:             row.UserID,
		StandingYawAngle:   row.StandingYawAngle,
		StandingPitchAngle: row.StandingPitchAngle,
		StandingRollAngle:  row.StandingRollAngle,
		UpdatedAt:          updated,
	}
}

func RegisterUsersRoutes(mux *http.ServeMux, logger log.Logger, store *data.Store) {
	// Get current user profile (protected)
	mux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		userID, err := middleware.GetDBUserIDFromClerkID(store, r)
		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, logger)
			return
		}

		user, err := store.Queries.GetUserByID(r.Context(), userID)
		if err != nil {
			logger.Error("GetUserByID: %v", err)
			responses.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "user not found"}, logger)
			return
		}

		created := ""
		if user.CreatedAt.Valid {
			created = user.CreatedAt.Time.Format(time.RFC3339)
		}

		updated := ""
		if user.UpdatedAt.Valid {
			updated = user.UpdatedAt.Time.Format(time.RFC3339)
		}

		dto := userDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: created,
			UpdatedAt: updated,
		}
		if user.ImageUrl.Valid {
			dto.ImageURL = user.ImageUrl.String
		}

		responses.JSONResponse(w, http.StatusOK, dto, logger)
	})

	// GET /api/users/me/calibration - get current user's calibration data
	mux.HandleFunc("GET /me/calibration", func(w http.ResponseWriter, r *http.Request) {
		userId, err := middleware.GetDBUserIDFromClerkID(store, r)
		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, logger)
			return
		}

		calibration, err := store.Queries.GetUserCalibrationByUserID(r.Context(), userId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				responses.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "no calibration data found for user"}, logger)
				return
			}
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user calibration"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusOK, dtoFromCalibrationRow(calibration), logger)
	})

	// POST /api/users/me/calibration - upsert current user's calibration data
	mux.HandleFunc("POST /me/calibration", func(w http.ResponseWriter, r *http.Request) {
		userId, err := middleware.GetDBUserIDFromClerkID(store, r)
		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, logger)
			return
		}

		var request calibrationRequestDTO
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"}, logger)
			return
		}

		calibration, err := store.Queries.UpsertUserCalibration(r.Context(), db.UpsertUserCalibrationParams{
			UserID:             userId,
			StandingYawAngle:   request.StandingYawAngle,
			StandingPitchAngle: request.StandingPitchAngle,
			StandingRollAngle:  request.StandingRollAngle,
		})
		if err != nil {
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to upsert user calibration"}, logger)
			return
		}

		responses.JSONResponse(w, http.StatusOK, dtoFromCalibrationRow(calibration), logger)
	})

	// DELETE /api/users/me/calibration - delete current user's calibration data
	mux.HandleFunc("DELETE /me/calibration", func(w http.ResponseWriter, r *http.Request) {
		userId, err := middleware.GetDBUserIDFromClerkID(store, r)
		if err != nil {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, logger)
			return
		}

		if err := store.Queries.DeleteUserCalibrationByUserID(r.Context(), userId); err != nil {
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete user calibration"}, logger)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
