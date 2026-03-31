//go:build !integration

package workout

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateCreateWorkoutSessionSetDTO(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		dto         CreateWorkoutSessionSetDTO
		expectedErr string
	}{
		{
			name: "invalid set number",
			dto: CreateWorkoutSessionSetDTO{
				SetNumber: 0,
			},
			expectedErr: "set_number must be greater than 0",
		},
		{
			name: "invalid rep number",
			dto: CreateWorkoutSessionSetDTO{
				SetNumber: 1,
				Reps: []CreateWorkoutSessionSetRepDTO{
					validRep(func(rep *CreateWorkoutSessionSetRepDTO) {
						rep.RepNumber = 0
					}),
				},
			},
			expectedErr: "rep_number must be greater than 0",
		},
		{
			name: "rep end before start",
			dto: CreateWorkoutSessionSetDTO{
				SetNumber: 1,
				Reps: []CreateWorkoutSessionSetRepDTO{
					validRep(func(rep *CreateWorkoutSessionSetRepDTO) {
						rep.EndTime = rep.StartTime - 1
					}),
				},
			},
			expectedErr: "rep end_time must be after start_time",
		},
		{
			name: "invalid samples json",
			dto: CreateWorkoutSessionSetDTO{
				SetNumber: 1,
				Reps: []CreateWorkoutSessionSetRepDTO{
					validRep(func(rep *CreateWorkoutSessionSetRepDTO) {
						rep.Samples = json.RawMessage(`{`)
					}),
				},
			},
			expectedErr: "samples must be valid JSON",
		},
		{
			name: "invalid metrics json",
			dto: CreateWorkoutSessionSetDTO{
				SetNumber: 1,
				Reps: []CreateWorkoutSessionSetRepDTO{
					validRep(func(rep *CreateWorkoutSessionSetRepDTO) {
						rep.Metrics = json.RawMessage(`{`)
					}),
				},
			},
			expectedErr: "metrics must be valid JSON",
		},
		{
			name: "valid dto",
			dto: CreateWorkoutSessionSetDTO{
				SetNumber: 1,
				Reps: []CreateWorkoutSessionSetRepDTO{
					validRep(nil),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validateCreateWorkoutSessionSetDTO(tc.dto)
			if tc.expectedErr == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tc.expectedErr)
		})
	}
}

func validRep(mutate func(*CreateWorkoutSessionSetRepDTO)) CreateWorkoutSessionSetRepDTO {
	rep := CreateWorkoutSessionSetRepDTO{
		RepNumber: 1,
		StartTime: 100,
		EndTime:   200,
		Samples:   json.RawMessage(`[{"angle": 90}]`),
		Metrics:   json.RawMessage(`{"tempo":"controlled"}`),
	}

	if mutate != nil {
		mutate(&rep)
	}

	return rep
}
