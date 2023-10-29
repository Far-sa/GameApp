package backofficeuserhandler

import (
	"game-app/service/authorizationservice"
	"game-app/service/authservice"
	"game-app/service/backofficeuserservice"
)

type Handler struct {
	authCfg           authservice.Config
	authSvc           authservice.Service
	backofficeUserSvc backofficeuserservice.Service
	authorizationSvc  authorizationservice.Service
}

func New(authCfg authservice.Config, authSvc authservice.Service,
	backofficeUserSvc backofficeuserservice.Service, authorizationSvc authorizationservice.Service) Handler {
	return Handler{
		authCfg:           authCfg,
		authSvc:           authSvc,
		backofficeUserSvc: backofficeUserSvc,
		authorizationSvc:  authorizationSvc,
	}
}
