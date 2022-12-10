package injector

import (
	"github.com/joisandresky/go-echo-mongodb-boilerplate/configs"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/handler"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/middleware"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/repository"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/routes"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/pkg/mongodb"
	"github.com/labstack/echo/v4"
)

type Injector interface {
	InjectModules()
}

type injector struct {
	mongoConn mongodb.MongoConnection
	cfg       *configs.Config
	server    *echo.Echo
}

func NewInjector(mongoConn mongodb.MongoConnection, cfg *configs.Config, server *echo.Echo) Injector {
	return &injector{mongoConn, cfg, server}
}

func (i *injector) InjectModules() {
	authMw := middleware.NewAuthMiddleware()

	humanRepo := repository.NewHumanRepository(i.mongoConn, i.cfg)
	humanHandler := handler.NewHumanHandler(humanRepo)
	humanRoutes := routes.NewHumanRoutes(humanHandler)

	humanRoutes.Install(i.server, authMw)
}
