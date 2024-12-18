package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Validator instance global
var Validate *validator.Validate

// Inisialisasi validator dan custom validation
func InitValidator() {
	Validate = validator.New()

	// Daftarkan custom validation
	_ = Validate.RegisterValidation("alphaSpace", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return regexp.MustCompile(`^[a-zA-Z ]+$`).MatchString(value)
	})

	_ = Validate.RegisterValidation("alphanumericSpace", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return regexp.MustCompile(`^[a-zA-Z0-9 ]+$`).MatchString(value)
	})
}
