package service

import (
	"context"
	"server/internal/data"
	"server/internal/db"
	"time"

	sqlc "server/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateWorkoutSessionDto struct {
	ID            string     `json:"id"`
	UserId        string     `json:"userId"`
	WorkoutTypeId string     `json:"workoutTypeId"`
	StartTime     *time.Time `json:"startTime"`
}

type UpdateWorkoutSessionDto struct {
	ID            string     `json:"id"`
	UserId        string     `json:"userId"`
	WorkoutTypeId string     `json:"workoutTypeId"`
	StartTime     *time.Time `json:"startTime"`
	EndTime       *time.Time `json:"endTime"`
}

type WorkoutSessionDto struct {
	ID            string     `json:"id"`
	UserId        string     `json:"userId"`
	WorkoutTypeId string     `json:"workoutTypeId"`
	StartTime     *time.Time `json:"startTime"`
	EndTime       *time.Time `json:"endTime"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func CreateWorkoutSession(ctx context.Context, store *data.Store, dto CreateWorkoutSessionDto) (sqlc.WorkoutSession, error) {

	startTime := time.Now()
	if dto.StartTime != nil && !dto.StartTime.IsZero() {
		startTime = *dto.StartTime
	}

	return store.Queries.CreateWorkoutSession(ctx, sqlc.CreateWorkoutSessionParams{
		UserID:        dto.UserId,
		WorkoutTypeID: dto.WorkoutTypeId,
		StartTime:     pgtype.Timestamptz{Time: startTime, Valid: true},
	})
}

func GetWorkoutSessionByID(ctx context.Context, store *data.Store, id string) (sqlc.WorkoutSession, error) {
	return store.Queries.GetWorkoutSessionByID(ctx, id)
}

func UpdateWorkoutSession(ctx context.Context, store *data.Store, id string, dto UpdateWorkoutSessionDto) (sqlc.WorkoutSession, error) {
	params := db.UpdateWorkoutSessionParams{
		ID: id,
		WorkoutTypeID: pgtype.Text{
			Valid: dto.WorkoutTypeId != "",
			String: func() string {
				if dto.WorkoutTypeId != "" {
					return dto.WorkoutTypeId
				}
				return ""
			}(),
		},
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
	}

	return store.Queries.UpdateWorkoutSession(ctx, params)

}

func GetWorkoutTypes(ctx context.Context, store *data.Store) ([]sqlc.WorkoutType, error) {
	return store.Queries.GetWorkoutTypeList(ctx)
}

func GetUserWorkoutSessionHistory(ctx context.Context, store *data.Store, userId string, limit int32) ([]db.GetUserWorkoutSessionHistoryRow, error) {
	if limit <= 0 {
		limit = 20
	}

	return store.Queries.GetUserWorkoutSessionHistory(ctx, db.GetUserWorkoutSessionHistoryParams{
		UserID: userId,
		Limit:  limit,
	})
}
