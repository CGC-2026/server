package service

import (
	"context"
	"encoding/json"
	"server/internal/data"
)

type WorkoutTypeDTO struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"description"`
	Config      json.RawMessage `json:"config"`
	ImageURL    *string         `json:"imageUrl"`
}

func GetWorkoutTypes(ctx context.Context, store *data.Store) ([]WorkoutTypeDTO, error) {
	rows, err := store.Queries.GetWorkoutTypeList(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]WorkoutTypeDTO, 0, len(rows))
	for _, row := range rows {
		dto := WorkoutTypeDTO{
			ID:     row.ID,
			Name:   row.Name,
			Config: json.RawMessage(row.Config),
		}

		if row.Description.Valid {
			desc := row.Description.String
			dto.Description = &desc
		}

		if row.ImageUrl.Valid {
			imageUrl := row.ImageUrl.String
			dto.ImageURL = &imageUrl
		}

		out = append(out, dto)
	}

	return out, nil
}
