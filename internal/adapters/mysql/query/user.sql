-- name: ListUsers :many
SELECT * FROM ms_user WHERE rec_status = 'A';

-- name: GetUserByEmail :one
SELECT * FROM ms_user WHERE email = ? AND rec_status = 'A' LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM ms_user WHERE id = ? AND rec_status = 'A';

-- name: CreateUser :execresult
INSERT INTO ms_user (name, email, password, rec_status) VALUES (?, ?, ?, 'A');

-- name: UpdateUser :execresult
UPDATE ms_user SET name = ?, email = ?, password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND rec_status = 'A';

-- name: DeleteUser :execresult
UPDATE ms_user SET rec_status = 'D', updated_at = CURRENT_TIMESTAMP WHERE id = ?;