package uservalidator

import (
	"fmt"
	"game-app/dto"
	"game-app/pkg/errs"
	"game-app/pkg/richerror"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateLoginRequest(req dto.LoginRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateLoginRequest"

	if err := validation.ValidateStruct(&req,

		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error("Invalid phone number"),
			//* use custom validation
			validation.By(v.doesPhoneNumberExsit)),
	); err != nil {
		//* helper
		fmt.Println("err valdiator", err)
		fmt.Printf("error type is : %T\n", err)

		fieldErros := make(map[string]string)
		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				if value != nil {
					fieldErros[key] = value.Error()
				}
			}
		}

		return fieldErros, richerror.New(op).WithMessage(errs.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	return nil, nil

}

// --> * custom validation
func (v Validator) doesPhoneNumberExsit(value interface{}) error {
	PhoneNumber := value.(string)

	_, err := v.repo.GetUserByPhoneNumber(PhoneNumber)
	if err != nil {
		return fmt.Errorf(errs.ErrorMsgNotFound)
	}

	return nil
}
