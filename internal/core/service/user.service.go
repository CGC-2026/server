package service

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/data"

	sqlc "server/internal/db"

	"server/internal/types"

	"github.com/jackc/pgx/v5/pgtype"
)

func HandleCreateUserFromClerk(ctx context.Context, store *data.Store, data json.RawMessage) error {
	var userData types.ClerkUserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return fmt.Errorf("failed to parse user data: %w", err)
	}

	// Extract primary email
	email := ""
	if len(userData.EmailAddresses) > 0 {
		email = userData.EmailAddresses[0].EmailAddress
	}

	// Create image URL as nullable text
	var imageURL pgtype.Text
	if userData.ImageURL != "" {
		imageURL = pgtype.Text{String: userData.ImageURL, Valid: true}
	}

	// Create user in database
	_, err := store.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		ClerkID:   userData.ID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     email,
		ImageUrl:  imageURL,
	})

	if err != nil {
		return fmt.Errorf("failed to create user in database: %w", err)
	}

	return nil
}

func HandleDeleteUserFromClerk(ctx context.Context, store *data.Store, data json.RawMessage) error {
	var userData types.ClerkUserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return fmt.Errorf("failed to parse user data: %w", err)
	}

	return store.Queries.DeleteUserByClerkID(ctx, userData.ID)
}
