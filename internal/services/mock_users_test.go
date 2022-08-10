package services

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
)

type usersInteractorMock struct{}

// Not implemented
func (ui *usersInteractorMock) Create(userDTO dto.User) (*dto.User, error) {
	return nil, nil
}

// Not implemented
func (ui *usersInteractorMock) UpdatePassword(id uint) (string, error) {

	return "", nil
}

// Not implemented
func (ui *usersInteractorMock) UpdateEnabledState(userDTO dto.User) error {
	return nil
}

// Not implemented
func (ui *usersInteractorMock) Delete(id uint) error {
	return nil
}

// Not implemented
func (ui *usersInteractorMock) existsOrErr(id uint) (svcErr *svcerr.ServiceError) {
	return svcErr
}

// Not implemented
/*
func (ui *usersInteractorMock) Retrieve(id uint) (*domain.User, error) {
	return nil, nil
}
*/

func (ui *usersInteractorMock) RetrieveByUsername(username string) (*dto.User, error) {
	return nil, nil
}

func (ui *usersInteractorMock) retrieveByUsername(username string) (user *domain.User, svcErr *svcerr.ServiceError) {

	switch username {
	case "usernameSuccess":
		return &domain.User{
			ID:       1,
			Username: username,
			Enabled:  true,
			// Raw pw: APIUserPassword
			PasswordHash: "$2a$10$pyxU3pYQ5HvoSfxJLIUzZuI4IQUQcnr8F8UIw.urC50k4BBxvy5E6",
			PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
		}, nil
	case "usernameSuccessStateDisabled":
		return &domain.User{
			ID:           2,
			Username:     username,
			Enabled:      false,
			PasswordHash: "",
			PasswordSalt: "",
		}, nil
	case "usernameSuccessInvalidPassword":
		return &domain.User{
			ID:           3,
			Username:     username,
			Enabled:      true,
			PasswordHash: "PWHash",
			PasswordSalt: "salt",
		}, nil
	case "usernameErrInvalidInputID":
		return nil, svcerr.New(
			errors.New("error msg invalid input"),
			internal.ErrorCategoryInvalidInputID,
		)
	case "usernameErrInternal":
		return nil, svcerr.New(
			errors.New("error msg internal"),
			internal.ErrorCategoryInternal,
		)
	}
	panic("incorrect username received on retrieveByUsername mocked func")
}
