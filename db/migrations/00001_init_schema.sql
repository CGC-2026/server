-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id text PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email text NOT NULL
);

CREATE TABLE IF NOT EXISTS exercise_types (
    id text PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE IF NOT EXISTS user_exercise_sessions (
    id text PRIMARY KEY,
    user_id text NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_type_id text NOT NULL REFERENCES exercise_types(id) ON DELETE CASCADE,
    start_time_utc timestamptz NOT NULL,
    end_time_utc timestamptz NULL
);

CREATE TABLE IF NOT EXISTS exercise_session_repetitions (
    id text PRIMARY KEY,
    exercise_session_id text NOT NULL REFERENCES user_exercise_sessions(id) ON DELETE CASCADE,
    rep_index integer NULL,
    avg_range_of_motion real NULL,
    peak_range_of_motion real NULL,
    avg_angle_velocity real NULL,
    peak_angle_velocity real NULL
);

CREATE TABLE IF NOT EXISTS exercise_session_joint_samples (
    id text PRIMARY KEY,
    exercise_session_id text NOT NULL,
    pitch_angle real NULL,
    yaw_angle real NULL,
    roll_angle real NULL,
    timestamp_ms integer NULL,
    ax_mps2 real NULL,
    ay_mps2 real NULL,
    az_mps2 real NULL,
    gx_dps real NULL,
    gy_dps real NULL,
    gz_dps real NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exercise_session_joint_samples;
DROP TABLE IF EXISTS exercise_session_repetitions;
DROP TABLE IF EXISTS user_exercise_sessions;
DROP TABLE IF EXISTS exercise_types;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS pgcrypto;
-- +goose StatementEnd
