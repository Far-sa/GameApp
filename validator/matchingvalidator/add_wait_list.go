package matchingvalidator

import (
	"fmt"
	"game-app/entity"
	"game-app/param"
	"game-app/pkg/errs"
	"game-app/pkg/richerror"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateWaitList(req param.AddToWaitingListRequest) (map[string]string, error) {
	const op = "matchingvalidator.AddToWaitingList"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Category,
			validation.Required,
			validation.By(v.isCategoryValid)),
	); err != nil {

		fieldErrors := make(map[string]string)

		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errs.ErrorMsgInvalidInput).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"req": req})
	}

	return nil, nil
}

func (c Validator) isCategoryValid(value interface{}) error {
	category := value.(entity.Category)

	if !category.IsValid() {
		return fmt.Errorf(errs.ErrorMsgInvalidCategory)
	}

	return nil

}
