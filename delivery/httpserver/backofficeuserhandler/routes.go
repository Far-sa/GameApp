package backofficeuserhandler

import (
	"game-app/delivery/httpserver/middleware"
	"game-app/entity"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroups := e.Group("/backoffice/users")

	userGroups.GET("/", h.ListUsers, middleware.Auth(h.authSvc, h.authCfg),
		middleware.AccessCheck(h.authorizationSvc, entity.UserListPermission))
}
