-- +goose Up
CREATE TABLE IF NOT EXISTS  blog_categories (
    blog_id INT REFERENCES blogs(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (blog_id, category_id)
);

-- +goose Down

