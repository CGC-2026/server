-- name: GetUserCalibrationByUserID :one
SELECT * FROM user_calibration WHERE user_id = $1;

-- name: UpsertUserCalibration :one
INSERT INTO user_calibration (
    user_id,
    standing_yaw_angle,
    standing_pitch_angle,
    standing_roll_angle,
    updated_at
)
VALUES ($1, $2, $3, $4, NOW())
-- If a calibration already exists for the user, update it with the new values and set updated_at to the current time
ON CONFLICT (user_id) DO UPDATE SET
    standing_yaw_angle = EXCLUDED.standing_yaw_angle,
    standing_pitch_angle = EXCLUDED.standing_pitch_angle,
    standing_roll_angle = EXCLUDED.standing_roll_angle,
    updated_at = NOW()
RETURNING *;

-- name: DeleteUserCalibrationByUserID :exec
DELETE FROM user_calibration WHERE user_id = $1;