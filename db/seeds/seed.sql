INSERT INTO exercise_types (id, name)
VALUES (gen_random_uuid(), 'squat');

INSERT INTO users (id, clerk_id, first_name, last_name, email)
VALUES (gen_random_uuid(), 'clerk_id', 'Test', 'User', 'test@example.com');