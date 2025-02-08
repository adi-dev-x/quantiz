package user

import (
	"context"
	"database/sql"
	"fmt"
	"myproject/pkg/handlers"
	"myproject/pkg/model"
	"strings"
)

type Repository interface {
	AddBlog(ctx context.Context, request handlers.AddBlog) error
	GetUserId(ctx context.Context, username string) (int64, error)
	InsertBlogCategoryAddBlog(ctx context.Context, blog handlers.AddBlog, id int64, tx *sql.Tx) error
	QueryExecute(ctx context.Context, query string) ([]model.Blog, error)
	UpdateBlog(ctx context.Context, id int64, req handlers.UpdateBlogRequest) error
	DeleteBlog(ctx context.Context, id int64) error
}

type repository struct {
	sql *sql.DB
}

func NewRepository(sqlDB *sql.DB) Repository {
	return &repository{
		sql: sqlDB,
	}
}
func (r *repository) QueryExecute(ctx context.Context, query string) ([]model.Blog, error) {
	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %v", err)
	}
	defer rows.Close()
	var blogs []model.Blog

	for rows.Next() {
		var blog model.Blog
		if err := blog.Scan(rows); err != nil {
			return nil, fmt.Errorf("could not scan blog: %v", err)
		}

		blogs = append(blogs, blog)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return blogs, nil
}
func (r *repository) UpdateBlog(ctx context.Context, id int64, req handlers.UpdateBlogRequest) error {
	query := "UPDATE blogs SET "
	updates := []string{}
	values := []interface{}{}
	paramCount := 1

	if req.Title != nil {
		updates = append(updates, fmt.Sprintf("title = $%d", paramCount))
		values = append(values, *req.Title)
		paramCount++
	}
	if req.Descriptions != nil {
		updates = append(updates, fmt.Sprintf("descriptions = $%d", paramCount))
		values = append(values, *req.Descriptions)
		paramCount++
	}
	if req.Content != nil {
		updates = append(updates, fmt.Sprintf("content = $%d", paramCount))
		values = append(values, *req.Content)
		paramCount++
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + fmt.Sprintf(" WHERE id = $%d", paramCount)
	values = append(values, id)
	fmt.Println(query)
	fmt.Println(values)
	_, err := r.sql.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to update blog: %w", err)
	}

	return nil
}

func (r *repository) GetUserId(ctx context.Context, email string) (int64, error) {
	var userID int64

	query := "SELECT id FROM users WHERE email = $1"
	err := r.sql.QueryRowContext(ctx, query, email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("failed to get user ID: %w", err)
	}

	return userID, nil
}

func (r *repository) InsertBlogCategoryAddBlog(ctx context.Context, blog handlers.AddBlog, id int64, tx *sql.Tx) error {
	query := `INSERT INTO blog_categories (blog_id, category_id) VALUES `
	values := []interface{}{}
	placeholderCount := 1

	for i, val := range blog.Categories {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d)", placeholderCount, placeholderCount+1)
		values = append(values, id, val)
		placeholderCount += 2
	}

	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute bulk insert: %w", err)
	}

	return nil
}

func (r *repository) AddBlog(ctx context.Context, blog handlers.AddBlog) error {
	tx, err := r.sql.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `
	   INSERT INTO blogs (title, descriptions, content, user_id)
	   VALUES ($1, $2, $3, $4)
	   RETURNING id
	`

	var id int64
	err = tx.QueryRowContext(ctx, query,
		blog.Title,
		blog.Descriptions,
		blog.Content,
		blog.UserID,
	).Scan(&id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert blog: %w", err)
	}

	err = r.InsertBlogCategoryAddBlog(ctx, blog, id, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert blog categories: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) DeleteBlog(ctx context.Context, id int64) error {
	query := "DELETE FROM blogs where id = $1 "

	fmt.Println(query)

	_, err := r.sql.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update blog: %w", err)
	}

	return nil
}
