-- +goose Up
CREATE TABLE workout_exercises (
    `id` VARCHAR(27) PRIMARY KEY,
    `workout_session_id` VARCHAR(27) NOT NULL,
    `exercise_id` VARCHAR(27) NOT NULL,
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`workout_session_id`) REFERENCES workout_sessions(`id`)
    ON DELETE CASCADE,
    FOREIGN KEY (`exercise_id`) REFERENCES exercises(`id`)
    ON DELETE CASCADE,

    INDEX idx_workout_session_id (workout_session_id),
    INDEX idx_exercise_id (exercise_id)
);

-- +goose Down
DROP TABLE workout_exercises;
