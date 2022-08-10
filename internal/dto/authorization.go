package dto

import (
	validator "github.com/go-playground/validator/v10"
)

type AuthorizationRequest struct {
	Token             string `json:"token" validate:"required,ascii"`
	ResourceRequested string `json:"resource_requested" validate:"required,ascii"`
	MethodRequested   string `json:"method_requested" validate:"required,alpha"`
}

func (vr *AuthorizationRequest) ValidateFormat() error {
	v := validator.New()
	return v.Struct(vr)
}

/*
type ValidationResponse struct {
	Message string `json:"message"`
	IsValid bool   `json:"is_valid"`
}
*/
