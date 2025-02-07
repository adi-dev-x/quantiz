package user

import (
	"context"
	"database/sql"
	"fmt"
	"myproject/pkg/handlers"
	//"reflect"
	//"strings"
)

// ListWish
type Repository interface {
	AddBlog(ctx context.Context, request handlers.AddBlog) error
	GetUserId(ctx context.Context, username string) (int64, error)
	InsertBlogCategoryAddBlog(ctx context.Context, blog handlers.AddBlog, id int64, tx *sql.Tx) error
}

type repository struct {
	sql *sql.DB
}

func NewRepository(sqlDB *sql.DB) Repository {
	return &repository{
		sql: sqlDB,
	}
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
	// Initialize the query and values slice
	query := `INSERT INTO blog_categories (blog_id, category_id) VALUES `
	values := []interface{}{}
	placeholderCount := 1

	// Build the query and values slice
	for i, val := range blog.Categories {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d)", placeholderCount, placeholderCount+1)
		values = append(values, id, val)
		placeholderCount += 2
	}

	// Execute the bulk insert query
	_, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute bulk insert: %w", err)
	}

	return nil
}

func (r *repository) AddBlog(ctx context.Context, blog handlers.AddBlog) error {
	// Start a transaction
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
