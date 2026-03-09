package service

import (
	"context"
	"fmt"
	"server/internal/core/service/helpers"
	"server/internal/data"
	"server/internal/db"
	"time"

	sqlc "server/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateWorkoutSessionDTO struct {
	ID                    string     `json:"id"`
	UserID                string     `json:"userId"`
	WorkoutTypeID         string     `json:"workoutTypeId"`
	CalibrationYawAngle   *float64   `json:"calibrationYawAngle"`
	CalibrationPitchAngle *float64   `json:"calibrationPitchAngle"`
	CalibrationRollAngle  *float64   `json:"calibrationRollAngle"`
	StartTime             *time.Time `json:"startTime"`
}

type UpdateWorkoutSessionDTO struct {
	ID                    string     `json:"id"`
	UserID                string     `json:"userId"`
	WorkoutTypeID         string     `json:"workoutTypeId"`
	CalibrationYawAngle   *float64   `json:"calibrationYawAngle"`
	CalibrationPitchAngle *float64   `json:"calibrationPitchAngle"`
	CalibrationRollAngle  *float64   `json:"calibrationRollAngle"`
	StartTime             *time.Time `json:"startTime"`
	EndTime               *time.Time `json:"endTime"`
}

type WorkoutSessionDTO struct {
	ID                    string     `json:"id"`
	UserID                string     `json:"userId"`
	WorkoutTypeID         string     `json:"workoutTypeId"`
	CalibrationYawAngle   *float64   `json:"calibrationYawAngle"`
	CalibrationPitchAngle *float64   `json:"calibrationPitchAngle"`
	CalibrationRollAngle  *float64   `json:"calibrationRollAngle"`
	StartTime             *time.Time `json:"startTime"`
	EndTime               *time.Time `json:"endTime"`
	CreatedAt             time.Time  `json:"createdAt"`
	UpdatedAt             time.Time  `json:"updatedAt"`
}

func CreateWorkoutSession(ctx context.Context, store *data.Store, userID string, dto CreateWorkoutSessionDTO) (sqlc.WorkoutSession, error) {

	startTime := time.Now()
	if dto.StartTime != nil && !dto.StartTime.IsZero() {
		startTime = *dto.StartTime
	}

	return store.Queries.CreateWorkoutSession(ctx, sqlc.CreateWorkoutSessionParams{
		UserID:                userID,
		WorkoutTypeID:         dto.WorkoutTypeID,
		CalibrationYawAngle:   helpers.NewNullFloat(dto.CalibrationYawAngle),
		CalibrationPitchAngle: helpers.NewNullFloat(dto.CalibrationPitchAngle),
		CalibrationRollAngle:  helpers.NewNullFloat(dto.CalibrationRollAngle),
		StartTime:             pgtype.Timestamptz{Time: startTime, Valid: true},
		EndTime:               pgtype.Timestamptz{Valid: false},
	})
}

func GetWorkoutSessionByID(ctx context.Context, store *data.Store, id string) (sqlc.WorkoutSession, error) {
	return store.Queries.GetWorkoutSessionByID(ctx, id)
}

func UpdateWorkoutSession(ctx context.Context, store *data.Store, userID string, sessionID string, dto UpdateWorkoutSessionDTO) (sqlc.WorkoutSession, error) {
	existingUser, err := store.Queries.GetWorkoutSessionByID(ctx, sessionID)
	if err != nil {
		return db.WorkoutSession{}, err
	}

	if existingUser.UserID != userID {
		return db.WorkoutSession{}, fmt.Errorf("forbidden")
	}

	params := db.UpdateWorkoutSessionParams{
		ID:                    sessionID,
		WorkoutTypeID:         helpers.NewNullString(dto.WorkoutTypeID),
		CalibrationYawAngle:   helpers.NewNullFloat(dto.CalibrationYawAngle),
		CalibrationPitchAngle: helpers.NewNullFloat(dto.CalibrationPitchAngle),
		CalibrationRollAngle:  helpers.NewNullFloat(dto.CalibrationRollAngle),
		StartTime:             helpers.NewNullTime(dto.StartTime),
		EndTime:               helpers.NewNullTime(dto.EndTime),
	}

	return store.Queries.UpdateWorkoutSession(ctx, params)

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
