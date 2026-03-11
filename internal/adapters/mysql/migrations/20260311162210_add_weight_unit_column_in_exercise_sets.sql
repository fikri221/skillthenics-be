-- +goose Up
ALTER TABLE exercise_sets ADD COLUMN weight_unit VARCHAR(10);

-- +goose Down
ALTER TABLE exercise_sets DROP COLUMN weight_unit;
