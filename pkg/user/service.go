package user

import (
	"context"
	"fmt"
	services "myproject/pkg/client"
	"myproject/pkg/common/utility"
	"myproject/pkg/config"
	db "myproject/pkg/database"
	"myproject/pkg/handlers"
)

type Service interface {
	AddBlog(ctx context.Context, request handlers.AddBlog, username string) error
	UpdateBlog(ctx context.Context, request handlers.UpdateBlogRequest, username string, id int64) error
}
type service struct {
	repo     Repository
	Config   config.Config
	services services.Services
}

func NewService(repo Repository, services services.Services) Service {
	return &service{
		repo:     repo,
		services: services,
	}
}
func (s *service) AddBlog(ctx context.Context, request handlers.AddBlog, username string) error {

	request.ID = utility.UniqueId()

	userId, err := s.repo.GetUserId(ctx, username)
	if err != nil {
		return err
	}
	request.UserID = userId
	return s.repo.AddBlog(ctx, request)

}
func (s *service) UpdateBlog(ctx context.Context, request handlers.UpdateBlogRequest, username string, id int64) error {

	userId, err := s.repo.GetUserId(ctx, username)
	if err != nil {
		return err
	}
	var conditions []db.WhereCondition
	conditions = append(conditions, db.WhereCondition{
		Key:       "id",
		Value:     string(id),
		Condition: "=",
		Table:     "blogs",
		Joins:     "",
	})
	conditions = append(conditions, db.WhereCondition{
		Key:       "id",
		Value:     string(userId),
		Condition: "=",
		Table:     "users",
		Joins:     "",
	})
	query, _ := db.QueryBuilder(conditions, "JOIN")
	query = "SELECT blogs.*,users.name FROM blogs JOIN users ON blogs.user_id = users.id " + query
	res, err := s.repo.QueryExecute(ctx, query)
	if err != nil {
		fmt.Println("this is not able to fetch details", err.Error())
		return fmt.Errorf("blogs not found: %w", err)
	}
	if len(res) == 0 {
		return fmt.Errorf("blogs not authorized to edit: %w", err)
	}

	return s.repo.UpdateBlog(ctx, id, request)

}
