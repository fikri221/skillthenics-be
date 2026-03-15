-- +goose Up
ALTER TABLE workout_sessions ADD COLUMN deleted_at TIMESTAMP NULL;
ALTER TABLE workout_exercises ADD COLUMN deleted_at TIMESTAMP NULL;
ALTER TABLE exercise_sets ADD COLUMN deleted_at TIMESTAMP NULL;

-- +goose Down
ALTER TABLE exercise_sets DROP COLUMN deleted_at;
ALTER TABLE workout_exercises DROP COLUMN deleted_at;
ALTER TABLE workout_sessions DROP COLUMN deleted_at;
