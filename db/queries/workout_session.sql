
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
    calibration_yaw_angle = COALESCE(sqlc.narg('calibration_yaw_angle'), calibration_yaw_angle),
    calibration_pitch_angle = COALESCE(sqlc.narg('calibration_pitch_angle'), calibration_pitch_angle),
    calibration_roll_angle = COALESCE(sqlc.narg('calibration_roll_angle'), calibration_roll_angle),
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
    ws.calibration_yaw_angle,
    ws.calibration_pitch_angle,
    ws.calibration_roll_angle,
    ws.start_time,
    ws.end_time,
    wt.name as workout_type_name
FROM workout_session ws
JOIN workout_type wt ON ws.workout_type_id = wt.id
WHERE ws.user_id = $1
ORDER BY ws.start_time DESC
LIMIT $2;