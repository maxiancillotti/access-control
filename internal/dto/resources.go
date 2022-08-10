package dto

import validator "github.com/go-playground/validator/v10"

type Resource struct {
	ID   *uint   `json:"id,omitempty"`
	Path *string `json:"path" validate:"required,ascii,max=100"`
}

func (uc *Resource) ValidateFormat() error {
	v := validator.New()
	return v.Struct(uc)
}
