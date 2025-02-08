package admin

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	services "myproject/pkg/client"
	"myproject/pkg/common/utility"
	"myproject/pkg/handlers"
)

type Service interface {
	Register(ctx context.Context, request handlers.Register) error
	DeleteBlog(ctx context.Context, username string, id int64) error
}

type service struct {
	repo     Repository
	services services.Services
}

func NewService(repo Repository, services services.Services) Service {
	return &service{
		repo:     repo,
		services: services,
	}
}

func (s *service) Register(ctx context.Context, request handlers.Register) error {
	var err error

	if !utility.IsValidPhoneNumber(request.Phone) {
		fmt.Println("this is in the service error invalid phone number")
		err = fmt.Errorf("invalid phone number")
		return err
	}

	existingUser, err := s.repo.Login(ctx, request.Email)

	if !existingUser.Verified && existingUser.Name != "" {
		return nil
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("this is in the service error checking existing user")
		err = fmt.Errorf("failed to check existing user: %w", err)
		return err
	}
	if existingUser.Email != "" {
		fmt.Println("this is in the service user already exists")
		err = fmt.Errorf("user already exists")
		return err
	}
	fmt.Println("this is in the service Register", request.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("this is in the service error hashing password")
		err = fmt.Errorf("failed to hash password: %w", err)
		return err
	}
	request.Password = string(hashedPassword)
	fmt.Println("this is in the service Register", request.Password)
	return s.repo.Register(ctx, request)
	//return nil
}

func (s *service) DeleteBlog(ctx context.Context, username string, id int64) error {

	userType, err := s.repo.GetUserDetails(ctx, username)
	if err != nil {
		return err
	}

	if userType == "admin" {
		return fmt.Errorf("blogs not authorized to edit: %w", err)
	}

	return s.repo.DeleteBlog(ctx, id)

}
