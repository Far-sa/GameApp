package userhandler

import (
	"fmt"
	"game-app/param"
	"game-app/pkg/httpmsg"
	"game-app/service/authservice"
	"game-app/service/userservice"
	"game-app/validator/uservalidator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	authConfig    authservice.Config
	authSrv       authservice.Service
	userSrv       userservice.Service
	userValidator uservalidator.Validator
}

func New(authConfig authservice.Config, authSrv authservice.Service, userSrv userservice.Service, userValidator uservalidator.Validator) Handler {
	return Handler{
		authConfig:    authConfig,
		authSrv:       authSrv,
		userSrv:       userSrv,
		userValidator: userValidator,
	}
}

func (h Handler) userRegister(c echo.Context) error {

	var req param.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	//use validator pkg
	if fieldErros, err := h.userValidator.ValidateRegisterRequest(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErros,
		})
	}

	resp, err := h.userSrv.Register(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h Handler) userLogin(c echo.Context) error {

	var req param.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	//* validate Login
	if _, err := h.userValidator.ValidateLoginRequest(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"messsage": msg,
			"errors":   err,
		})
	}

	resp, err := h.userSrv.Login(req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
		//return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)

}

func getClaims(c echo.Context) *authservice.Claims {
	claims := c.Get("user")
	fmt.Println("claims", claims)
	cl, ok := claims.(*authservice.Claims)
	if !ok {
		panic("claim was not found")
	}

	return cl
}

func (h Handler) userProfile(c echo.Context) error {

	claims := getClaims(c)

	resp, err := h.userSrv.Profile(param.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
		//return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

// func (h Handler) healthCheck(c echo.Context) error {
// 	return c.JSON(http.StatusOK, echo.Map{
// 		"message": "everything is fine",
// 	})
// }
