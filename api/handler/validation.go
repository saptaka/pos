package handler

import (
	"github.com/go-playground/validator"
	"github.com/saptaka/pos/model"
)

func cashierValidation(validation *validator.Validate,
	processType string) *validator.Validate {

	customValidator := validation
	if processType == model.CREATE {
		customValidator.RegisterStructValidation(validateCashierCreation, model.Cashier{})
	} else {
		customValidator.RegisterStructValidation(validateCashierUpdate, model.Cashier{})
	}

	return customValidator
}

func validateCashierCreation(sl validator.StructLevel) {

	cashier := sl.Current().Interface().(model.Cashier)
	if len(cashier.Name) == 0 && len(cashier.Passcode) == 0 {
		sl.ReportError(cashier.Name, "name", "name", "any.required", " \"name\" is required.")
		sl.ReportError(cashier.Passcode, "passcode", "passcode", "any.required", " \"passcode\" is required.")
		return
	}

	if len(cashier.Passcode) > 6 {
		sl.ReportError(cashier.Passcode, "passcode", "passcode", "any.required", "max 6 character")
		return
	}
}

func validateCashierUpdate(sl validator.StructLevel) {
	user := sl.Current().Interface().(model.Cashier)
	if len(user.Name) == 0 && len(user.Passcode) == 0 {
		sl.ReportError(user.Name, "name", "name", "object.missing", " \"value\" must contain at least one of [name]")
		return
	}
}

func categoryValidation(validation *validator.Validate,
	processType string) *validator.Validate {

	customValidator := validation
	if processType == model.CREATE {
		customValidator.RegisterStructValidation(validateCategoryCreation, model.Cashier{})
	} else {
		customValidator.RegisterStructValidation(validateCategoryUpdate, model.Cashier{})
	}

	return customValidator
}

func validateCategoryCreation(sl validator.StructLevel) {

	category := sl.Current().Interface().(model.Category)
	if len(category.Name) == 0 {
		sl.ReportError(category.Name, "name", "name", "any.required", " \"name\" is required.")
		return
	}

}

func validateCategoryUpdate(sl validator.StructLevel) {
	user := sl.Current().Interface().(model.Cashier)
	if len(user.Name) == 0 && len(user.Passcode) == 0 {
		sl.ReportError(user.Name, "name", "name", "object.missing", " \"value\" must contain at least one of [name]")
		return
	}
}
