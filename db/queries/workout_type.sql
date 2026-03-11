
-- name: CreateWorkoutType :one
INSERT INTO workout_type(name, description, config, image_url)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: GetWorkoutType :one
SELECT * FROM workout_type WHERE id = $1;

-- name: GetWorkoutTypeList :many
SELECT * FROM workout_type ORDER BY name ASC;