package middleware

import (
	"game-app/pkg/claims"
	"game-app/pkg/errs"
	"game-app/service/authorizationservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AccessCheck(service authorizationservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claims.GetClaimFromEchoCTX(c)
			isAllowed, err := service.CheckAccess(claims.UserID, claims.Role)
			if err != nil {
				// TODO -> log unexpected error
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errs.ErrorMsgSomethingWrong,
				})
			}

			if !isAllowed {
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": errs.ErrorMsgAccessDenied,
				})
			}

			return next(c)
		}
	}
}
