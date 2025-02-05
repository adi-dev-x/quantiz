package bootserver

import (
	"github.com/gofiber/fiber/v2"
	"myproject/pkg/admin"
	"myproject/pkg/config"
	"myproject/pkg/user"
	common "myproject/pkg/usualprivilage"
)

type ServerHttp struct {
	app *fiber.App
}

func NewServerHttp(commonHandler common.Handler, userHandler user.Handler, adminHandler admin.Handler) *ServerHttp {
	app := fiber.New()

	// Mount user routes
	commonHandler.MountRoutes(app)
	userHandler.MountRoutes(app)

	// Mount vendor routes

	adminHandler.MountRoutes(app)
	//return &ServerHttp{Engine: engine}
	return &ServerHttp{app}
}

func (s *ServerHttp) Start(conf config.Config) {
	s.app.Listen(conf.Host + ":" + conf.ServerPort)
}
func (s *ServerHttp) Kill() error {
	return s.app.Shutdown()
}
