package userhandler

import (
	"game-app/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetUserRoutes(e *echo.Echo) {

	//* Group
	userGroup := e.Group("/users")

	userGroup.GET("/profile", h.userProfile, middleware.Auth(h.authSrv, h.authConfig))
	userGroup.POST("/login", h.userLogin)
	userGroup.POST("/register", h.userRegister)

}
