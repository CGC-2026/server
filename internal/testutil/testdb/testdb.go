package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"server/internal/data"

	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

var nonSchemaChars = regexp.MustCompile(`[^a-z0-9_]+`)

type Harness struct {
	Store       *data.Store
	SQLDB       *sql.DB
	Schema      string
	DatabaseURL string
}

type resolvedTestDatabaseURL struct {
	URL    string
	Source string
}

func New(t testing.TB) *Harness {
	t.Helper()

	resolved := mustResolveTestDatabaseURL(t)

	adminDB := openSQLDB(t, resolved, "integration test admin database")
	schema := testSchemaName(t.Name())
	createSchema(t, adminDB, schema)

	isolatedURL := urlWithSearchPath(t, resolved.URL, schema)
	sqlDB := openSQLDB(t, resolvedTestDatabaseURL{
		URL:    isolatedURL,
		Source: resolved.Source,
	}, "integration test schema")

	require.NoError(t, goose.SetDialect("postgres"))
	require.NoError(t, goose.Up(sqlDB, migrationsDir(t)))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := data.NewStoreFromURL(ctx, isolatedURL)
	require.NoError(t, err)

	harness := &Harness{
		Store:       store,
		SQLDB:       sqlDB,
		Schema:      schema,
		DatabaseURL: isolatedURL,
	}

	t.Cleanup(func() {
		if harness.Store != nil {
			harness.Store.Close()
		}
		if harness.SQLDB != nil {
			require.NoError(t, harness.SQLDB.Close())
		}
		_, err := adminDB.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema))
		require.NoError(t, err)
		require.NoError(t, adminDB.Close())
	})

	return harness
}

func mustResolveTestDatabaseURL(t testing.TB) resolvedTestDatabaseURL {
	t.Helper()

	root := repoRoot(t)
	resolved, err := resolveTestDatabaseURL(root)
	require.NoError(t, err)
	return resolved
}

func resolveTestDatabaseURL(root string) (resolvedTestDatabaseURL, error) {
	if databaseURL := os.Getenv("TEST_DATABASE_URL"); databaseURL != "" {
		return resolvedTestDatabaseURL{
			URL:    databaseURL,
			Source: "process environment",
		}, nil
	}

	envPath := filepath.Join(root, ".env")
	envMap, err := godotenv.Read(envPath)
	if err != nil {
		return resolvedTestDatabaseURL{}, fmt.Errorf("TEST_DATABASE_URL is not set in the process environment, and failed to load %s: %w", envPath, err)
	}

	databaseURL := envMap["TEST_DATABASE_URL"]
	if databaseURL == "" {
		return resolvedTestDatabaseURL{}, fmt.Errorf("TEST_DATABASE_URL is not set in the process environment and not found in %s", envPath)
	}

	return resolvedTestDatabaseURL{
		URL:    databaseURL,
		Source: envPath,
	}, nil
}

func createSchema(t testing.TB, db *sql.DB, schema string) {
	t.Helper()

	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schema))
	require.NoError(t, err)
}

func openSQLDB(t testing.TB, resolved resolvedTestDatabaseURL, purpose string) *sql.DB {
	t.Helper()

	db, err := sql.Open("pgx", resolved.URL)
	require.NoErrorf(t, err, "failed to open %s using TEST_DATABASE_URL from %s (%s)", purpose, resolved.Source, databaseTarget(resolved.URL))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoErrorf(t, db.PingContext(ctx), "failed to connect to %s using TEST_DATABASE_URL from %s (%s)%s", purpose, resolved.Source, databaseTarget(resolved.URL), connectionHint(resolved.URL))
	return db
}

func migrationsDir(t testing.TB) string {
	t.Helper()

	root := repoRoot(t)
	return filepath.Join(root, "db", "migrations")
}

func repoRoot(t testing.TB) string {
	t.Helper()

	dir, err := os.Getwd()
	require.NoError(t, err)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("failed to locate repository root from test working directory")
		}

		dir = parent
	}
}

func testSchemaName(testName string) string {
	base := strings.ToLower(nonSchemaChars.ReplaceAllString(testName, "_"))
	base = strings.Trim(base, "_")
	if base == "" {
		base = "integration"
	}
	if len(base) > 24 {
		base = base[:24]
	}

	return fmt.Sprintf("it_%s_%d", base, time.Now().UnixNano())
}

func urlWithSearchPath(t testing.TB, databaseURL string, schema string) string {
	t.Helper()

	parsed, err := url.Parse(databaseURL)
	require.NoError(t, err)

	query := parsed.Query()
	query.Set("search_path", schema)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

func databaseTarget(databaseURL string) string {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return "unparseable TEST_DATABASE_URL"
	}

	parts := make([]string, 0, 4)
	if parsed.User != nil && parsed.User.Username() != "" {
		parts = append(parts, "user="+parsed.User.Username())
	}
	if host := parsed.Hostname(); host != "" {
		parts = append(parts, "host="+host)
	}
	if port := parsed.Port(); port != "" {
		parts = append(parts, "port="+port)
	}
	if dbName := strings.TrimPrefix(parsed.Path, "/"); dbName != "" {
		parts = append(parts, "db="+dbName)
	}

	if len(parts) == 0 {
		return "TEST_DATABASE_URL"
	}

	return strings.Join(parts, ", ")
}

func connectionHint(databaseURL string) string {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return ""
	}

	switch parsed.Hostname() {
	case "db":
		return " Hint: hostname `db` only resolves inside Docker Compose. For local integration tests, set TEST_DATABASE_URL to use 127.0.0.1 and the configured TEST_POSTGRES_PORT."
	case "127.0.0.1", "localhost", "::1":
		return " If this is a local host-side run, confirm the Postgres container is up and listening on the configured TEST_POSTGRES_PORT."
	default:
		return ""
	}
}

