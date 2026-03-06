-- +goose Up
CREATE TABLE workout_sessions (
    `id` VARCHAR(27) PRIMARY KEY,
    `user_id` INT NOT NULL,
    `session_date` DATE NOT NULL,
    `duration_minutes` INT,
    `calories_burned` INT,
    `notes` TEXT,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES ms_user(`id`)
    ON DELETE CASCADE,
    INDEX idx_user_id (user_id)
);

-- +goose Down
DROP TABLE workout_sessions;
