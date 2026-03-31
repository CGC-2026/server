//go:build !integration

package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestAuthMiddlewareRejectsMissingAuthorizationHeader(t *testing.T) {
	t.Parallel()

	logger := testutil.NewTestLogger()
	handler := AuthMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)

	var body map[string]string
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "missing authorization header", body["error"])
}

func TestAuthMiddlewareRejectsInvalidBearerFormat(t *testing.T) {
	t.Parallel()

	logger := testutil.NewTestLogger()
	handler := AuthMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	req.Header.Set("Authorization", "Token invalid")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusUnauthorized, recorder.Code)

	var body map[string]string
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&body))
	require.Equal(t, "invalid authorization header format", body["error"])
}

func TestAuthMiddlewareAllowsDevelopmentBypass(t *testing.T) {
	logger := testutil.NewTestLogger()

	t.Setenv("APP_ENV", "development")
	t.Setenv("DEV_BYPASS_SECRET", "local-secret")

	handler := AuthMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r.Context())
		require.True(t, ok)
		require.Equal(t, "user-dev-123", userID)
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	req.Header.Set("X-Dev-Secret", "local-secret")
	req.Header.Set("X-Dev-User-Id", "user-dev-123")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusNoContent, recorder.Code)

	entry, ok := logger.LastEntry()
	require.True(t, ok)
	require.Equal(t, "info", entry.Level)
	require.Contains(t, entry.Message, "Development mode: bypassing auth for user")
}
