package matchinghandler

import (
	"game-app/param"
	"game-app/pkg/claims"
	"game-app/pkg/httpmsg"
	"game-app/service/authservice"
	"game-app/service/matchingservice"
	"game-app/service/presenceservice"
	"game-app/validator/matchingvalidator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	authConfig     authservice.Config
	authSrv        authservice.Service
	matchSrv       matchingservice.Service
	matchValidator matchingvalidator.Validator
	presenceSvc    presenceservice.Service
}

func New(authConfig authservice.Config, authSrv authservice.Service,
	matchSrv matchingservice.Service, matchValidator matchingvalidator.Validator,
	presenceSvc presenceservice.Service) Handler {
	return Handler{
		authConfig:     authConfig,
		authSrv:        authSrv,
		matchSrv:       matchSrv,
		matchValidator: matchValidator,
		presenceSvc:    presenceSvc,
	}
}

func (h Handler) addToWaitingList(c echo.Context) error {

	var req param.AddToWaitingListRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims := claims.GetClaimFromEchoCTX(c)
	req.UserID = claims.UserID

	if fieldErrors, err := h.matchValidator.ValidateWaitList(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	resp, err := h.matchSrv.AddToWaitingList(req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)

}
