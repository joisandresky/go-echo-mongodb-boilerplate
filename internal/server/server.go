package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joisandresky/go-echo-mongodb-boilerplate/configs"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/internal/injector"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/pkg/helper"
	"github.com/joisandresky/go-echo-mongodb-boilerplate/pkg/mongodb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server interface {
	Run() error
}

type server struct {
	mongoConn mongodb.MongoConnection
	cfg       *configs.Config
}

func NewServer(mongoConn mongodb.MongoConnection, cfg *configs.Config) Server {
	return &server{mongoConn, cfg}
}

func (srv *server) Run() error {
	port := srv.cfg.App.Port

	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.Use((middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		},
	)))

	server.GET("/", func(c echo.Context) error {
		return helper.OkResponse(c, helper.Response{
			Message: fmt.Sprintf("Welcome to %s", srv.cfg.App.ServiceName),
		})
	})

	server.Any("/*", func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"code":    http.StatusNotFound,
			"message": "Route is Not Exist",
		})
	})

	injection := injector.NewInjector(srv.mongoConn, srv.cfg, server)
	injection.InjectModules()

	log.Println(fmt.Printf("%v is Running at port %v ... ", srv.cfg.App.ServiceName, port))

	// Start Server
	go func() {
		if err := server.Start(":" + port); err != nil && err != http.ErrServerClosed {
			server.Logger.Fatal("shutting down the My Service")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatal(err)
		return err
	}

	return nil
}
