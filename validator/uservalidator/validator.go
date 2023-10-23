package uservalidator

import (
	"fmt"
	"game-app/dto"
	"game-app/pkg/errs"
	"game-app/pkg/richerror"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Repository interface {
	UniquenePhonenumber(phoneNumer string) (bool, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) (error, map[string]string) {
	const op = "uservalidator.ValidateRegisterRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 40).
			Error(errs.ErrorMsgNameLengthError)),

		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile("^09[0-9]{9,}$")).Error("Invalid phone number"),
			//* use custom validation
			validation.By(v.checkPhoneNumberUniqueess)),

		validation.Field(&req.Password, validation.Required,
			validation.Match(regexp.MustCompile("^[A-Za-z0-9!@#^%<>()]{4,}$"))),
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

		return richerror.New(op).WithMessage(errs.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).WithErr(err).
			WithMeta(map[string]interface{}{"req": req}), fieldErros
	}

	return nil, nil

}

// --> * custom validation
func (v Validator) checkPhoneNumberUniqueess(value interface{}) error {
	PhoneNumber := value.(string)

	if isUnigue, err := v.repo.UniquenePhonenumber(PhoneNumber); err != nil || !isUnigue {
		if err != nil {
			return err
		}
		if !isUnigue {
			return fmt.Errorf(errs.ErrorMsgPhoneNumberIsNotUnique)
		}
	}

	return nil
}
