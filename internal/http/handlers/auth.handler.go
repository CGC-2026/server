package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"server/internal/core/service"
	"server/internal/data"
	"server/internal/http/responses"
	"server/internal/log"
	"strconv"
	"time"

	svix "github.com/svix/svix-webhooks/go"
)

// ClerkWebhookEvent represents the structure of a Clerk webhook event
type ClerkWebhookEvent struct {
	Type   string          `json:"type"`
	Object string          `json:"object"`
	Data   json.RawMessage `json:"data"`
}

func RegisterAuthRoutes(mux *http.ServeMux, logger log.Logger, store *data.Store) {
	// Clerk webhook handler for user events
	mux.HandleFunc("POST /webhooks/clerk", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Failed to read webhook body: %v", err)
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"}, logger)
			return
		}
		defer r.Body.Close()

		secret := os.Getenv("CLERK_WEBHOOK_SECRET") // must be the Signing secret: whsec_...
		if secret == "" {
			logger.Error("CLERK_WEBHOOK_SECRET not set")
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "server misconfigured"}, logger)
			return
		}

		if r.Header.Get("svix-id") == "" || r.Header.Get("svix-timestamp") == "" || r.Header.Get("svix-signature") == "" {
			logger.Error("Missing webhook signature headers from request")
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "missing signature headers"}, logger)
			return
		}

		if ts, err := strconv.ParseInt(r.Header.Get("svix-timestamp"), 10, 64); err == nil {
			now := time.Now().Unix()
			if abs := now - ts; abs > 300 || abs < -300 {
				logger.Error("Webhook timestamp outside 5 min skew. now=%d ts=%d", now, ts)
				responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "stale or future timestamp"}, logger)
				return
			}
		}

		wh, err := svix.NewWebhook(secret)
		if err != nil {
			logger.Error("svix.NewWebhook: %v", err)
			responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "server error"}, logger)
			return
		}

		if err := wh.Verify(body, r.Header); err != nil {
			logger.Error("Invalid Svix signature: %v", err)
			responses.JSONResponse(w, http.StatusUnauthorized, map[string]string{"error": "invalid signature"}, logger)
			return
		}

		var event ClerkWebhookEvent
		if err := json.Unmarshal(body, &event); err != nil {
			logger.Error("parse event: %v", err)
			responses.JSONResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid event data"}, logger)
			return
		}
		logger.Info("Clerk webhook ok: %s", event.Type)

		// Handle different event types
		switch event.Type {
		case "user.created":
			if err := service.HandleCreateUserFromClerk(r.Context(), store, event.Data); err != nil {
				logger.Error("Failed to handle user.created: %v", err)
				responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"}, logger)
				return
			}
			logger.Info("Created user from Clerk webhook")
		case "user.deleted":
			if err := service.HandleDeleteUserFromClerk(r.Context(), store, event.Data); err != nil {
				logger.Error("Failed to handle user.deleted: %v", err)
				responses.JSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete user"}, logger)
				return
			}
			logger.Info("Deleted user from Clerk webhook")
		default:
			logger.Info("Unhandled webhook event type: %s", event.Type)
		}

		responses.JSONResponse(w, http.StatusOK, map[string]string{"status": "success"}, logger)
	})
}
