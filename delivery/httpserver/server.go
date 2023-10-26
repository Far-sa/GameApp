package httpserver

import (
	"fmt"
	"game-app/config"
	"game-app/delivery/httpserver/userhandler"
	"game-app/service/authservice"
	"game-app/service/userservice"
	"game-app/validator/uservalidator"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config      config.Config
	userHandler userhandler.Handler
}

func New(config config.Config, authSrv authservice.Service,
	userSrv userservice.Service, userValidator uservalidator.Validator) Server {
	return Server{
		config:      config,
		userHandler: userhandler.New(config.Auth, authSrv, userSrv, userValidator),
	}
}

func (s Server) Serve() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//* User Routes
	s.userHandler.SetUserRoutes(e)

	// Start server
	log.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
