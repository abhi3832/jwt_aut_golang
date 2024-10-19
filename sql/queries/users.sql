-- name: CreateUser :one

INSERT INTO users(user_id,first_name,last_name,email,phone,passward,user_type,created_at,updated_at,api_key)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,encode(sha256(random()::text::bytea),'hex'))
RETURNING *;

-- name: GetUserByApiKey :one
SELECT * FROM users WHERE api_key = $1;

-- name: CheckUserAlreadyExits :one
SELECT * FROM users WHERE email = $1 or phone = $2;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByUserId :one
SELECT * FROM users WHERE user_id = $1;
