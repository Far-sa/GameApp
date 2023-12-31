package userhandler

import (
	"context"
	"game-app/param"
	"game-app/pkg/claims"
	"game-app/pkg/httpmsg"
	"game-app/service/authservice"
	"game-app/service/presenceservice"
	"game-app/service/userservice"
	"game-app/validator/uservalidator"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	authConfig    authservice.Config
	authSrv       authservice.Service
	userSrv       userservice.Service
	userValidator uservalidator.Validator
	presenceSvc   presenceservice.Service
}

func New(authConfig authservice.Config, authSrv authservice.Service,
	userSrv userservice.Service, userValidator uservalidator.Validator,
	presenceSvc presenceservice.Service,
) Handler {
	return Handler{
		authConfig:    authConfig,
		authSrv:       authSrv,
		userSrv:       userSrv,
		userValidator: userValidator,
		presenceSvc:   presenceSvc,
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

func (h Handler) userProfile(c echo.Context) error {

	claims := claims.GetClaimFromEchoCTX(c)

	ctx := c.Request().Context()
	ctxWithTime, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := h.userSrv.Profile(ctxWithTime, param.ProfileRequest{UserID: claims.UserID})
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
