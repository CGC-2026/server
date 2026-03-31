package testdb

import (
	"context"
	"fmt"
	"testing"
	"time"

	"server/internal/data"
	sqlc "server/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func CreateUser(t testing.TB, store *data.Store, suffix string) sqlc.User {
	t.Helper()

	if suffix == "" {
		suffix = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	user, err := store.Queries.CreateUser(context.Background(), sqlc.CreateUserParams{
		ClerkID:   fmt.Sprintf("clerk_%s", suffix),
		FirstName: "Test",
		LastName:  "User",
		Email:     fmt.Sprintf("test_%s@example.com", suffix),
		ImageUrl:  pgtype.Text{},
	})
	require.NoError(t, err)

	return user
}

func CreateWorkoutType(t testing.TB, store *data.Store, suffix string) sqlc.WorkoutType {
	t.Helper()

	if suffix == "" {
		suffix = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	workoutType, err := store.Queries.CreateWorkoutType(context.Background(), sqlc.CreateWorkoutTypeParams{
		Name:        fmt.Sprintf("Workout %s", suffix),
		Description: pgtype.Text{String: "Integration test workout", Valid: true},
		Config:      []byte(`{"sampleRate":120}`),
		ImageUrl:    pgtype.Text{},
	})
	require.NoError(t, err)

	return workoutType
}

func UpsertCalibration(t testing.TB, store *data.Store, userID string, yaw float64, pitch float64, roll float64) sqlc.UserCalibration {
	t.Helper()

	calibration, err := store.Queries.UpsertUserCalibration(context.Background(), sqlc.UpsertUserCalibrationParams{
		UserID:             userID,
		StandingYawAngle:   yaw,
		StandingPitchAngle: pitch,
		StandingRollAngle:  roll,
	})
	require.NoError(t, err)

	return calibration
}
