

-- name: CreateWorkoutSessionSet :one
INSERT INTO workout_session_set (
    workout_session_id,
    set_number,
    start_time,
    end_time
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWorkoutSessionSetsBySessionID :many
SELECT *
FROM workout_session_set
WHERE workout_session_id = $1
ORDER BY set_number ASC;