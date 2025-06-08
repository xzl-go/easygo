package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Validate(obj interface{}) error {
	return validate.Struct(obj)
}

func RegisterValidation(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}

func RegisterCustomValidation(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}
