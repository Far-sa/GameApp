package httpserver

import (
	"fmt"
	"game-app/config"
	"game-app/delivery/httpserver/backofficeuserhandler"
	"game-app/delivery/httpserver/userhandler"
	"game-app/service/authorizationservice"
	"game-app/service/authservice"
	"game-app/service/backofficeuserservice"
	"game-app/service/userservice"
	"game-app/validator/uservalidator"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config                config.Config
	userHandler           userhandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
}

func New(config config.Config, authSrv authservice.Service,
	userSrv userservice.Service, userValidator uservalidator.Validator,
	backofficeUserSvc backofficeuserservice.Service,
	authorizationSvc authorizationservice.Service) Server {
	return Server{
		config:                config,
		userHandler:           userhandler.New(config.Auth, authSrv, userSrv, userValidator),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSrv, backofficeUserSvc, authorizationSvc),
	}
}

func (s Server) Serve() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//* User Routes
	s.userHandler.SetRoutes(e)
	s.backofficeUserHandler.SetRoutes(e)

	// Start server
	log.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
