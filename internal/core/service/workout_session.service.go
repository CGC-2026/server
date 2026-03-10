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
	WorkoutTypeID string     `json:"workout_type_id"`
	StartTime     *time.Time `json:"start_time"`
}

type UpdateWorkoutSessionDTO struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	WorkoutTypeID string     `json:"workout_type_id"`
	StartTime     *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
}

type WorkoutSessionDTO struct {
	ID                    string     `json:"id"`
	UserID                string     `json:"user_id"`
	WorkoutTypeID         string     `json:"workout_type_id"`
	CalibrationYawAngle   *float64   `json:"calibration_yaw_angle"`
	CalibrationPitchAngle *float64   `json:"calibration_pitch_angle"`
	CalibrationRollAngle  *float64   `json:"calibration_roll_angle"`
	StartTime             *time.Time `json:"start_time"`
	EndTime               *time.Time `json:"end_time"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func CreateWorkoutSession(ctx context.Context, store *data.Store, userID string, dto CreateWorkoutSessionDTO) (sqlc.WorkoutSession, error) {

	startTime := time.Now()
	if dto.StartTime != nil && !dto.StartTime.IsZero() {
		startTime = *dto.StartTime
	}

	calibrations, err := store.Queries.GetUserCalibrationByUserID(ctx, userID)
	if err != nil {
		calibrations = sqlc.UserCalibration{
			StandingYawAngle:   0,
			StandingPitchAngle: 0,
			StandingRollAngle:  0,
		}
	}

	return store.Queries.CreateWorkoutSession(ctx, sqlc.CreateWorkoutSessionParams{
		UserID:                userID,
		WorkoutTypeID:         dto.WorkoutTypeID,
		CalibrationYawAngle:   helpers.NewNullFloat(&calibrations.StandingYawAngle),
		CalibrationPitchAngle: helpers.NewNullFloat(&calibrations.StandingPitchAngle),
		CalibrationRollAngle:  helpers.NewNullFloat(&calibrations.StandingRollAngle),
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
		ID:            sessionID,
		WorkoutTypeID: helpers.NewNullString(dto.WorkoutTypeID),
		StartTime:     helpers.NewNullTime(dto.StartTime),
		EndTime:       helpers.NewNullTime(dto.EndTime),
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
