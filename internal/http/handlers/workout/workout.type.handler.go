package workout

import (
	"net/http"
	service "server/internal/core/service/workout"
	"server/internal/data"
	"server/internal/http/responses"
	"server/internal/log"
)

func RegisterWorkoutTypeRoutes(mux *http.ServeMux, logger log.Logger, store *data.Store) {
	// GET /api/workouts/types - Retrieve all workout types
	mux.HandleFunc("GET /types", func(w http.ResponseWriter, r *http.Request) {
		types, err := service.GetWorkoutTypes(r.Context(), store)
		if err != nil {
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Could not fetch workout types"}, logger)
			return
		}
		responses.JSONResponse(w, http.StatusOK, types, logger)
	})

}
