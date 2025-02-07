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
	//Listing(ctx context.Context) ([]model.Coupon, error)
	Login(ctx context.Context, email string) (model.UserDetails, error)
	VerifyUser(ctx context.Context, email string) error
	QueryExecute(ctx context.Context, query string) ([]model.Blog, error)
	//Product listing

}

type repository struct {
	sql *sql.DB
}

func NewRepository(sqlDB *sql.DB) Repository {
	return &repository{
		sql: sqlDB,
	}
}

func (r *repository) ListCategoriesBlog(ctx context.Context, id int64) ([]model.Category, error) {
	query := `SELECT c.* FROM blog_categories bc LEFT JOIN categories c 
             ON bc.category_id = c.id WHERE bc.blog_id = $1;`
	var categories []model.Category
	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var category model.Category
		if err := category.Scan(rows); err != nil {
			return nil, fmt.Errorf("could not scan blog: %v", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
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
		categories, err := r.ListCategoriesBlog(ctx, blog.ID)
		if err != nil {
			fmt.Printf("could not list categories blog: %v", err)
		}
		fmt.Println("printing categories ", categories)
		blog.Categories = categories
		blogs = append(blogs, blog)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return blogs, nil
}
func (r *repository) VerifyUser(ctx context.Context, email string) error {
	fmt.Println("this is in the repository VerifyUser")

	query := `UPDATE users SET verified = TRUE WHERE email = $1`
	_, err := r.sql.ExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("failed to update user verification status: %w", err)
	}

	return nil
}

func (r *repository) Register(ctx context.Context, request handlers.Register) error {
	fmt.Println("this is in the repository Register")
	query := `INSERT INTO users (name, email, password,phone,user_type) VALUES ($1, $2, $3, $4,$5)`
	_, err := r.sql.ExecContext(ctx, query, request.Name, request.Email, request.Password, request.Phone, "admin")
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

func (r *repository) Login(ctx context.Context, email string) (model.UserDetails, error) {
	fmt.Println("theee !!!!!!!!!!!  LLLLoginnnnnn  ", email)
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
