//go:build integration

package service

import (
	"context"
	"encoding/json"
	"testing"

	"server/internal/testutil/testdb"

	"github.com/stretchr/testify/require"
)

func TestHandleCreateUserFromClerkPersistsExpectedFields(t *testing.T) {
	harness := testdb.New(t)

	payload := json.RawMessage(`{
		"id": "clerk_user_123",
		"first_name": "Ada",
		"last_name": "Lovelace",
		"email_addresses": [
			{"email_address": "ada@example.com", "id": "email_1"},
			{"email_address": "other@example.com", "id": "email_2"}
		],
		"image_url": "https://example.com/ada.png"
	}`)

	err := HandleCreateUserFromClerk(context.Background(), harness.Store, payload)
	require.NoError(t, err)

	user, err := harness.Store.Queries.GetUserByClerkID(context.Background(), "clerk_user_123")
	require.NoError(t, err)
	require.Equal(t, "Ada", user.FirstName)
	require.Equal(t, "Lovelace", user.LastName)
	require.Equal(t, "ada@example.com", user.Email)
	require.True(t, user.ImageUrl.Valid)
	require.Equal(t, "https://example.com/ada.png", user.ImageUrl.String)
}
