-- +goose Up

CREATE TABLE IF NOT EXISTS workout_session_set (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    workout_session_id TEXT NOT NULL REFERENCES workout_session(id) ON DELETE CASCADE,
    set_number INT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    end_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (workout_session_id, set_number)
);

CREATE TABLE IF NOT EXISTS workout_session_set_rep (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    workout_session_set_id TEXT NOT NULL REFERENCES workout_session_set(id) ON DELETE CASCADE,
    rep_number INT NOT NULL,
    start_time BIGINT NOT NULL,
    end_time BIGINT NOT NULL,
    samples JSONB NOT NULL,
    metrics JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (workout_session_set_id, rep_number)
);

-- +goose Down
DROP TABLE IF EXISTS workout_session_set;
DROP TABLE IF EXISTS workout_session_set_rep;
