-- name: CreateExerciseSession :one
INSERT INTO user_exercise_sessions (
  user_id, exercise_type_id, start_time_utc, end_time_utc
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: FinishExerciseSession :one
UPDATE user_exercise_sessions
SET end_time_utc = $2
WHERE id = $1
RETURNING *;

-- name: ListExerciseSessionsByUser :many
SELECT *
FROM user_exercise_sessions
WHERE user_id = $1
ORDER BY start_time_utc DESC
LIMIT $2 OFFSET $3;

-- name: GetExerciseSession :one
SELECT * FROM user_exercise_sessions WHERE id = $1;