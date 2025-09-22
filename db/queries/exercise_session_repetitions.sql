-- name: CreateRepetition :one
INSERT INTO exercise_session_repetitions (
  exercise_session_id, rep_index,
  avg_range_of_motion, peak_range_of_motion,
  avg_angle_velocity,  peak_angle_velocity
) VALUES ($1,$2,$3,$4,$5,$6)
ON CONFLICT (exercise_session_id, rep_index) DO UPDATE
  SET avg_range_of_motion  = EXCLUDED.avg_range_of_motion,
      peak_range_of_motion = EXCLUDED.peak_range_of_motion,
      avg_angle_velocity   = EXCLUDED.avg_angle_velocity,
      peak_angle_velocity  = EXCLUDED.peak_angle_velocity
RETURNING *;

-- name: ListRepetitionsForSession :many
SELECT *
FROM exercise_session_repetitions
WHERE exercise_session_id = $1
ORDER BY rep_index NULLS LAST;