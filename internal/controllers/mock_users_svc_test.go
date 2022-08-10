package controllers

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
)

type usersServiceMock struct {
	// Implementing this interface embedding it so the compiler can recognize unexported methods
	services.UsersServices
}

func (ui *usersServiceMock) Create(userDTO dto.User) (*dto.User, error) {
	switch *userDTO.Username {

	case "username1":
		uintOne := uint(1)
		passwordStr := "password"

		return &dto.User{
			ID:       &uintOne,
			Username: userDTO.Username,
			Password: &passwordStr,
		}, nil

	case "username0":
		return nil, errors.New("invalid_input")
	case "username2":
		return nil, errors.New("internal_error")
	}
	panic("invalid username received in mocked Create method")
}

func (ui *usersServiceMock) UpdatePassword(id uint) (string, error) {
	switch id {
	case 1:
		return "password", nil
	case 0:
		return "", errors.New("invalid_input")
	case 2:
		return "", errors.New("internal_error")
	}
	panic("invalid userID received in mocked UpdatePassword method")
}

func (ui *usersServiceMock) UpdateEnabledState(userDTO dto.User) error {
	switch *userDTO.ID {
	case 1:
		return nil
	case 0:
		return errors.New("invalid_input")
	case 2:
		return errors.New("internal_error")
	}
	panic("invalid userID received in mocked UpdateEnabledState method")
}

func (ui *usersServiceMock) Delete(id uint) error {
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
func (ui *usersServiceMock) existsOrErr(id uint) (svcErr *svcerr.ServiceError) {
	return nil
}

func (ui *usersServiceMock) RetrieveByUsername(username string) (*dto.User, error) {
	switch username {

	case "username1":
		uintOne := uint(1)
		boolTrue := true

		return &dto.User{
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
func (ui *usersServiceMock) retrieveByUsername(username string) (user *domain.User, svcErr *svcerr.ServiceError) {
	return nil, nil
}
