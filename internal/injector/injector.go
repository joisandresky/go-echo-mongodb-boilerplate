package injector

import (
	"github.com/joisandresky/go-echo-mongodb-boilerplate/database"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/handler"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/middleware"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/repository"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/routes"
	"github.com/labstack/echo/v4"
)

type Injector interface {
	InjectModules()
}

type injector struct {
	conn   database.Connection
	server *echo.Echo
}

func NewInjector(conn database.Connection, server *echo.Echo) Injector {
	return &injector{conn, server}
}

func (i *injector) InjectModules() {
	authMw := middleware.NewAuthMiddleware()

	humanRepo := repository.NewHumanRepository(i.conn)
	humanHandler := handler.NewHumanHandler(humanRepo)
	humanRoutes := routes.NewHumanRoutes(humanHandler)

	humanRoutes.Install(i.server, authMw)
}
