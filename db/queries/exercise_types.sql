-- name: UpsertExerciseTypeByName :one
INSERT INTO exercise_types (name)
VALUES ($1)
ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
RETURNING *;

-- name: ListExerciseTypes :many
SELECT * FROM exercise_types ORDER BY name;