package workout

import (
	"net/http"
	"server/internal/data"
	"server/internal/log"
)

func RegisterWorkoutRoutes(mux *http.ServeMux, logger log.Logger, store *data.Store) {
	RegisterWorkoutTypeRoutes(mux, logger, store)
	RegisterWorkoutSessionRoutes(mux, logger, store)
	RegisterWorkoutSessionSetRoutes(mux, logger, store)
}
