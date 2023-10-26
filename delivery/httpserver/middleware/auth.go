package middleware

import (
	"game-app/service/authservice"

	mw "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth(service authservice.Service, config authservice.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {

			claims, err := service.VerifyToken(auth)
			if err != nil {
				return nil, err
			}

			return claims, nil
		},
		ContextKey: "user",
		SigningKey: []byte(config.SignKey),
	})
}
