/**
 * Seed data for development environment
 * Run with: psql -d your_database_name -f db/seeds/seed.dev.sql
 */
TRUNCATE TABLE 
    workout_session,
    workout_type,
    users
RESTART IDENTITY CASCADE;


-- Create a test user
INSERT INTO users (id, clerk_id, first_name, last_name, email)
VALUES ('test-user','user_396cTR8sDnxPxTk8VS4tMfCgd8k', 'Test', 'User', 'test@example.com')
ON CONFLICT (id) DO UPDATE SET
  clerk_id = EXCLUDED.clerk_id,
  first_name = EXCLUDED.first_name,
  last_name = EXCLUDED.last_name,
  email = EXCLUDED.email;

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

-- Create a fake session from yesterday
INSERT INTO workout_session (id, user_id, workout_type_id, start_time, end_time)
VALUES (
    'dev-session-001', 
    'test-user', 
    'workout-squat-001', 
    NOW() - INTERVAL '1 day', 
    NOW() - INTERVAL '23 hours'
) ON CONFLICT (id) DO NOTHING;
