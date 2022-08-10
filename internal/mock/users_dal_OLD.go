package mock

/*
import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/services"

	"github.com/pkg/errors"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////////

type UserDALMock struct{}

func (d *UserDALMock) SelectUserByUsername(username string) (*domain.User, error) {

	switch username {
	case "userExists": // PW plain text: APIUserPassword
		return &domain.User{
			ID:           1,
			Username:     username,
			Password:     "$2a$10$AA92NZpfyuYlANHoXePlG.GackNctcOiBsA6wCegUAKcTqVkrgRLC",
			PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
		}, nil
	case "userExistsPermissionsError": // PW plain text: APIUserPassword
		return &domain.User{
			ID:           2,
			Username:     username,
			Password:     "$2a$10$AA92NZpfyuYlANHoXePlG.GackNctcOiBsA6wCegUAKcTqVkrgRLC",
			PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
		}, nil
	case "userExistsPermissionsEmpty": // PW plain text: APIUserPassword
		return &domain.User{
			ID:           3,
			Username:     username,
			Password:     "$2a$10$AA92NZpfyuYlANHoXePlG.GackNctcOiBsA6wCegUAKcTqVkrgRLC",
			PasswordSalt: "sv8TM3WJcLQ4GXIsCBhUSS0964L4ZA7S",
		}, nil
	case "userDoesNotExist":
		return nil, errors.Wrap(services.ErrorDataAccessEmptyResult, "cannot select User")
	case "userDALerr":
		return nil, errors.Wrap(errors.New("DAL error"), "cannot select User")
	}
	return nil, nil
}

func (d *UserDALMock) SelectUserPermissionsByUserID(ID int) (domain.UserPermissions, error) {

	usrPerm := make(domain.UserPermissions)
	if ID == 1 {
		//usrPerm["resource"] = append(usrPerm["resource"], "permission")
		usrPerm = domain.UserPermissions{
			domain.RESTpermissionCategory: domain.RESTPermissionsPathsMethods{
				"/resource": {"POST"},
			},
		}
		return usrPerm, nil
	}
	if ID == 3 {
		return usrPerm, nil
	}
	return nil, errors.Wrap(errors.New("sql: error"), "SP Result Scan failed")
}
*/
