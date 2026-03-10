package workout

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/core/service/helpers"
	"server/internal/data"
	"server/internal/db"
	"time"
)

type CreateWorkoutSessionSetRepDTO struct {
	RepNumber int32           `json:"rep_number"`
	StartTime int64           `json:"start_time"`
	EndTime   int64           `json:"end_time"`
	Samples   json.RawMessage `json:"samples"`
	Metrics   json.RawMessage `json:"metrics"`
}

type CreateWorkoutSessionSetDTO struct {
	SetNumber int32                           `json:"set_number"`
	StartTime *time.Time                      `json:"start_time"`
	EndTime   *time.Time                      `json:"end_time"`
	Reps      []CreateWorkoutSessionSetRepDTO `json:"reps"`
}

func CreateWorkoutSessionSet(ctx context.Context, store *data.Store, userID string, sessionID string, dto CreateWorkoutSessionSetDTO) (db.WorkoutSessionSet, error) {
	session, err := store.Queries.GetWorkoutSessionByID(ctx, sessionID)
	if err != nil {
		return db.WorkoutSessionSet{}, err
	}

	if session.UserID != userID {
		return db.WorkoutSessionSet{}, fmt.Errorf("forbidden")
	}

	if dto.SetNumber <= 0 {
		return db.WorkoutSessionSet{}, fmt.Errorf("set_number must be greater than 0")
	}

	var createdSet db.WorkoutSessionSet

	startTime := time.Now()
	if dto.StartTime != nil && !dto.StartTime.IsZero() {
		startTime = *dto.StartTime
	}
	endTime := time.Now()
	if dto.EndTime != nil && !dto.EndTime.IsZero() {
		endTime = *dto.EndTime
	}

	err = store.WithTransaction(ctx, func(q *db.Queries) error {
		setRow, err := q.CreateWorkoutSessionSet(ctx, db.CreateWorkoutSessionSetParams{
			WorkoutSessionID: sessionID,
			SetNumber:        dto.SetNumber,
			StartTime:        helpers.NewNullTime(&startTime),
			EndTime:          helpers.NewNullTime(&endTime),
		})
		if err != nil {
			return err
		}

		for _, rep := range dto.Reps {
			if rep.RepNumber <= 0 {
				return fmt.Errorf("rep_number must be greater than 0")
			}
			if rep.EndTime < rep.StartTime {
				return fmt.Errorf("rep end_time must be after start_time")
			}
			if !json.Valid(rep.Samples) {
				return fmt.Errorf("samples must be valid JSON")
			}
			if !json.Valid(rep.Metrics) {
				return fmt.Errorf("metrics must be valid JSON")
			}

			_, err := q.CreateWorkoutSessionSetRep(ctx, db.CreateWorkoutSessionSetRepParams{
				WorkoutSessionSetID: setRow.ID,
				RepNumber:           rep.RepNumber,
				StartTime:           rep.StartTime,
				EndTime:             rep.EndTime,
				Samples:             []byte(rep.Samples),
				Metrics:             []byte(rep.Metrics),
			})
			if err != nil {
				return err
			}
		}

		createdSet = setRow
		return nil
	})

	if err != nil {
		return db.WorkoutSessionSet{}, err
	}

	return createdSet, nil
}
