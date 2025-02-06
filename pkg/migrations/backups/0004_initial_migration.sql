-- +goose Up
CREATE TABLE IF NOT EXISTS  blogs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    descriptions TEXT NOT NULL,
    content TEXT NOT NULL,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- +goose Down

