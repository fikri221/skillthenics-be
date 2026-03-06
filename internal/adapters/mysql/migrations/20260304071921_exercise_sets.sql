-- +goose Up
CREATE TABLE exercise_sets (
    `id` VARCHAR(27) PRIMARY KEY,
    `workout_exercise_id` VARCHAR(27) NOT NULL,
    `set_number` INT NOT NULL,
    `reps` INT,
    `weight` DECIMAL(10, 2),
    `rest_seconds` INT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`workout_exercise_id`) REFERENCES workout_exercises(`id`)
    ON DELETE CASCADE,

    INDEX idx_workout_exercise_id (workout_exercise_id)
);

-- +goose Down
DROP TABLE exercise_sets;
