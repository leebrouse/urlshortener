package validator

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator  {
	return &CustomValidator{validator: validator.New(),}
}

func (c *CustomValidator)Validate(i interface{}) error {
	return c.validator.Struct(i)
}