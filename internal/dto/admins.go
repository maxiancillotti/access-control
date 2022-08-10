package dto

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Admin struct {
	ID           *uint   `json:"id,omitempty"`
	Username     *string `json:"username,omitempty" validate:"required,alphanum,max=100"`
	Password     *string `json:"password,omitempty"`
	EnabledState *bool   `json:"enabled_state,omitempty"`
}

func (a *Admin) ValidateFormat() error {
	v := validator.New()
	return v.Struct(a)
}

func (a *Admin) ValidateEmptyEnabledState() error {
	if a.EnabledState == nil {
		return errors.New("state field cannot be empty")
	}
	return nil
}

type AdminCredentials struct {
	Username string `json:"username" validate:"required,alphanum,max=100"`
	Password string `json:"password" validate:"required,ascii,max=64"`
}

func (ac *AdminCredentials) ValidateFormat() error {
	v := validator.New()
	return v.Struct(ac)
}
