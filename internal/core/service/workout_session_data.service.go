package service

import (
	"context"
	"encoding/json"
	"server/internal/data"
	"server/internal/db"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type SaveWorkoutSetDto struct {
	WorkoutSessionID string         `json:"workoutSessionId"`
	SetNumber        int32          `json:"setNumber"`
	StartTime        *time.Time     `json:"startTime"`
	EndTime          *time.Time     `json:"endTime"`
	CoachingScore    *float64       `json:"coachingScore"`
	Reps             []CreateRepDto `json:"reps"`
}

type CreateRepDto struct {
	RepNumber  int32           `json:"repNumber"`
	StartTime  *time.Time      `json:"startTime"`
	EndTime    *time.Time      `json:"endTime"`
	Quality    string          `json:"quality"`
	PeakAngle  float64         `json:"peakAngle"`
	AvgFlex    float64         `json:"avgFlex"`
	SensorData json.RawMessage `json:"sensorData"`
}

func SaveWorkoutSet(ctx context.Context, store *data.Store, dto SaveWorkoutSetDto) error {
	return store.WithTransaction(ctx, func(qtx *db.Queries) error {
		set, err := qtx.CreateWorkoutSet(ctx, db.CreateWorkoutSetParams{
			WorkoutSessionID: dto.WorkoutSessionID,
			SetNumber:        dto.SetNumber,
			StartTime: pgtype.Timestamptz{
				Valid: dto.StartTime != nil,
				Time: func() time.Time {
					if dto.StartTime != nil {
						return *dto.StartTime
					}
					return time.Time{}
				}(),
			},
		})
		if err != nil {
			return err
		}

		for _, r := range dto.Reps {
			err := qtx.CreateWorkoutRep(ctx, db.CreateWorkoutRepParams{
				WorkoutSetID: set.ID,
				RepNumber:    r.RepNumber,
				StartTime: pgtype.Timestamptz{
					Valid: dto.StartTime != nil,
					Time: func() time.Time {
						if dto.StartTime != nil {
							return *dto.StartTime
						}
						return time.Time{}
					}(),
				},
				EndTime: pgtype.Timestamptz{
					Valid: dto.EndTime != nil,
					Time: func() time.Time {
						if dto.EndTime != nil {
							return *dto.EndTime
						}
						return time.Time{}
					}(),
				},
				Quality:          pgtype.Text{String: r.Quality, Valid: true},
				PeakAngle:        pgtype.Float8{Float64: r.PeakAngle, Valid: true},
				AvgFlex:          pgtype.Float8{Float64: r.AvgFlex, Valid: true},
				SensorDataPoints: r.SensorData,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
