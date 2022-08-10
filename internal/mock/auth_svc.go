package mock

import (
	"github.com/maxiancillotti/access-control/internal/dto"
)

type AuthServiceMock struct{}

func (s *AuthServiceMock) AuthenticateAdmin(adminCredentials *dto.AdminCredentials) (uint, error) {

	switch adminCredentials.Username {
	case "userSuccess":
		return 1, nil
	case "userInvalidCredentials":
		return 0, errMockInvalidCredentials
	}
	return 0, errMockInternal
}

func (s *AuthServiceMock) Authenticate(userCredentials *dto.UserCredentials) (uint, error) {

	switch userCredentials.Username {
	case "userSuccess":
		return 1, nil
	case "userInvalidCredentials":
		return 0, errMockInvalidCredentials
	}
	return 0, errMockInternal
}

func (s *AuthServiceMock) CreateToken(userID uint) ([]byte, error) {

	switch userID {
	case 1:
		return []byte("header.payload.sign"), nil
		// case 2:
		// 	return nil, errMockInvalidCredentials
	}
	return nil, errMockInternal
}

func (s *AuthServiceMock) Authorize(authorizationData dto.AuthorizationRequest, permissionCategory string) error {

	switch authorizationData.Token {
	case "valid":
		return nil
	case "invalid":
		return errMockInvalidToken
	case "not_enough_permissions":
		return errMockPermissions
	case "unprocessable":
		return errMockUnprocessable
	}
	return errMockInternal
}
