-- name: CreateSession :one

INSERT INTO sessions(user_id, refresh_token, created_at, updated_at, expires_at)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;

-- name: GetSessionByUserId :one
SELECT * FROM sessions where user_id = $1;

-- name: UpdateSession :exec
UPDATE sessions
SET
    refresh_token = 'new_refresh_token',
    expires_at = '2024-12-31 23:59:59',
    updated_at = NOW()
WHERE
    user_id = $1; 

-- name: UpdateSessionByDelete :exec

DELETE FROM sessions where user_id = $1;
