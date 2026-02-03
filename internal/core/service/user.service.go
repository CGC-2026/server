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

type ClerkUserData struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	EmailAddresses []struct {
		EmailAddress string `json:"email_address"`
		ID           string `json:"id"`
	} `json:"email_addresses"`
	ImageURL string `json:"image_url"`
}

type CalibrationDataDto struct {
	UserId        string     `json:"userId"`
	StandingAngle float64    `json:"standingAngle"`
	StandingFlex  float64    `json:"standingFlex"`
	LastUpdated   *time.Time `json:"lastUpdated"`
}

func HandleCreateUserFromClerk(ctx context.Context, store *data.Store, data json.RawMessage) error {
	var userData ClerkUserData
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
	_, err := store.Queries.CreateUser(ctx, db.CreateUserParams{
		ID:        userData.ID,
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

func GetUserCalibrationData(ctx context.Context, store *data.Store, userID string) (CalibrationDataDto, error) {
	cal, err := store.Queries.GetUserCalibrationByUserID(ctx, userID)
	if err != nil {
		return CalibrationDataDto{}, err
	}

	return CalibrationDataDto{
		UserId:        cal.UserID,
		StandingAngle: cal.StandingAngle,
		StandingFlex:  cal.StandingFlex,
		LastUpdated:   &cal.UpdatedAt.Time,
	}, nil
}

func SaveUserCalibration(ctx context.Context, store *data.Store, dto CalibrationDataDto) error {
	_, err := store.Queries.UpsertUserCalibration(ctx, db.UpsertUserCalibrationParams{
		UserID:        dto.UserId,
		StandingAngle: dto.StandingAngle,
		StandingFlex:  dto.StandingFlex,
	})
	return err
}

func DeleteUserCalibration(ctx context.Context, store *data.Store, userID string) error {
	return store.Queries.DeleteUserCalibration(ctx, userID)
}
