-- +goose Up
CREATE TABLE IF NOT EXISTS users (
   id SERIAL PRIMARY KEY,
     name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_type VARCHAR(100) NOT NULL CHECK (user_type IN ('admin', 'customer')),
    password VARCHAR(100) NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE
    );

-- +goose Down
-- DROP TABLE IF EXISTS users;
