
-- name: CreateWorkoutRep :exec
INSERT INTO workout_rep (workout_set_id, rep_number, start_time, end_time, quality, peak_angle, avg_flex, sensor_data_points)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetRepsByWorkoutSetId :many
SELECT * FROM workout_rep WHERE workout_set_id = $1 ORDER BY rep_number ASC;