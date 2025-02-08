package admin

import (
	"context"
	"database/sql"
	"fmt"
	"myproject/pkg/handlers"
	"myproject/pkg/model"
)

type Repository interface {
	Register(ctx context.Context, request handlers.Register) error
	Login(ctx context.Context, email string) (model.UserDetails, error)
	VerifyUser(ctx context.Context, email string) error
	GetUserDetails(ctx context.Context, username string) (string, error)
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

func (r *repository) VerifyUser(ctx context.Context, email string) error {
	query := `UPDATE users SET verified = TRUE WHERE email = $1`
	_, err := r.sql.ExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("failed to update user verification status: %w", err)
	}

	return nil
}

func (r *repository) Register(ctx context.Context, request handlers.Register) error {

	query := `INSERT INTO users (name, email, password,phone,user_type) VALUES ($1, $2, $3, $4,$5)`
	_, err := r.sql.ExecContext(ctx, query, request.Name, request.Email, request.Password, request.Phone, "admin")
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

func (r *repository) Login(ctx context.Context, email string) (model.UserDetails, error) {
	query := `SELECT name, email, password,verified FROM users WHERE email = $1`

	var user model.UserDetails
	err := r.sql.QueryRowContext(ctx, query, email).Scan(&user.Name, &user.Email, &user.Password, &user.Verified)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.UserDetails{}, nil
		}
		return model.UserDetails{}, fmt.Errorf("failed to find user by email: %w", err)
	}
	fmt.Println("the data !!!! ", user)

	return user, nil
}
func (r *repository) GetUserDetails(ctx context.Context, username string) (string, error) {
	var userType string

	query := "SELECT user_type FROM users WHERE email = $1"
	err := r.sql.QueryRowContext(ctx, query, username).Scan(&userType)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to get user ID: %w", err)
	}

	return userType, nil
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
