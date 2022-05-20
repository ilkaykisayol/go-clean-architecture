package validator

import (
	"github.com/go-playground/validator/v10"
)

type IValidator interface {
	ValidateStruct(s interface{}) error
}

type Validator struct {
	validate *validator.Validate
}

// New
// Returns a new Validator.
func New() IValidator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}
