package middleware

import (
	"fmt"
	"game-app/param"
	"game-app/pkg/claims"
	"game-app/pkg/errs"
	"game-app/pkg/timestamp"
	"game-app/service/presenceservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpsertPresence(service presenceservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			claims := claims.GetClaimFromEchoCTX(c)

			_, err = service.UpsertPresence(c.Request().Context(), param.UpsertPresenceRequest{
				UserID:    claims.UserID,
				Timestamp: timestamp.Now(),
			})
			if err != nil {
				//TODO : log an unexpected error
				fmt.Println("UpsertPresence err :", err.Error())

				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errs.ErrorMsgSomethingWrong,
				})
			}

			return next(c)
		}
	}
}
