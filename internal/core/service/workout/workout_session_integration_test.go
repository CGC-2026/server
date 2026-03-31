//go:build integration

package workout

import (
	"context"
	"testing"
	"time"

	"server/internal/testutil/testdb"

	"github.com/stretchr/testify/require"
)

func TestCreateWorkoutSessionPersistsCalibrationValues(t *testing.T) {
	t.Parallel()

	harness := testdb.New(t)
	user := testdb.CreateUser(t, harness.Store, "create_session")
	workoutType := testdb.CreateWorkoutType(t, harness.Store, "create_session")
	testdb.UpsertCalibration(t, harness.Store, user.ID, 1.25, 2.5, 3.75)

	startTime := time.Date(2026, 3, 30, 8, 0, 0, 0, time.UTC)
	session, err := CreateWorkoutSession(context.Background(), harness.Store, user.ID, CreateWorkoutSessionDTO{
		WorkoutTypeID: workoutType.ID,
		StartTime:     &startTime,
	})

	require.NoError(t, err)
	require.Equal(t, user.ID, session.UserID)
	require.Equal(t, workoutType.ID, session.WorkoutTypeID)
	require.True(t, session.CalibrationYawAngle.Valid)
	require.True(t, session.CalibrationPitchAngle.Valid)
	require.True(t, session.CalibrationRollAngle.Valid)
	require.Equal(t, 1.25, session.CalibrationYawAngle.Float64)
	require.Equal(t, 2.5, session.CalibrationPitchAngle.Float64)
	require.Equal(t, 3.75, session.CalibrationRollAngle.Float64)
	require.True(t, session.StartTime.Valid)
	require.True(t, session.StartTime.Time.Equal(startTime))
	require.False(t, session.EndTime.Valid)
}

func TestGetWorkoutSessionByIDReturnsHydratedSetsAndReps(t *testing.T) {
	t.Parallel()

	harness := testdb.New(t)
	user := testdb.CreateUser(t, harness.Store, "get_by_id")
	workoutType := testdb.CreateWorkoutType(t, harness.Store, "get_by_id")

	session, err := CreateWorkoutSession(context.Background(), harness.Store, user.ID, CreateWorkoutSessionDTO{
		WorkoutTypeID: workoutType.ID,
	})
	require.NoError(t, err)

	_, err = CreateWorkoutSessionSet(context.Background(), harness.Store, user.ID, session.ID, CreateWorkoutSessionSetDTO{
		SetNumber: 1,
		Reps: []CreateWorkoutSessionSetRepDTO{
			{
				RepNumber: 1,
				StartTime: 10,
				EndTime:   20,
				Samples:   []byte(`[{"angle": 90}]`),
				Metrics:   []byte(`{"velocity": 1.5}`),
			},
		},
	})
	require.NoError(t, err)

	details, err := GetWorkoutSessionByID(context.Background(), harness.Store, user.ID, session.ID)
	require.NoError(t, err)
	require.Equal(t, session.ID, details.ID)
	require.Len(t, details.Sets, 1)
	require.Equal(t, int32(1), details.Sets[0].SetNumber)
	require.Len(t, details.Sets[0].Reps, 1)
	require.Equal(t, int32(1), details.Sets[0].Reps[0].RepNumber)

	samples, ok := details.Sets[0].Reps[0].Samples.([]any)
	require.True(t, ok)
	require.Len(t, samples, 1)

	metrics, ok := details.Sets[0].Reps[0].Metrics.(map[string]any)
	require.True(t, ok)
	require.Equal(t, 1.5, metrics["velocity"])
}

func TestGetUserWorkoutSessionHistoryReturnsOrderedCounts(t *testing.T) {
	t.Parallel()

	harness := testdb.New(t)
	user := testdb.CreateUser(t, harness.Store, "history")
	workoutType := testdb.CreateWorkoutType(t, harness.Store, "history")

	olderStart := time.Date(2026, 3, 29, 8, 0, 0, 0, time.UTC)
	olderSession, err := CreateWorkoutSession(context.Background(), harness.Store, user.ID, CreateWorkoutSessionDTO{
		WorkoutTypeID: workoutType.ID,
		StartTime:     &olderStart,
	})
	require.NoError(t, err)

	_, err = CreateWorkoutSessionSet(context.Background(), harness.Store, user.ID, olderSession.ID, CreateWorkoutSessionSetDTO{
		SetNumber: 1,
		Reps: []CreateWorkoutSessionSetRepDTO{
			{
				RepNumber: 1,
				StartTime: 1,
				EndTime:   2,
				Samples:   []byte(`[]`),
				Metrics:   []byte(`{}`),
			},
		},
	})
	require.NoError(t, err)

	newerStart := time.Date(2026, 3, 30, 8, 0, 0, 0, time.UTC)
	newerSession, err := CreateWorkoutSession(context.Background(), harness.Store, user.ID, CreateWorkoutSessionDTO{
		WorkoutTypeID: workoutType.ID,
		StartTime:     &newerStart,
	})
	require.NoError(t, err)

	_, err = CreateWorkoutSessionSet(context.Background(), harness.Store, user.ID, newerSession.ID, CreateWorkoutSessionSetDTO{
		SetNumber: 1,
		Reps: []CreateWorkoutSessionSetRepDTO{
			{
				RepNumber: 1,
				StartTime: 3,
				EndTime:   4,
				Samples:   []byte(`[]`),
				Metrics:   []byte(`{}`),
			},
			{
				RepNumber: 2,
				StartTime: 5,
				EndTime:   6,
				Samples:   []byte(`[]`),
				Metrics:   []byte(`{}`),
			},
		},
	})
	require.NoError(t, err)

	_, err = CreateWorkoutSessionSet(context.Background(), harness.Store, user.ID, newerSession.ID, CreateWorkoutSessionSetDTO{
		SetNumber: 2,
		Reps: []CreateWorkoutSessionSetRepDTO{
			{
				RepNumber: 1,
				StartTime: 7,
				EndTime:   8,
				Samples:   []byte(`[]`),
				Metrics:   []byte(`{}`),
			},
		},
	})
	require.NoError(t, err)

	history, err := GetUserWorkoutSessionHistory(context.Background(), harness.Store, user.ID, 10)
	require.NoError(t, err)
	require.Len(t, history, 2)

	require.Equal(t, newerSession.ID, history[0].ID)
	require.Equal(t, int32(2), history[0].TotalSets)
	require.Equal(t, int32(3), history[0].TotalReps)

	require.Equal(t, olderSession.ID, history[1].ID)
	require.Equal(t, int32(1), history[1].TotalSets)
	require.Equal(t, int32(1), history[1].TotalReps)
}
