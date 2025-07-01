package utilities_validator

import (
	"github.com/go-playground/validator/v10"
)

type validatorImpl struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	return &validatorImpl{
		validate: validator.New(),
	}
}

func (v *validatorImpl) ValidateStruct(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		return err
	}

	return nil
}
