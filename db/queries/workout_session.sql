
-- name: CreateWorkoutSession :one
INSERT INTO workout_session (
    user_id,
    workout_type_id,
    calibration_yaw_angle,
    calibration_pitch_angle,
    calibration_roll_angle,
    start_time,
    end_time)
VALUES ($1, $2, $3, $4, $5, $6, $7)
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

-- name: GetUserWorkoutSessionHistory :many
SELECT
    ws.id,
    ws.user_id,
    ws.workout_type_id,
    ws.start_time,
    ws.end_time,
    wt.name as workout_type_name,
    COUNT(DISTINCT s.id)::int AS total_sets,
    COUNT(r.id)::int AS total_reps
FROM workout_session ws
JOIN workout_type wt 
    ON ws.workout_type_id = wt.id
LEFT JOIN workout_session_set s
    ON ws.id = s.workout_session_id
LEFT JOIN workout_session_set_rep r
    ON s.id = r.workout_session_set_id
WHERE ws.user_id = $1
GROUP BY
    ws.id,
    ws.user_id,
    ws.workout_type_id,
    wt.name,
    ws.start_time,
    ws.end_time
ORDER BY ws.start_time DESC
LIMIT $2;