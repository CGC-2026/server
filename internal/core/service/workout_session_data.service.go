package service

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/data"
	"server/internal/db"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CoachingScoreDto struct {
	OverallQuality  string  `json:"overallQuality"`
	GoodReps        int32   `json:"goodReps"`
	OkayReps        int32   `json:"okayReps"`
	BadReps         int32   `json:"badReps"`
	AverageDepth    float64 `json:"averageDepth"`
	AverageDuration float64 `json:"averageDuration"`
}

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
		if dto.SetNumber <= 0 {
			return fmt.Errorf("setNumber must be greater than 0")
		}
		if len(dto.Reps) == 0 {
			return fmt.Errorf("at least one rep is required")
		}

		// ensure unique, sequential reps and valid time ranges
		seenRep := make(map[int32]struct{}, len(dto.Reps))
		for _, r := range dto.Reps {
			if r.RepNumber <= 0 {
				return fmt.Errorf("repNumber must be positive")
			}
			if _, exists := seenRep[r.RepNumber]; exists {
				return fmt.Errorf("duplicate repNumber %d in set", r.RepNumber)
			}
			seenRep[r.RepNumber] = struct{}{}

			if r.StartTime != nil && r.EndTime != nil && r.EndTime.Before(*r.StartTime) {
				return fmt.Errorf("rep %d has endTime before startTime", r.RepNumber)
			}

			if len(r.SensorData) == 0 {
				return fmt.Errorf("rep %d has empty sensorData", r.RepNumber)
			}

			// validate sensorData JSON and enforce a soft size limit of 1MB
			if len(r.SensorData) > 1*1024*1024 {
				return fmt.Errorf("rep %d sensorData too large", r.RepNumber)
			}
			var tmp any
			if err := json.Unmarshal(r.SensorData, &tmp); err != nil {
				return fmt.Errorf("rep %d has invalid sensorData JSON: %w", r.RepNumber, err)
			}
		}

		var coachingScoreBytes []byte
		if dto.CoachingScore != nil {
			b, err := json.Marshal(dto.CoachingScore)
			if err != nil {
				return err
			}
			coachingScoreBytes = b
		}

		startTime := time.Now()
		if dto.StartTime != nil && !dto.StartTime.IsZero() {
			startTime = *dto.StartTime
		}

		set, err := qtx.CreateWorkoutSet(ctx, db.CreateWorkoutSetParams{
			WorkoutSessionID: dto.WorkoutSessionID,
			SetNumber:        dto.SetNumber,
			StartTime:        pgtype.Timestamptz{Time: startTime, Valid: true},
			EndTime: pgtype.Timestamptz{
				Valid: dto.EndTime != nil && !dto.EndTime.IsZero(),
				Time: func() time.Time {
					if dto.EndTime != nil {
						return *dto.EndTime
					}
					return time.Time{}
				}(),
			},
			CoachingScore: coachingScoreBytes,
		})
		if err != nil {
			return err
		}

		for _, r := range dto.Reps {
			err := qtx.CreateWorkoutRep(ctx, db.CreateWorkoutRepParams{
				WorkoutSetID: set.ID,
				RepNumber:    r.RepNumber,
				StartTime: pgtype.Timestamptz{
					Valid: r.StartTime != nil,
					Time: func() time.Time {
						if r.StartTime != nil {
							return *r.StartTime
						}
						return time.Time{}
					}(),
				},
				EndTime: pgtype.Timestamptz{
					Valid: r.EndTime != nil,
					Time: func() time.Time {
						if r.EndTime != nil {
							return *r.EndTime
						}
						return time.Time{}
					}(),
				},
				Quality:          pgtype.Text{String: r.Quality, Valid: true},
				PeakAngle:        pgtype.Float8{Float64: r.PeakAngle, Valid: true},
				AvgFlex:          pgtype.Float8{Float64: r.AvgFlex, Valid: true},
				SensorDataPoints: []byte(r.SensorData),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
