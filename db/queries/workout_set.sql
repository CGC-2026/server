
-- name: CreateWorkoutSet :one
INSERT INTO workout_set (workout_session_id, set_number, start_time, end_time, coaching_score)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSetsByWorkoutSessionId :many
SELECT * FROM workout_set WHERE workout_session_id = $1 ORDER BY set_number ASC;