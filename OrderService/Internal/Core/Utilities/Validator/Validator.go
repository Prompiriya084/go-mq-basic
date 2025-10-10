package utilities_validator

type Validator interface {
	ValidateStruct(s interface{}) error
}
