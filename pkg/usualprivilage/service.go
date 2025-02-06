package admin

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	services "myproject/pkg/client"
	"myproject/pkg/common/utility"
	db "myproject/pkg/database"
	"myproject/pkg/handlers"
)

type Service interface {
	Register(ctx context.Context, request handlers.Register) error
	Login(ctx context.Context, request handlers.Login) error
	Listing(ctx context.Context, conditions []db.WhereCondition, pageCount int) (interface{}, error)
	OtpLogin(ctx context.Context, request handlers.Otp) error

	///product listing
	//ProductListing(ctx context.Context) ([]model.ProductListingUsers, error)
	//PlowListing(ctx context.Context, id string) ([]model.ProductListingUsers, error)

	///orders
	/// Singlevendor
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

func (s *service) Login(ctx context.Context, request handlers.Login) error {
	fmt.Println("this is in the service Login", request.Password)
	var err error

	storedUser, err := s.repo.Login(ctx, request.Email)
	fmt.Println("thisss is the dataaa ", storedUser)
	if err != nil {
		fmt.Println("this is in the service user not found")
		return fmt.Errorf("user not found: %w", err)
	}
	if !storedUser.Verified {
		return fmt.Errorf("user not found: %w", "User is not verified")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(request.Password)); err != nil {
		fmt.Println("this is in the service incorrect password")
		return fmt.Errorf("incorrect password: %w", err)
	}

	return nil

	//return s.repo.Login(ctx, request)
}
func (s *service) OtpLogin(ctx context.Context, request handlers.Otp) error {

	var err error
	if request.Email == "" || request.Otp == "" {
		fmt.Println("this is in the service error value missing")
		err = fmt.Errorf("missing values")
		return err
	}
	err = s.repo.VerifyUser(ctx, request.Email)
	if err != nil {

		return fmt.Errorf("user not found: %w", err)
	}

	return nil

	//return s.repo.Login(ctx, request)
}

func (s *service) Listing(ctx context.Context, conditions []db.WhereCondition, pageCount int) (interface{}, error) {

	query, Joins := db.QueryBuilder(conditions, "JOIN")
	if Joins != "" {
		query = "SELECT blogs.*,users.name FROM blogs JOIN users ON blogs.user_id = users.id " + query
	} else {
		query = "SELECT blogs.*,users.name FROM blogs JOIN users ON blogs.user_id = users.id " + query
	}
	fmt.Println("Generated query !!!", query)
	res, err := s.repo.QueryExecute(ctx, query)
	if err != nil {
		fmt.Println("this is in the service error querying blogs", err.Error())
		return nil, fmt.Errorf("blogs not found: %w", err)
	}

	return res, nil
	//	return s.repo.Listing(ctx)
}
