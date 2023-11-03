package claims

import (
	"game-app/service/authservice"

	"github.com/labstack/echo/v4"
)

func GetClaimFromEchoCTX(c echo.Context) *authservice.Claims {
	claims := c.Get("user")
	//fmt.Println("claims", claims)
	cl, ok := claims.(*authservice.Claims)
	if !ok {
		panic("claim was not found")
	}

	return cl

}
