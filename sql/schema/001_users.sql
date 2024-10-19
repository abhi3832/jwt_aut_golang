-- +goose Up

CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    passward TEXT NOT NULL, 
    user_type VARCHAR(50),
    token TEXT,
    refresh_token TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT(encode(sha256(random()::text::bytea),'hex'))
);

-- +goose Down

DROP TABLE users;