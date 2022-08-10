package controllers

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
)

type adminsServiceMock struct {
	// Implementing this interface embedding it so the compiler can recognize unexported methods
	services.AdminsServices
}

func (ai *adminsServiceMock) Create(adminDTO dto.Admin) (*dto.Admin, error) {
	switch *adminDTO.Username {

	case "username1":
		uintOne := uint(1)
		passwordStr := "password"

		return &dto.Admin{
			ID:       &uintOne,
			Username: adminDTO.Username,
			Password: &passwordStr,
		}, nil

	case "username0":
		return nil, errors.New("invalid_input")
	case "username2":
		return nil, errors.New("internal_error")
	}
	panic("invalid username received in mocked Create method")
}

func (ai *adminsServiceMock) UpdatePassword(id uint) (string, error) {
	switch id {
	case 1:
		return "password", nil
	case 0:
		return "", errors.New("invalid_input")
	case 2:
		return "", errors.New("internal_error")
	}
	panic("invalid adminID received in mocked UpdatePassword method")
}

func (ai *adminsServiceMock) UpdateEnabledState(adminDTO dto.Admin) error {
	switch *adminDTO.ID {
	case 1:
		return nil
	case 0:
		return errors.New("invalid_input")
	case 2:
		return errors.New("internal_error")
	}
	panic("invalid adminID received in mocked UpdateEnabledState method")
}

func (ai *adminsServiceMock) Delete(id uint) error {
	switch id {
	case 1:
		return nil
	case 0:
		return errors.New("invalid_input")
	case 2:
		return errors.New("internal_error")
	}
	panic("invalid id received in mocked Delete method")
}

// Not implemented. Private method.
func (ai *adminsServiceMock) existsOrErr(id uint) (svcErr *svcerr.ServiceError) {
	return nil
}

func (ai *adminsServiceMock) RetrieveByUsername(username string) (*dto.Admin, error) {
	switch username {

	case "username1":
		uintOne := uint(1)
		boolTrue := true

		return &dto.Admin{
			ID:           &uintOne,
			Username:     &username,
			EnabledState: &boolTrue,
		}, nil

	case "username0":
		return nil, errors.New("invalid_input")
	case "username2":
		return nil, errors.New("internal_error")
	}
	panic("invalid username received in mocked Create method")
}

// Not implemented. Private method.
func (ai *adminsServiceMock) retrieveByUsername(username string) (admin *domain.Admin, svcErr *svcerr.ServiceError) {
	return nil, nil
}
