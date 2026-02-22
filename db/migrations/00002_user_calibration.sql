-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS user_calibration (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
    user_id TEXT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,

    standing_yaw_angle DOUBLE PRECISION NOT NULL,
    standing_pitch_angle DOUBLE PRECISION NOT NULL,
    standing_roll_angle DOUBLE PRECISION NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_calibration;
-- +goose StatementEnd
