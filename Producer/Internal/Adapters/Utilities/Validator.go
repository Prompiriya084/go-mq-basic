package adapters_utilities

import (
	ports_utilities "github.com/Prompiriya084/go-mq/Producer/Internal/Core/Utilities"
	"github.com/go-playground/validator/v10"
)

type validatorImpl struct {
	validate *validator.Validate
}

func NewValidator() ports_utilities.Validator {
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
