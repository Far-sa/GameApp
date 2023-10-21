package httpmsg

import (
	"game-app/pkg/errs"
	"game-app/pkg/richerror"
	"net/http"
)

// * mapper fn to transfer errors from Repo layer
func Error(err error) (message string, code int) {
	switch err.(type) {
	case richerror.RichError:
		re := err.(richerror.RichError)
		msg := re.Message()

		//* should not expose unexported error messages
		code := mapKindToHttpStatsusCode(re.Kind())
		if code >= 500 {
			msg = errs.ErrorMsgSomethingWrong
		}
		return msg, code
	default:
		return err.Error(), http.StatusBadRequest
	}
}

func mapKindToHttpStatsusCode(kind richerror.Kind) int {
	switch kind {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindForbidden:
		return http.StatusForbidden
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
