package di

import (
	"github.com/go-playground/validator/v10"
	"log"
	"myproject/pkg/admin"
	bootserver "myproject/pkg/boot"
	services "myproject/pkg/client"
	"myproject/pkg/config"
	db "myproject/pkg/database"
	"myproject/pkg/middleware"
	"myproject/pkg/migrations"
	"myproject/pkg/user"
	common "myproject/pkg/usualprivilage"
)

func InitializeEvent(conf config.Config) (*bootserver.ServerHttp, error) {

	sqlDB, err := db.ConnectPGDB(conf)
	if err != nil {
		return nil, err
	}
	err = migrations.RunMigrations(conf)
	if err != nil {
		log.Fatalln("Issue during migration 🧘: ", err.Error())
	}

	userRepository := user.NewRepository(sqlDB)
	myService := services.MyService{Config: conf}
	userService := user.NewService(userRepository, myService)
	validate := validator.New()
	admjwt := middleware.MiddlewareJWT{Config: conf}
	userHandler := user.NewHandler(userService, myService, admjwt, conf)
	adminRepository := admin.NewRepository(sqlDB)
	myService2 := services.MyService{Config: conf}
	adminService := admin.NewService(adminRepository, myService2)
	adminHandler := admin.NewHandler(adminService, myService, admjwt, validate)

	commonRepository := common.NewRepository(sqlDB)
	//myService2 := services.MyService{Config: conf}
	commonService := common.NewService(commonRepository, myService2)
	commonHandler := common.NewHandler(commonService, myService, admjwt, validate)
	serverHttp := bootserver.NewServerHttp(*commonHandler, *userHandler, *adminHandler)

	return serverHttp, nil
}
