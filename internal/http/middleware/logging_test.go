//go:build !integration

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestStatusTextCategory(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		status   int
		expected string
	}{
		{name: "unknown", status: 0, expected: "unknown"},
		{name: "informational", status: 102, expected: "informational"},
		{name: "success", status: 200, expected: "success"},
		{name: "redirect", status: 302, expected: "redirect"},
		{name: "client error", status: 404, expected: "client_error"},
		{name: "server error", status: 503, expected: "server_error"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, statusTextCategory(tc.status))
		})
	}
}

func TestLoggingMiddlewareLogsSuccessfulResponses(t *testing.T) {
	t.Parallel()

	logger := testutil.NewTestLogger()
	handler := LoggingMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/health?check=1", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code)

	entry, ok := logger.LastEntry()
	require.True(t, ok)
	require.Equal(t, "info", entry.Level)
	require.Contains(t, entry.Message, "GET /health?check=1 200")
	require.Contains(t, entry.Message, "5B")
	require.Contains(t, entry.Message, "category: success")
}

func TestLoggingMiddlewareLogsErrorResponses(t *testing.T) {
	t.Parallel()

	logger := testutil.NewTestLogger()
	handler := LoggingMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("oops"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/workouts/sessions", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusServiceUnavailable, recorder.Code)

	entry, ok := logger.LastEntry()
	require.True(t, ok)
	require.Equal(t, "error", entry.Level)
	require.Contains(t, entry.Message, "POST /api/workouts/sessions 503")
	require.Contains(t, entry.Message, "4B")
	require.Contains(t, entry.Message, "category: server_error")
}
