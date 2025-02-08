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
	if req.Title != nil {
		updates = append(updates, "title = ?")
		values = append(values, *req.Title)
	}
	if req.Descriptions != nil {
		updates = append(updates, "descriptions = ?")
		values = append(values, *req.Descriptions)
	}
	if req.Author != nil {
		updates = append(updates, "author = ?")
		values = append(values, *req.Author)
	}
	if req.Content != nil {
		updates = append(updates, "content = ?")
		values = append(values, *req.Content)
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields provided for update")
	}

	query += strings.Join(updates, ", ") + " WHERE id = ?"
	values = append(values, id)

	_, err := r.sql.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to update blog: %w", err)
	}

	if req.Categories != nil {
		err = r.updateBlogCategories(ctx, id, req.Categories)
		if err != nil {
			return fmt.Errorf("failed to update blog categories: %w", err)
		}
	}

	return nil
}
func (r *repository) updateBlogCategories(ctx context.Context, blogID int64, categories []int64) error {
	_, err := r.sql.ExecContext(ctx, "DELETE FROM blog_categories WHERE blog_id = ?", blogID)
	if err != nil {
		return fmt.Errorf("failed to delete existing categories: %w", err)
	}
	query := "INSERT INTO blog_categories (blog_id, category_id) VALUES "
	values := []interface{}{}
	for i, categoryID := range categories {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?)"
		values = append(values, blogID, categoryID)
	}

	_, err = r.sql.ExecContext(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to insert new categories: %w", err)
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
