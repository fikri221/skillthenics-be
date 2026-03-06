-- name: CreateSession :execresult
INSERT INTO ms_sessions (
    id, user_id, refresh_token, user_agent, client_ip, expires_at
) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetSessionByRefreshToken :one
SELECT * FROM ms_sessions WHERE refresh_token = ? LIMIT 1;

-- name: GetSession :one
SELECT * FROM ms_sessions WHERE id = ? LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM ms_sessions WHERE id = ?;

-- name: CleanupExpiredSessions :exec
DELETE FROM ms_sessions WHERE expires_at < NOW();

-- name: BlockSession :exec
UPDATE ms_sessions SET is_blocked = TRUE WHERE id = ?;
