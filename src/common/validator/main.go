package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidations registers all custom validation tags
func RegisterCustomValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("strongpassword", StrongPassword)
		// Example: add more in future
		// v.RegisterValidation("customemail", CustomEmailValidator)
	}
}

// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
// 		v.RegisterValidation("strongpassword", regaxValidation.StrongPassword)
// 	}
