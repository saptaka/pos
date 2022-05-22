package handler

import (
	"github.com/go-playground/validator"
	"github.com/saptaka/pos/model"
)

func cashierValidation(processType string) *validator.Validate {

	customValidator := validator.New()
	if processType == model.CREATE {
		customValidator.RegisterStructValidation(validateCashierCreation, model.Cashier{})
	} else {
		customValidator.RegisterStructValidation(validateCashierUpdate, model.Cashier{})
	}

	return customValidator
}

func validateCashierCreation(sl validator.StructLevel) {

	cashier := sl.Current().Interface().(model.Cashier)
	if len(cashier.Name) == 0 {
		sl.ReportError(cashier.Name, "name", "name", "any.required", " \"name\" is required.")
	}

	if len(cashier.Passcode) == 0 {
		sl.ReportError(cashier.Passcode, "passcode", "passcode", "any.required", " \"passcode\" is required.")
	}

	if len(cashier.Passcode) > 6 {
		sl.ReportError(cashier.Passcode, "passcode", "passcode", "any.required", "max 6 character")
	}
}

func validateCashierUpdate(sl validator.StructLevel) {
	cashier := sl.Current().Interface().(model.Cashier)
	if len(cashier.Name) == 0 && len(cashier.Passcode) == 0 {
		sl.ReportError(cashier.Name, "name", "name", "object.missing", " \"value\" must contain at least one of [name]")
		return
	}
}

func categoryValidation(processType string) *validator.Validate {

	customValidator := validator.New()
	if processType == model.CREATE {
		customValidator.RegisterStructValidation(validateCategoryCreation, model.Category{})
	} else {
		customValidator.RegisterStructValidation(validateCategoryUpdate, model.Category{})
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
	category := sl.Current().Interface().(model.Category)
	if len(category.Name) == 0 {
		sl.ReportError(category.Name, "name", "name", "object.missing", " \"value\" must contain at least one of [name]")
		return
	}
}

func paymentValidation(processType string) *validator.Validate {

	customValidator := validator.New()
	if processType == model.CREATE {
		customValidator.RegisterStructValidation(validatePaymentCreation, model.Payment{})
	} else {
		customValidator.RegisterStructValidation(validatePaymentUpdate, model.Payment{})
	}

	return customValidator
}

func validatePaymentCreation(sl validator.StructLevel) {

	payment := sl.Current().Interface().(model.Payment)
	if len(payment.Name) == 0 {
		sl.ReportError(payment.Name, "name", "name", "any.required", " \"name\" is required.")

	}

	if len(payment.Type) == 0 {
		sl.ReportError(payment.Type, "type", "type", "any.required", " \"type\" is required.")

	}

	if !model.PaymentType[payment.Type] {
		sl.ReportError(payment.Type, "type", "type", "any.required", " \"type\" is CASH,E-WALLET,EDC")
	}

}

func validatePaymentUpdate(sl validator.StructLevel) {
	payment := sl.Current().Interface().(model.Payment)
	if len(payment.Name) == 0 && len(payment.Type) == 0 && len(payment.Logo) == 0 {
		sl.ReportError(payment.Name, "name,logo,type", "name", "object.missing", "\"value\" must contain at least one of [name, logo, type]")
		return
	}
}

func productValidation(processType string) *validator.Validate {
	customValidator := validator.New()
	if processType == model.CREATE {
		customValidator.RegisterStructValidation(validateProductCreation, model.ProductCreateRequest{})
	} else {
		customValidator.RegisterStructValidation(validateProductUpdate, model.Product{})
	}

	return customValidator
}

func validateProductCreation(sl validator.StructLevel) {

	product := sl.Current().Interface().(model.ProductCreateRequest)
	if len(product.Name) == 0 {
		sl.ReportError(product.Name, "name", "name", "any.required", " \"name\" is required.")
	}

	if product.CategoryId == nil {
		sl.ReportError(product.CategoryId, "categoryId", "categoryId", "any.required", " \"categoryId\" is required.")
	}

	if product.Stock == 0 {
		sl.ReportError(product.Stock, "stock", "stock", "any.required", " \"stock\" is required.")
	}

	if product.Price == 0 {
		sl.ReportError(product.Price, "price", "price", "any.required", " \"price\" is required.")
	}

	if len(product.Image) == 0 {
		sl.ReportError(product.Image, "image", "image", "any.required", " \"image\" is required.")
	}

}

func validateProductUpdate(sl validator.StructLevel) {
	payment := sl.Current().Interface().(model.Product)
	if len(payment.Name) == 0 &&
		payment.CategoryId == nil &&
		payment.Price == 0 &&
		payment.Stock == 0 &&
		len(payment.Image) == 0 {

		sl.ReportError(payment.Name, "categoryId,name,image,price,stock", "name", "object.missing", "\"value\" must contain at least one of [categoryId,name,image,price,stock]")
		return
	}
}

func orderValidation(processType string) *validator.Validate {
	customValidator := validator.New()
	if processType == model.CREATE {
		customValidator.RegisterStructValidation(validateOrderCreation, model.AddOrderRequest{})
	} else {
		customValidator.RegisterStructValidation(validateSubOrder, model.OrderedProducts{})
	}

	return customValidator
}

func validateOrderCreation(sl validator.StructLevel) {

	order := sl.Current().Interface().(model.AddOrderRequest)
	if order.TotalPaid == 0 {
		sl.ReportError(order.TotalPaid, "totalPaid", "totalPaid", "any.required", " \"totalPaid\" is required.")
	}

	if order.PaymentID == 0 {
		sl.ReportError(order.PaymentID, "paymentId", "paymentId", "any.required", " \"paymentId\" is required.")
	}

	if len(order.OrderedProduct) == 0 || order.OrderedProduct == nil {
		sl.ReportError(order.OrderedProduct, "products", "products", "any.required", " \"products\" is required.")
	}

}

func validateSubOrder(sl validator.StructLevel) {
	order := sl.Current().Interface().(model.OrderedProducts)
	if len(order.Products) == 0 || order.Products == nil {
		sl.ReportError(order, "array", "array", "array.base", "\"value\" must be an array")
		return
	}
}
