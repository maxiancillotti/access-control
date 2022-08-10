package services

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
)

type adminsInteractorMock struct{}

// Not implemented
func (ai *adminsInteractorMock) Create(userDTO dto.Admin) (*dto.Admin, error) {
	return nil, nil
}

// Not implemented
func (ai *adminsInteractorMock) UpdatePassword(id uint) (string, error) {

	return "", nil
}

// Not implemented
func (ai *adminsInteractorMock) UpdateEnabledState(userDTO dto.Admin) error {
	return nil
}

// Not implemented
func (ai *adminsInteractorMock) Delete(id uint) error {
	return nil
}

// Not implemented
func (ai *adminsInteractorMock) existsOrErr(id uint) (svcErr *svcerr.ServiceError) {
	return svcErr
}

// Not implemented
/*
func (ai *adminsInteractorMock) Retrieve(id uint) (*domain.Admin, error) {
	return nil, nil
}
*/

func (ai *adminsInteractorMock) RetrieveByUsername(username string) (*dto.Admin, error) {
	return nil, nil
}

func (ai *adminsInteractorMock) retrieveByUsername(username string) (user *domain.Admin, svcErr *svcerr.ServiceError) {

	switch username {
	case "usernameSuccess":
		return &domain.Admin{
			ID:       1,
			Username: username,
			Enabled:  true,
			// Raw pw: APIUserPassword
			PasswordHash: "$2a$10$pyxU3pYQ5HvoSfxJLIUzZuI4IQUQcnr8F8UIw.urC50k4BBxvy5E6",
			PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
		}, nil
	case "usernameSuccessStateDisabled":
		return &domain.Admin{
			ID:           2,
			Username:     username,
			Enabled:      false,
			PasswordHash: "",
			PasswordSalt: "",
		}, nil
	case "usernameSuccessInvalidPassword":
		return &domain.Admin{
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
