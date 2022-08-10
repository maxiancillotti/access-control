package mock

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/services/dal"

	"github.com/pkg/errors"
)

type AdminsDALMock struct{}

// Inserts a Admin and returns the created ID
func (*AdminsDALMock) Insert(admin *domain.Admin) (uint, error) {

	switch admin.Username {
	case "userDoesNotExist":
		return 2, nil
	case "userDoesNotExistsReturnInsertErr":
		return 0, errors.New("error msg")
	}
	panic("incorrect username received on Insert mocked func")
}

// Updates a password for the given admin
func (*AdminsDALMock) UpdatePassword(admin *domain.Admin) error {

	switch admin.ID {
	case 1:
		return nil
	case 4:
		return errors.New("error msg admin 4")
	}
	panic("incorrect adminID received on UpdatePassword mocked func")
}

// Updates enabled state for the given admin
func (*AdminsDALMock) UpdateEnabledState(admin *domain.Admin) error {

	switch admin.ID {
	case 1:
		return nil
	case 4:
		return errors.New("error msg admin 4")
	}
	panic("incorrect adminID received on UpdateEnabledState mocked func")
}

// Deletes a Admin logically by its ID
func (*AdminsDALMock) Delete(adminID uint) error {

	switch adminID {
	case 1:
		return nil
	case 4:
		return errors.New("error msg admin 4")
	}
	panic("incorrect adminID received on Delete mocked func")
}

// Retrieves a Admin by its ID
func (*AdminsDALMock) Select(adminID uint) (*domain.Admin, error) {
	switch adminID {
	case 1:
		return &domain.Admin{ID: 1, Username: "username1", PasswordHash: "password", PasswordSalt: "salt"}, nil
	case 2:
		return nil, dal.ErrorDataAccessEmptyResult
	case 3:
		return nil, errors.New("another error msg")
	}
	panic("incorrect adminID received on Select mocked func")
}

// Checks if a adminID exists
func (*AdminsDALMock) Exists(adminID uint) (bool, error) {

	switch adminID {
	case 1:
		return true, nil
	case 2:
		return false, nil
	case 3:
		return false, errors.New("error msg")
	// returns ok now but you can use it to return err later on
	case 4:
		return true, nil
	// returns ok now but you can use it to return err later on
	case 5:
		return true, nil
	}

	panic("incorrect adminID received on Exists mocked func")
}

// Checks if a username exists even if it was soft-deleted
func (*AdminsDALMock) UsernameExists(username string) (bool, error) {

	switch username {
	case "userExists":
		return true, nil
	case "userDoesNotExist":
		return false, nil
	case "userDALerr":
		return false, errors.New("DAL error")
	case "userDoesNotExistsReturnInsertErr":
		return false, nil
	}
	panic("incorrect username received on UsernameExists mocked func")
}

func (d *AdminsDALMock) SelectByUsername(username string) (*domain.Admin, error) {

	switch username {
	case "userExists": // PW plain text: APIUserPassword
		return &domain.Admin{
			ID:           1,
			Username:     username,
			PasswordHash: "$2a$10$AA92NZpfyuYlANHoXePlG.GackNctcOiBsA6wCegUAKcTqVkrgRLC",
			PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
		}, nil
	// case "userExistsPermissionsError": // PW plain text: APIUserPassword
	// 	return &domain.Admin{
	// 		ID:           2,
	// 		Username:     username,
	// 		Password:     "$2a$10$AA92NZpfyuYlANHoXePlG.GackNctcOiBsA6wCegUAKcTqVkrgRLC",
	// 		PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
	// 	}, nil
	// case "userExistsPermissionsEmpty": // PW plain text: APIUserPassword
	// 	return &domain.Admin{
	// 		ID:           3,
	// 		Username:     username,
	// 		Password:     "$2a$10$AA92NZpfyuYlANHoXePlG.GackNctcOiBsA6wCegUAKcTqVkrgRLC",
	// 		PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
	// 	}, nil
	case "userDoesNotExist":
		return nil, dal.ErrorDataAccessEmptyResult
	case "userDALerr":
		return nil, errors.New("DAL error")
	}
	return nil, nil
}
