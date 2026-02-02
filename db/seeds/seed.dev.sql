INSERT INTO users (id, first_name, last_name, email)
VALUES ('user_396cTR8sDnxPxTk8VS4tMfCgd8k', 'Test', 'User', 'test@example.com')
ON CONFLICT (id) DO NOTHING;

-- 1. Create a fake session from yesterday
INSERT INTO workout_session (id, user_id, workout_type_id, start_time, end_time)
VALUES (
    'dev-session-001', 
    'user_396cTR8sDnxPxTk8VS4tMfCgd8k', 
    'workout-squat-001', 
    NOW() - INTERVAL '1 day', 
    NOW() - INTERVAL '23 hours'
) ON CONFLICT DO NOTHING;

-- 2. Create a Set for that session
INSERT INTO workout_set (id, workout_session_id, set_number, start_time, end_time, coaching_score)
VALUES (
    'dev-set-001', 
    'dev-session-001', 
    1, 
    NOW() - INTERVAL '23 hours 55 minutes', 
    NOW() - INTERVAL '23 hours 50 minutes',
    '{"overallQuality": "good", "averageDepth": 82.5}'::jsonb
) ON CONFLICT DO NOTHING;

-- 3. Create 3 Reps for that set
INSERT INTO workout_rep (workout_set_id, rep_number, start_time, end_time, quality, peak_angle, avg_flex, sensor_data_points)
VALUES 
('dev-set-001', 1, NOW() - INTERVAL '23 hours 54 minutes', NOW() - INTERVAL '23 hours 53 minutes 58 seconds', 'good', 85.0, 110.0, '[{"t":0, "v":0}, {"t":1, "v":85}]'::jsonb),
('dev-set-001', 2, NOW() - INTERVAL '23 hours 53 minutes', NOW() - INTERVAL '23 hours 52 minutes 58 seconds', 'good', 82.0, 112.0, '[{"t":0, "v":0}, {"t":1, "v":82}]'::jsonb),
('dev-set-001', 3, NOW() - INTERVAL '23 hours 52 minutes', NOW() - INTERVAL '23 hours 51 minutes 58 seconds', 'okay', 70.0, 115.0, '[{"t":0, "v":0}, {"t":1, "v":70}]'::jsonb)
ON CONFLICT DO NOTHING;

