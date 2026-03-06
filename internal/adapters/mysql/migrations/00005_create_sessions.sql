-- +goose Up
-- ms_sessions: Stores refresh tokens for rotation
CREATE TABLE IF NOT EXISTS ms_sessions (
    id VARCHAR(27) PRIMARY KEY,
    user_id INT NOT NULL,
    refresh_token TEXT NOT NULL,
    user_agent TEXT,
    client_ip VARCHAR(45),
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS ms_sessions;
