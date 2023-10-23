package httpserver

import (
	"fmt"
	"game-app/config"
	"game-app/service/authservice"
	"game-app/service/userservice"
	"game-app/validator/uservalidator"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config        config.Config
	authSrv       authservice.Service
	userSrv       userservice.Service
	userValidator uservalidator.Validator
}

func New(config config.Config, authSrv authservice.Service, userSrv userservice.Service, userValidator uservalidator.Validator) Server {
	return Server{
		config:        config,
		authSrv:       authSrv,
		userSrv:       userSrv,
		userValidator: userValidator}
}

func (s Server) Serve() {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//* Group
	userGroup := e.Group("/users")

	e.GET("/health", s.healthCheck)
	userGroup.POST("/login", s.userLogin)
	userGroup.POST("/register", s.userRegister)
	userGroup.GET("/profile", s.userProfile)

	// Start server
	log.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HTTPServer.Port)))
}
