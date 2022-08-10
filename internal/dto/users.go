package dto

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type User struct {
	ID           *uint   `json:"id,omitempty"`
	Username     *string `json:"username,omitempty" validate:"required,alphanum,max=100"`
	Password     *string `json:"password,omitempty"`
	EnabledState *bool   `json:"enabled_state,omitempty"`
}

func (u *User) ValidateFormat() error {
	v := validator.New()
	return v.Struct(u)
}

func (u *User) ValidateEmptyEnabledState() error {
	if u.EnabledState == nil {
		return errors.New("state field cannot be empty")
	}
	return nil
}

type UserCredentials struct {
	Username string `json:"username" validate:"required,alphanum,max=100"`
	Password string `json:"password" validate:"required,ascii,max=64"`
}

func (uc *UserCredentials) ValidateFormat() error {
	v := validator.New()
	return v.Struct(uc)
}

/*
type UserState struct {
	ID          uint `json:"id" validate:"required"`
	EnableState bool `json:"enable_state" validate:"required"`
}

func (uc *UserState) ValidateFormat() error {
	v := validator.New()
	return v.Struct(uc)
}
*/
