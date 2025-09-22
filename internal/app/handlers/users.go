package handlers

import (
	"encoding/json"
	"net/http"

	"server/internal/data"
	"server/internal/http/responses"
	"server/internal/log"

	sqlc "server/internal/db"
)

type (
	userDTO struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	createUserDTO struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)

func RegisterUsersRoutes(mux *http.ServeMux, logger log.Logger, store *data.Store) {
	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		users, err := store.Queries.ListUsers(r.Context())
		if err != nil {
			logger.Error("ListUsers: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to list users"}, logger)
			return
		}

		out := make([]userDTO, 0, len(users))
		for _, u := range users {
			out = append(out, userDTO{
				ID: u.ID,
				FirstName: u.FirstName,
				LastName: u.LastName,
				Email: u.Email,
			})
		}
		responses.JSONResponse(w, http.StatusOK, out, logger)
	})

	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var in createUserDTO
		dec := jsonNewDecoder(r)
		if err := dec.Decode(&in); err != nil {
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON: " + err.Error()}, logger)
			return
		}
		u, err := store.Queries.CreateUser(r.Context(), sqlc.CreateUserParams{
			FirstName: in.FirstName,
			LastName: in.LastName,
			Email: in.Email,
		})
		if err != nil {
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"}, logger)
			return
		}
		responses.JSONResponse(w, http.StatusCreated, userDTO{
			ID: u.ID,
			FirstName: u.FirstName,
			LastName: u.LastName,
			Email: u.Email,
		}, logger)
	})
}


func jsonNewDecoder(r *http.Request) *json.Decoder {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec
}