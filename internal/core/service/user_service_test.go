//go:build !integration

package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandleCreateUserFromClerkReturnsErrorOnMalformedJSON(t *testing.T) {
	t.Parallel()

	err := HandleCreateUserFromClerk(context.Background(), nil, json.RawMessage(`{"id":`))

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse user data")
}
