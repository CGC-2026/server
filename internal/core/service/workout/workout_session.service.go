package workout

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/internal/core/service/helpers"
	"server/internal/data"
	"server/internal/db"
	"time"

	sqlc "server/internal/db"

	"github.com/jackc/pgx/v5"
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

type WorkoutSessionSetRepDTO struct {
	ID        string `json:"id"`
	RepNumber int32  `json:"rep_number"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Samples   any    `json:"samples"`
	Metrics   any    `json:"metrics"`
}

type WorkoutSessionSetDTO struct {
	ID               string                    `json:"id"`
	WorkoutSessionID string                    `json:"workout_session_id"`
	SetNumber        int32                     `json:"set_number"`
	StartTime        time.Time                 `json:"start_time"`
	EndTime          *time.Time                `json:"end_time"`
	Reps             []WorkoutSessionSetRepDTO `json:"reps"`
}

type WorkoutSessionDTO struct {
	ID                    string                 `json:"id"`
	UserID                string                 `json:"user_id"`
	WorkoutTypeID         string                 `json:"workout_type_id"`
	CalibrationYawAngle   *float64               `json:"calibration_yaw_angle"`
	CalibrationPitchAngle *float64               `json:"calibration_pitch_angle"`
	CalibrationRollAngle  *float64               `json:"calibration_roll_angle"`
	Sets                  []WorkoutSessionSetDTO `json:"sets"`
	StartTime             time.Time              `json:"start_time"`
	EndTime               *time.Time             `json:"end_time"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
}

type WorkoutSessionHistoryItemDTO struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	WorkoutTypeID   string     `json:"workout_type_id"`
	WorkoutTypeName string     `json:"workout_type_name"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	TotalSets       int32      `json:"total_sets"`
	TotalReps       int32      `json:"total_reps"`
}

func CreateWorkoutSession(ctx context.Context, store *data.Store, userID string, dto CreateWorkoutSessionDTO) (sqlc.WorkoutSession, error) {

	startTime := time.Now()
	if dto.StartTime != nil && !dto.StartTime.IsZero() {
		startTime = *dto.StartTime
	}

	calibrations, err := store.Queries.GetUserCalibrationByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			calibrations = sqlc.UserCalibration{
				StandingYawAngle:   0,
				StandingPitchAngle: 0,
				StandingRollAngle:  0,
			}
		} else {
			return db.WorkoutSession{}, err
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

func GetWorkoutSessionByID(ctx context.Context, store *data.Store, userID string, sessionID string) (WorkoutSessionDTO, error) {
	session, err := store.Queries.GetWorkoutSessionByID(ctx, sessionID)
	if err != nil {
		return WorkoutSessionDTO{}, err
	}

	if session.UserID != userID {
		return WorkoutSessionDTO{}, fmt.Errorf("forbidden")
	}

	setRows, err := store.Queries.GetWorkoutSessionSetsBySessionID(ctx, sessionID)
	if err != nil {
		return WorkoutSessionDTO{}, err
	}

	repRows, err := store.Queries.GetWorkoutSessionSetRepsBySessionID(ctx, sessionID)
	if err != nil {
		return WorkoutSessionDTO{}, err
	}

	repsBySetID := make(map[string][]db.WorkoutSessionSetRep, len(setRows))
	for _, repRow := range repRows {
		repsBySetID[repRow.WorkoutSessionSetID] = append(repsBySetID[repRow.WorkoutSessionSetID], repRow)
	}

	sets := make([]WorkoutSessionSetDTO, 0, len(setRows))

	for _, setRow := range setRows {
		reps := make([]WorkoutSessionSetRepDTO, 0, len(repsBySetID[setRow.ID]))

		for _, repRow := range repsBySetID[setRow.ID] {
			var samples any
			var metrics any

			if err := json.Unmarshal(repRow.Samples, &samples); err != nil {
				return WorkoutSessionDTO{}, err
			}
			if err := json.Unmarshal(repRow.Metrics, &metrics); err != nil {
				return WorkoutSessionDTO{}, err
			}

			reps = append(reps, WorkoutSessionSetRepDTO{
				ID:        repRow.ID,
				RepNumber: repRow.RepNumber,
				StartTime: repRow.StartTime,
				EndTime:   repRow.EndTime,
				Samples:   samples,
				Metrics:   metrics,
			})
		}

		var setEndTime *time.Time
		if setRow.EndTime.Valid {
			t := setRow.EndTime.Time
			setEndTime = &t
		}

		sets = append(sets, WorkoutSessionSetDTO{
			ID:               setRow.ID,
			WorkoutSessionID: setRow.WorkoutSessionID,
			SetNumber:        setRow.SetNumber,
			StartTime:        setRow.StartTime.Time,
			EndTime:          setEndTime,
			Reps:             reps,
		})
	}

	var sessionEndTime *time.Time
	if session.EndTime.Valid {
		t := session.EndTime.Time
		sessionEndTime = &t
	}

	var yaw *float64
	if session.CalibrationYawAngle.Valid {
		v := session.CalibrationYawAngle.Float64
		yaw = &v
	}

	var pitch *float64
	if session.CalibrationPitchAngle.Valid {
		v := session.CalibrationPitchAngle.Float64
		pitch = &v
	}

	var roll *float64
	if session.CalibrationRollAngle.Valid {
		v := session.CalibrationRollAngle.Float64
		roll = &v
	}

	return WorkoutSessionDTO{
		ID:                    session.ID,
		UserID:                session.UserID,
		WorkoutTypeID:         session.WorkoutTypeID,
		CalibrationYawAngle:   yaw,
		CalibrationPitchAngle: pitch,
		CalibrationRollAngle:  roll,
		Sets:                  sets,
		StartTime:             session.StartTime.Time,
		EndTime:               sessionEndTime,
		CreatedAt:             session.CreatedAt.Time,
		UpdatedAt:             session.UpdatedAt.Time,
	}, nil
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

func GetUserWorkoutSessionHistory(ctx context.Context, store *data.Store, userId string, limit int32) ([]WorkoutSessionHistoryItemDTO, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := store.Queries.GetUserWorkoutSessionHistory(ctx, db.GetUserWorkoutSessionHistoryParams{
		UserID: userId,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	items := make([]WorkoutSessionHistoryItemDTO, 0, len(rows))
	for _, row := range rows {
		var endTime *time.Time
		if row.EndTime.Valid {
			t := row.EndTime.Time
			endTime = &t
		}

		items = append(items, WorkoutSessionHistoryItemDTO{
			ID:              row.ID,
			UserID:          row.UserID,
			WorkoutTypeID:   row.WorkoutTypeID,
			WorkoutTypeName: row.WorkoutTypeName,
			StartTime:       row.StartTime.Time,
			EndTime:         endTime,
			TotalSets:       row.TotalSets,
			TotalReps:       row.TotalReps,
		})
	}

	return items, nil
}
