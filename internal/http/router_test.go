//go:build !integration

package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/http/handlers"
	"server/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestRouterHandlerServesHealthRoute(t *testing.T) {
	t.Parallel()

	logger := testutil.NewTestLogger()
	router := New(logger, nil)
	router.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.Handler().ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	var response handlers.Response
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	require.Equal(t, "OK", response.Status)

	entry, ok := logger.LastEntry()
	require.True(t, ok)
	require.Equal(t, "info", entry.Level)
	require.Contains(t, entry.Message, "GET /health 200")
}
