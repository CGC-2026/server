-- name: GetUserCalibrationByUserID :one
SELECT * FROM user_calibration WHERE user_id = $1;

-- name: UpsertUserCalibration :one
INSERT INTO user_calibration (user_id, standing_angle, standing_flex, updated_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (user_id)
DO UPDATE SET
    standing_angle = EXCLUDED.standing_angle,
    standing_flex = EXCLUDED.standing_flex,
    updated_at = NOW()
RETURNING *;

-- name: DeleteUserCalibration :exec
DELETE FROM user_calibration WHERE user_id = $1;