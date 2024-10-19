-- +goose Up

ALTER TABLE users
DROP COLUMN token,
DROP COLUMN refresh_token;


-- +goose Down

ALTER TABLE users
ADD COLUMN token TEXT,
ADD COLUMN refresh_token TEXT;
