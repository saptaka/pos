package handler

import (
	"github.com/go-playground/validator"
	"github.com/saptaka/pos/model"
)

func generateCashierCreateValidation(validation *validator.Validate) *validator.Validate {

	customValidator := validation
	customValidator.RegisterStructValidation(validateCreateCashier, model.Cashier{})

	return customValidator
}

func validateCreateCashier(sl validator.StructLevel) {

	user := sl.Current().Interface().(model.Cashier)

	if len(user.Name) == 0 && len(user.Passcode) == 0 {
		sl.ReportError(user.Name, "name", "name", "any.required", " \"name\" is required.")
		sl.ReportError(user.Passcode, "passcode", "passcode", "any.required", " \"passcode\" is required.")
		return
	}

	if len(user.Passcode) > 6 {
		sl.ReportError(user.Passcode, "passcode", "passcode", "any.required", "max 6 character")
		return
	}

}
