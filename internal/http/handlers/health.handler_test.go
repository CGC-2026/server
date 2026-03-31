//go:build !integration

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestRegisterHealthRoutesReturnsExpectedResponse(t *testing.T) {
	t.Parallel()

	logger := testutil.NewTestLogger()
	mux := http.NewServeMux()
	RegisterHealthRoutes(mux, logger)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	var response Response
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	require.Equal(t, "OK", response.Status)
	require.Equal(t, "The server is running", response.Description)
	_, err := time.Parse(time.RFC3339, response.Timestamp)
	require.NoError(t, err)
}
