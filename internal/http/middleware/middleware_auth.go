package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"server/internal/http/responses"
	"server/internal/log"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
)

type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware validates Clerk JWT tokens and extracts user information
func AuthMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Development mode bypass (Requires development environment and secondary bypass secret)
			isDev := os.Getenv("APP_ENV") == "development"

			devSecret := os.Getenv("DEV_BYPASS_SECRET")
			providedSecret := r.Header.Get("X-Dev-Secret")
			devUserId := r.Header.Get("X-Dev-User-Id")

			if isDev && devSecret != "" && providedSecret == devSecret && devUserId != "" {
				logger.Info("Development mode: bypassing auth for user %s", devUserId)
				ctx := context.WithValue(r.Context(), UserIDKey, devUserId)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Get the authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "missing authorization header"}, logger)
				return
			}

			// Extract the token (Bearer <token>)
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"}, logger)
				return
			}

			token := parts[1]

			// Verify the token with Clerk
			claims, err := jwt.Verify(r.Context(), &jwt.VerifyParams{
				Token: token,
			})
			if err != nil {
				logger.Error("Token verification failed: %v", err)
				responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"}, logger)
				return
			}

			// Extract user ID from claims
			userID := claims.Subject
			if userID == "" {
				responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid token claims"}, logger)
				return
			}

			// Add user ID to request context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the user ID from the request context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// InitClerk initializes the Clerk client with the secret key
func InitClerk() error {
	secretKey := os.Getenv("CLERK_SECRET_KEY")
	if secretKey == "" {
		return fmt.Errorf("CLERK_SECRET_KEY is not set")
	}
	clerk.SetKey(secretKey)
	return nil
}
