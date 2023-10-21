package uservalidator

import (
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

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) error {
	const op = "uservalidator.ValidateRegisterRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 40).
			Error(errs.ErrorMsgNameLengthError)),

		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile("^09[0-9]{9,}$"))),

		validation.Field(&req.Password, validation.Required,
			validation.Match(regexp.MustCompile("^[A-Za-z]{4,}$"))),
	); err != nil {
		return richerror.New(op).WithMessage(errs.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	if isUnigue, err := v.repo.UniquenePhonenumber(req.PhoneNumber); err != nil || !isUnigue {
		if err != nil {
			return richerror.New(op).WithErr(err)
		}
		if !isUnigue {
			return richerror.New(op).WithMessage(errs.ErrorMsgPhoneNumberIsNotUnique).
				WithKind(richerror.KindInvalid)
		}
	}

	return nil

}
