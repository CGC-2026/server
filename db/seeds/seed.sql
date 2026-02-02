
-- Seed data for workout_type table
INSERT INTO workout_type (id, name, description, config)
VALUES (
    'workout-squat-001',
    'Squat',
    'Track your squat form and depth with real-time coaching',
    '{
        "sampleRate": 120,
        "minRepDuration": 800,
        "maxRepDuration": 8000,
        "minDepthAngle": 30
    }'::jsonb
)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    config = EXCLUDED.config;