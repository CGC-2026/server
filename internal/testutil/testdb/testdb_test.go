//go:build !integration

package testdb

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveTestDatabaseURLPrefersProcessEnvironment(t *testing.T) {
	t.Setenv("TEST_DATABASE_URL", "postgres://from-process")

	root := t.TempDir()
	envPath := filepath.Join(root, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("TEST_DATABASE_URL=postgres://from-dotenv\n"), 0o644))

	resolved, err := resolveTestDatabaseURL(root)
	require.NoError(t, err)
	require.Equal(t, "postgres://from-process", resolved.URL)
	require.Equal(t, "process environment", resolved.Source)
}

func TestResolveTestDatabaseURLLoadsRepoDotEnv(t *testing.T) {
	require.NoError(t, os.Unsetenv("TEST_DATABASE_URL"))
	t.Cleanup(func() {
		_ = os.Unsetenv("TEST_DATABASE_URL")
	})

	root := t.TempDir()
	envPath := filepath.Join(root, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("TEST_DATABASE_URL=postgres://from-dotenv\n"), 0o644))

	resolved, err := resolveTestDatabaseURL(root)
	require.NoError(t, err)
	require.Equal(t, "postgres://from-dotenv", resolved.URL)
	require.Equal(t, envPath, resolved.Source)
}

func TestResolveTestDatabaseURLFailsWhenMissingEverywhere(t *testing.T) {
	require.NoError(t, os.Unsetenv("TEST_DATABASE_URL"))
	t.Cleanup(func() {
		_ = os.Unsetenv("TEST_DATABASE_URL")
	})

	root := t.TempDir()
	envPath := filepath.Join(root, ".env")
	require.NoError(t, os.WriteFile(envPath, []byte("DATABASE_URL=postgres://app-db\n"), 0o644))

	resolved, err := resolveTestDatabaseURL(root)
	require.Error(t, err)
	require.Empty(t, resolved.URL)
	require.Contains(t, err.Error(), "TEST_DATABASE_URL is not set in the process environment and not found in")
	require.Contains(t, err.Error(), envPath)
}

func TestResolveTestDatabaseURLFailsWhenRepoDotEnvCannotBeLoaded(t *testing.T) {
	require.NoError(t, os.Unsetenv("TEST_DATABASE_URL"))
	t.Cleanup(func() {
		_ = os.Unsetenv("TEST_DATABASE_URL")
	})

	root := t.TempDir()

	resolved, err := resolveTestDatabaseURL(root)
	require.Error(t, err)
	require.Empty(t, resolved.URL)
	require.Contains(t, err.Error(), "TEST_DATABASE_URL is not set in the process environment, and failed to load")
	require.Contains(t, err.Error(), filepath.Join(root, ".env"))
}

func TestConnectionHintForComposeHostname(t *testing.T) {
	hint := connectionHint("postgres://capstone:capstone@db:5432/capstone?sslmode=disable")
	require.Contains(t, hint, "hostname `db` only resolves inside Docker Compose")
	require.Contains(t, hint, "TEST_POSTGRES_PORT")
}

func TestConnectionHintForHostSideURL(t *testing.T) {
	hint := connectionHint("postgres://capstone:capstone@127.0.0.1:5432/capstone?sslmode=disable")
	require.Contains(t, hint, "local host-side run")
	require.Contains(t, hint, "TEST_POSTGRES_PORT")
}

