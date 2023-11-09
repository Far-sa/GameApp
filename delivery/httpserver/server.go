package httpserver

import (
	"fmt"
	"game-app/config"
	"game-app/delivery/httpserver/backofficeuserhandler"
	"game-app/delivery/httpserver/matchinghandler"
	"game-app/delivery/httpserver/userhandler"
	"game-app/logger"
	"game-app/service/authorizationservice"
	"game-app/service/authservice"
	"game-app/service/backofficeuserservice"
	"game-app/service/matchingservice"
	"game-app/service/presenceservice"
	"game-app/service/userservice"
	"game-app/validator/matchingvalidator"
	"game-app/validator/uservalidator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Server struct {
	config                config.Config
	userHandler           userhandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
	matchingHandler       matchinghandler.Handler
	Router                *echo.Echo
}

func New(config config.Config,
	authSrv authservice.Service,
	userSrv userservice.Service,
	userValidator uservalidator.Validator,
	backofficeUserSvc backofficeuserservice.Service,
	authorizationSvc authorizationservice.Service,
	matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator,
	presenceSvc presenceservice.Service) Server {
	return Server{
		Router:                echo.New(),
		config:                config,
		userHandler:           userhandler.New(config.Auth, authSrv, userSrv, userValidator, presenceSvc),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSrv, backofficeUserSvc, authorizationSvc),
		matchingHandler:       matchinghandler.New(config.Auth, authSrv, matchingSvc, matchingValidator, presenceSvc),
	}
}

func (s Server) Serve() {

	// Middleware
	//s.Router.Use(middleware.Logger())

	s.Router.Use(middleware.RequestID())

	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID:     true,
		LogLatency:       true,
		LogContentLength: true,
		LogURI:           true,
		LogHost:          true,
		LogStatus:        true,
		LogMethod:        true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogError:         true,
		LogResponseSize:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			errMsg := ""
			if v.Error != nil {
				errMsg = v.Error.Error()
			}
			logger.Logger.Named("http-server").Info("request",
				zap.String("REQUEST_ID", v.RequestID),
				zap.String("HOST", v.Host),
				zap.String("CONTENT_LENGHT", v.ContentLength),
				zap.String("PROTOCOLE", v.Protocol),
				zap.String("METHOD", v.Method),
				zap.String("URI", v.URI),
				zap.Duration("LATENCY", v.Latency),
				zap.String("ERROR", errMsg),
				zap.String("REMOTE_IP", v.RemoteIP),
				zap.Int64("RESPONSE_SIZE", v.ResponseSize),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))

	s.Router.Use(middleware.Recover())

	//* User Routes
	s.userHandler.SetRoutes(s.Router)
	s.backofficeUserHandler.SetRoutes(s.Router)
	s.matchingHandler.SetRoutes(s.Router)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router error:", err)
	}
}
