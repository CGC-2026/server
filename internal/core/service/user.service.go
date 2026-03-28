package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/internal/data"
	"time"

	sqlc "server/internal/db"

	"server/internal/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserProfileDTO struct {
	ID        string `json:"id"`
	ClerkID   string `json:"clerk_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	ImageURL  string `json:"image_url,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserCalibrationDTO struct {
	UserID             string  `json:"user_id"`
	StandingYawAngle   float64 `json:"standing_yaw_angle"`
	StandingPitchAngle float64 `json:"standing_pitch_angle"`
	StandingRollAngle  float64 `json:"standing_roll_angle"`
	UpdatedAt          string  `json:"updated_at"`
}

type UpsertUserCalibrationDTO struct {
	StandingYawAngle   float64 `json:"standing_yaw_angle"`
	StandingPitchAngle float64 `json:"standing_pitch_angle"`
	StandingRollAngle  float64 `json:"standing_roll_angle"`
}

func dtoFromCalibrationRow(row sqlc.UserCalibration) UserCalibrationDTO {
	updated := ""
	if row.UpdatedAt.Valid {
		updated = row.UpdatedAt.Time.Format(time.RFC3339)
	}

	return UserCalibrationDTO{
		UserID:             row.UserID,
		StandingYawAngle:   row.StandingYawAngle,
		StandingPitchAngle: row.StandingPitchAngle,
		StandingRollAngle:  row.StandingRollAngle,
		UpdatedAt:          updated,
	}
}

func dtoFromUserRow(row sqlc.User) UserProfileDTO {
	created := ""
	if row.UpdatedAt.Valid {
		created = row.UpdatedAt.Time.Format(time.RFC3339)
	}

	updated := ""
	if row.UpdatedAt.Valid {
		updated = row.UpdatedAt.Time.Format(time.RFC3339)
	}

	dto := UserProfileDTO{
		ID:        row.ID,
		ClerkID:   row.ClerkID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Email:     row.Email,
		CreatedAt: created,
		UpdatedAt: updated,
	}

	if row.ImageUrl.Valid {
		dto.ImageURL = row.ImageUrl.String
	}

	return dto
}

func GetCurrentUserProfile(ctx context.Context, store *data.Store, userID string) (UserProfileDTO, error) {
	user, err := store.Queries.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserProfileDTO{}, fmt.Errorf("user not found: %w", err)
		}
		return UserProfileDTO{}, fmt.Errorf("error fetching user: %w", err)
	}

	return dtoFromUserRow(user), nil
}

func GetCurrentUserCalibration(ctx context.Context, store *data.Store, userID string) (UserCalibrationDTO, error) {
	calibration, err := store.Queries.GetUserCalibrationByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserCalibrationDTO{}, fmt.Errorf("calibration not found: %w", err)
		}
		return UserCalibrationDTO{}, fmt.Errorf("error fetching calibration: %w", err)
	}

	return dtoFromCalibrationRow(calibration), nil
}

func UpsertUserCalibration(ctx context.Context, store *data.Store, userID string, dto UpsertUserCalibrationDTO) (UserCalibrationDTO, error) {
	calibration, err := store.Queries.UpsertUserCalibration(ctx, sqlc.UpsertUserCalibrationParams{
		UserID:             userID,
		StandingYawAngle:   dto.StandingYawAngle,
		StandingPitchAngle: dto.StandingPitchAngle,
		StandingRollAngle:  dto.StandingRollAngle,
	})
	if err != nil {
		return UserCalibrationDTO{}, fmt.Errorf("error upserting calibration: %w", err)
	}

	return dtoFromCalibrationRow(calibration), nil
}

func DeleteUserCalibration(ctx context.Context, store *data.Store, userID string) error {
	if err := store.Queries.DeleteUserCalibrationByUserID(ctx, userID); err != nil {
		return fmt.Errorf("error deleting calibration: %w", err)
	}
	return nil
}

func HandleCreateUserFromClerk(ctx context.Context, store *data.Store, data json.RawMessage) error {
	var userData types.ClerkUserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return fmt.Errorf("failed to parse user data: %w", err)
	}

	// Extract primary email
	email := ""
	if len(userData.EmailAddresses) > 0 {
		email = userData.EmailAddresses[0].EmailAddress
	}

	// Create image URL as nullable text
	var imageURL pgtype.Text
	if userData.ImageURL != "" {
		imageURL = pgtype.Text{String: userData.ImageURL, Valid: true}
	}

	// Create user in database
	_, err := store.Queries.CreateUser(ctx, sqlc.CreateUserParams{
		ClerkID:   userData.ID,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     email,
		ImageUrl:  imageURL,
	})

	if err != nil {
		return fmt.Errorf("failed to create user in database: %w", err)
	}

	return nil
}

func HandleDeleteUserFromClerk(ctx context.Context, store *data.Store, data json.RawMessage) error {
	var userData types.ClerkUserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return fmt.Errorf("failed to parse user data: %w", err)
	}

	return store.Queries.DeleteUserByClerkID(ctx, userData.ID)
}
