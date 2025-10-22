package handlers

import (
	"net/http"

	"server/internal/data"
	"server/internal/http/middleware"
	"server/internal/http/responses"
	"server/internal/log"
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

func RegisterUsersRoutes(mux *http.ServeMux, logger log.Logger, store *data.Store) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		users, err := store.Queries.ListUsers(r.Context())
		if err != nil {
			logger.Error("ListUsers: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to list users"}, logger)
			return
		}

		out := make([]userDTO, 0, len(users))
		for _, u := range users {
			dto := userDTO{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
				CreatedAt: u.CreatedAt.Time.Format("YYYY-MM-DDTHH:MM:SSZ"),
				UpdatedAt: u.UpdatedAt.Time.Format("YYYY-MM-DDTHH:MM:SSZ"),
			}
			if u.ImageUrl.Valid {
				dto.ImageURL = u.ImageUrl.String
			}
			out = append(out, dto)
		}
		responses.JSONResponse(w, http.StatusOK, out, logger)
	})

	// Get current user profile (protected)
	mux.HandleFunc("GET /me", func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.GetUserID(r.Context())
		if !ok {
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"}, logger)
			return
		}

		user, err := store.Queries.GetUserByID(r.Context(), userID)
		if err != nil {
			logger.Error("GetUserByID: %v", err)
			responses.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "user not found"}, logger)
			return
		}

		dto := userDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Time.Format("YYYY-MM-DDTHH:MM:SSZ"),
			UpdatedAt: user.UpdatedAt.Time.Format("YYYY-MM-DDTHH:MM:SSZ"),
		}
		if user.ImageUrl.Valid {
			dto.ImageURL = user.ImageUrl.String
		}

		responses.JSONResponse(w, http.StatusOK, dto, logger)
	})

	// Get user by ID (protected)
	mux.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "user ID is required"}, logger)
			return
		}

		user, err := store.Queries.GetUserByID(r.Context(), id)
		if err != nil {
			logger.Error("GetUserByID: %v", err)
			responses.JSONResponse(w, http.StatusNotFound, map[string]string{"error": "user not found"}, logger)
			return
		}

		dto := userDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Time.Format("YYYY-MM-DDTHH:MM:SSZ"),
			UpdatedAt: user.UpdatedAt.Time.Format("YYYY-MM-DDTHH:MM:SSZ"),
		}
		if user.ImageUrl.Valid {
			dto.ImageURL = user.ImageUrl.String
		}

		responses.JSONResponse(w, http.StatusOK, dto, logger)
	})
}
