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

func (r *repository) AddBlog(ctx context.Context, blog handlers.AddBlog) error {
	// Define the SQL query
	//query := `
	//    INSERT INTO blogs (title, descriptions, content, user_id, created_at, updated_at)
	//    VALUES ($1, $2, $3, $4, $5, $6)
	//    RETURNING id
	//`

	// Handle NULL descriptions

	// Execute the query
	//	var id int
	//err := r.sql.QueryRowContext(ctx, query,
	//	blog.Title,     // $1
	//	descriptions,   // $2
	//	blog.Content,   // $3
	//	blog.UserID,    // $4
	//	blog.CreatedAt, // $5
	//	blog.UpdatedAt, // $6
	//).Scan(&id)
	//if err != nil {
	//	return 0, fmt.Errorf("could not insert blog: %v", err)
	//}

	// Return the ID of the newly inserted blog
	return nil
}
