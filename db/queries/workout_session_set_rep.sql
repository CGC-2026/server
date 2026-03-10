-- name: CreateWorkoutSessionSetRep :one
INSERT INTO workout_session_set_rep (
    workout_session_set_id,
    rep_number,
    start_time,
    end_time,
    samples,
    metrics
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetWorkoutSessionSetRepsBySessionSetID :many
SELECT *
FROM workout_session_set_rep
WHERE workout_session_set_id = $1
ORDER BY rep_number ASC;