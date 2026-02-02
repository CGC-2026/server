
-- name: CreateWorkoutSession :one
INSERT INTO workout_session (user_id, workout_type_id, start_time, end_time)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWorkoutSessionByID :one
SELECT * FROM workout_session WHERE id = $1;

-- name: UpdateWorkoutSession :one
UPDATE workout_session
SET workout_type_id = COALESCE(sqlc.narg('workout_type_id'), workout_type_id),
    start_time = COALESCE(sqlc.narg('start_time'), start_time),
    end_time = COALESCE(sqlc.narg('end_time'), end_time),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: EndWorkoutSession :exec
UPDATE workout_session
SET end_time = $2, updated_at = NOW()
WHERE id = $1;

-- name: GetUserWorkoutSessionHistory :many
SELECT
    ws.id,
    ws.start_time,
    ws.end_time,
    wt.name as workout_type_name
FROM workout_session ws
JOIN workout_type wt ON ws.workout_type_id = wt.id
WHERE ws.user_id = $1
ORDER BY ws.start_time DESC
LIMIT $2;