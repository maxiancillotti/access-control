package dto

import validator "github.com/go-playground/validator/v10"

type HttpMethod struct {
	ID   *uint   `json:"id,omitempty"`
	Name *string `json:"name" validate:"required,alphanum,max=100"`
}

func (uc *HttpMethod) ValidateFormat() error {
	v := validator.New()
	return v.Struct(uc)
}
