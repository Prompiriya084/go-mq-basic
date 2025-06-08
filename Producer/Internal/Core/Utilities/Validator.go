package ports_utilities

type Validator interface {
	ValidateStruct(s interface{}) error
}
