package user

import (
	"context"
	services "myproject/pkg/client"
	"myproject/pkg/common/utility"
	"myproject/pkg/config"
	"myproject/pkg/handlers"
)

type Service interface {
	AddBlog(ctx context.Context, request handlers.AddBlog, username string) error
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
