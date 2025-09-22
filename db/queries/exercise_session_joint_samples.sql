-- name: InsertJointSample :exec
INSERT INTO exercise_session_joint_samples (
  exercise_session_id, timestamp_ms,
  pitch_angle, yaw_angle, roll_angle,
  ax_mps2, ay_mps2, az_mps2,
  gx_dps, gy_dps, gz_dps
) VALUES (
  $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11
);

-- name: GetSamplesInRange :many
SELECT *
FROM exercise_session_joint_samples
WHERE exercise_session_id = $1
  AND timestamp_ms BETWEEN $2 AND $3
ORDER BY timestamp_ms;

-- name: DeleteSamplesForSession :exec
DELETE FROM exercise_session_joint_samples WHERE exercise_session_id = $1;