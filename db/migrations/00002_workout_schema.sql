-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS workout_type (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    name TEXT NOT NULL,
    description TEXT,
    config JSONB NOT NULL, -- Store workout-specific configuration in JSON format
    image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workout_session (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    workout_type_id TEXT NOT NULL REFERENCES workout_type(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workout_set (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    workout_session_id TEXT NOT NULL REFERENCES workout_session(id) ON DELETE CASCADE,
    set_number INT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    coaching_score JSONB, -- Store coaching feedback in JSON format
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workout_rep (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    workout_set_id TEXT NOT NULL REFERENCES workout_set(id) ON DELETE CASCADE,
    rep_number INT NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    quality TEXT, -- "good", "okay", "bad"
    peak_angle FLOAT,
    avg_flex FLOAT, -- Average flex sensor value
    sensor_data_points JSONB, -- Full array of SensorData[] for the rep
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_calibration (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id TEXT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    standing_angle FLOAT NOT NULL,
    standing_flex FLOAT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_calibration;
DROP TABLE IF EXISTS workout_rep;
DROP TABLE IF EXISTS workout_set;
DROP TABLE IF EXISTS workout_session;
DROP TABLE IF EXISTS workout_type;
-- +goose StatementEnd
