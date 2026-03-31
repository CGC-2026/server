//go:build integration

package workout

import (
	"context"
	"testing"

	"server/internal/testutil/testdb"

	"github.com/stretchr/testify/require"
)

func TestCreateWorkoutSessionSetPersistsSetAndReps(t *testing.T) {
	harness := testdb.New(t)
	user := testdb.CreateUser(t, harness.Store, "create_set")
	workoutType := testdb.CreateWorkoutType(t, harness.Store, "create_set")

	session, err := CreateWorkoutSession(context.Background(), harness.Store, user.ID, CreateWorkoutSessionDTO{
		WorkoutTypeID: workoutType.ID,
	})
	require.NoError(t, err)

	setRow, err := CreateWorkoutSessionSet(context.Background(), harness.Store, user.ID, session.ID, CreateWorkoutSessionSetDTO{
		SetNumber: 1,
		Reps: []CreateWorkoutSessionSetRepDTO{
			{
				RepNumber: 1,
				StartTime: 100,
				EndTime:   150,
				Samples:   []byte(`[{"frame": 1}]`),
				Metrics:   []byte(`{"peak": 10}`),
			},
			{
				RepNumber: 2,
				StartTime: 200,
				EndTime:   260,
				Samples:   []byte(`[{"frame": 2}]`),
				Metrics:   []byte(`{"peak": 12}`),
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, session.ID, setRow.WorkoutSessionID)
	require.Equal(t, int32(1), setRow.SetNumber)

	sets, err := harness.Store.Queries.GetWorkoutSessionSetsBySessionID(context.Background(), session.ID)
	require.NoError(t, err)
	require.Len(t, sets, 1)

	reps, err := harness.Store.Queries.GetWorkoutSessionSetRepsBySessionID(context.Background(), session.ID)
	require.NoError(t, err)
	require.Len(t, reps, 2)
	require.Equal(t, int32(1), reps[0].RepNumber)
	require.Equal(t, int32(2), reps[1].RepNumber)
	require.JSONEq(t, `[{"frame": 1}]`, string(reps[0].Samples))
	require.JSONEq(t, `{"peak": 12}`, string(reps[1].Metrics))
}
