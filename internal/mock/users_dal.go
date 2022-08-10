package mock

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/services/dal"

	"github.com/pkg/errors"
)

type UsersDALMock struct{}

// Inserts a User and returns the created ID
func (*UsersDALMock) Insert(user *domain.User) (uint, error) {

	switch user.Username {
	case "userDoesNotExist":
		return 2, nil
	case "userDoesNotExistsReturnInsertErr":
		return 0, errors.New("error msg")
	}
	panic("incorrect username received on Insert mocked func")
}

// Updates a password for the given user
func (*UsersDALMock) UpdatePassword(user *domain.User) error {

	switch user.ID {
	case 1:
		return nil
	case 4:
		return errors.New("error msg user 4")
	}
	panic("incorrect userID received on UpdatePassword mocked func")
}

// Updates enabled state for the given user
func (*UsersDALMock) UpdateEnabledState(user *domain.User) error {

	switch user.ID {
	case 1:
		return nil
	case 4:
		return errors.New("error msg user 4")
	}
	panic("incorrect userID received on UpdateEnabledState mocked func")
}

// Deletes a User logically by its ID
func (*UsersDALMock) Delete(userID uint) error {

	switch userID {
	case 1:
		return nil
	case 4:
		return errors.New("error msg user 4")
	}
	panic("incorrect userID received on Delete mocked func")
}

// Retrieves a User by its ID
func (*UsersDALMock) Select(userID uint) (*domain.User, error) {
	switch userID {
	case 1:
		return &domain.User{ID: 1, Username: "username1", PasswordHash: "passwordHash", PasswordSalt: "salt"}, nil
	case 2:
		return nil, dal.ErrorDataAccessEmptyResult
	case 3:
		return nil, errors.New("another error msg")
	}
	panic("incorrect userID received on Select mocked func")
}

// Checks if a UserID exists
func (*UsersDALMock) Exists(userID uint) (bool, error) {

	switch userID {
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

	panic("incorrect userID received on Exists mocked func")
}

// Checks if a username exists even if it was soft-deleted
func (*UsersDALMock) UsernameExists(username string) (bool, error) {

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

func (d *UsersDALMock) SelectByUsername(username string) (*domain.User, error) {

	switch username {
	case "userExists": // PW plain text: APIUserPassword
		return &domain.User{
			ID:           1,
			Username:     username,
			PasswordHash: "$2a$10$AA92NZpfyuYlANHoXePlG.GackNctcOiBsA6wCegUAKcTqVkrgRLC",
			PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
		}, nil
	// case "userExistsPermissionsError": // PW plain text: APIUserPassword
	// 	return &domain.User{
	// 		ID:           2,
	// 		Username:     username,
	// 		PasswordHash:     "$2a$10$AA92NZpfyuYlANHoXePlG.GackNctcOiBsA6wCegUAKcTqVkrgRLC",
	// 		PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
	// 	}, nil
	// case "userExistsPermissionsEmpty": // PW plain text: APIUserPassword
	// 	return &domain.User{
	// 		ID:           3,
	// 		Username:     username,
	// 		PasswordHash:     "$2a$10$AA92NZpfyuYlANHoXePlG.GackNctcOiBsA6wCegUAKcTqVkrgRLC",
	// 		PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
	// 	}, nil
	case "userDoesNotExist":
		return nil, dal.ErrorDataAccessEmptyResult
	case "userDALerr":
		return nil, errors.New("DAL error")
	}
	return nil, nil
}

/*
// No need to use it anymore
func (d *UsersDALMock) SelectUserPermissionsByUserID(ID uint) (domain.UserPermissions, error) {

	switch ID {
	case 1: // Success
		usrPerm := make(domain.UserPermissions)
		usrPerm["resource"] = append(usrPerm["resource"], "permission")
		return usrPerm, nil
	case 2: // Empty result from database
		return nil, ErrorDataAccessEmptyResult
	case 3: // Returning empty result without err
		usrPerm := make(domain.UserPermissions)
		//usrPerm["resource"] = append(usrPerm["resource"], "permission")
		return usrPerm, nil
	default: // Unexpected error
		return nil, errors.Wrap(errors.New("sql: error"), "SP Result Scan failed")
	}
}
*/
