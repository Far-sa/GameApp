package httpserver

import (
	"fmt"
	"game-app/config"
	"game-app/service/authservice"
	"game-app/service/userservice"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config  config.Config
	authSrv authservice.Service
	userSrv userservice.Service
}

func New(config config.Config, authSrv authservice.Service, userSrv userservice.Service) Server {
	return Server{config: config, authSrv: authSrv, userSrv: userSrv}
}

func (s Server) Serve() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/health", s.healthCheck)
	e.POST("/users/register", s.userRegister)

	// Start server
	log.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
