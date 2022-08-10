package mock

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"
	"github.com/pkg/errors"
)

type UsersRESTPermissionsDALMock struct{}

// Inserts a User and returns the created ID
func (d *UsersRESTPermissionsDALMock) Insert(permission *domain.UserRESTPermission) error {

	switch permission.UserID {
	case 2:
		return nil
	case 5:
		return errors.New("error insert user 5")
	}
	panic("incorrect permission.UserID received on Insert mocked func")
}

func (d *UsersRESTPermissionsDALMock) Delete(permission *domain.UserRESTPermission) error {

	switch permission.UserID {
	case 1:
		return nil
	case 4:
		return errors.New("delete err")
	}
	panic("incorrect permission.UserID received on Delete mocked func")
}

func (d *UsersRESTPermissionsDALMock) SelectAllByUserID(userID uint) ([]dto.PermissionsIDs, error) {

	switch userID {
	case 1:
		return []dto.PermissionsIDs{
			{ResourceID: 1, MethodsIDs: []uint{1, 2, 3}},
		}, nil
	case 4:
		return nil, errors.New("SelectAllByUserID err userID 4")
	case 5:
		return nil, dal.ErrorDataAccessEmptyResult
	}

	panic("incorrect UserID received on SelectAllByUserID mocked func")
}

func (d *UsersRESTPermissionsDALMock) SelectAllWithDescriptionsByUserID(userID uint) ([]dto.PermissionsWithDescriptions, error) {

	switch userID {
	case 1:
		resourceID := uint(10)
		resourcePath := "/customers"
		methodID := uint(15)
		methodName := "POST"

		return []dto.PermissionsWithDescriptions{
			{
				Resource: dto.Resource{
					ID:   &resourceID,
					Path: &resourcePath,
				},
				Methods: []dto.HttpMethod{
					{ID: &methodID,
						Name: &methodName,
					},
				},
			},
		}, nil
	case 4:
		return nil, errors.New("SelectAllWithDescriptionsByUserID err userID 4")
	case 5:
		return nil, dal.ErrorDataAccessEmptyResult
	}

	panic("incorrect UserID received on SelectAllWithDescriptionsByUserID mocked func")
}

func (d *UsersRESTPermissionsDALMock) SelectAllPathMethodsByUserID(userID uint) (domain.RESTPermissionsPathsMethods, error) {

	switch userID {
	case 1:
		return domain.RESTPermissionsPathsMethods{"/customers": {"POST", "GET"}}, nil
	case 4:
		return nil, errors.New("SelectAllPathMethodsByUserID err userID 4")
	case 5:
		return nil, dal.ErrorDataAccessEmptyResult
	}

	panic("incorrect UserID received on SelectAllPathMethodsByUserID mocked func")
}

func (d *UsersRESTPermissionsDALMock) Exists(permission *domain.UserRESTPermission) (bool, error) {

	switch permission.UserID {
	case 1: // Exists
		return true, nil

	case 2: // Does not exists
		return false, nil

	case 3: // Dal err
		return false, errors.New("err user 3")

	case 4: // Exists case to return err later on to the caller
		return true, nil

	case 5: // Does not exists case to return err later on to the caller
		return false, nil

	case 6: // Does not exists case to return err later on to the caller
		return false, nil
	}

	panic("incorrect permission.UserID received on Exists mocked func")
}

func (d *UsersRESTPermissionsDALMock) RelationshipsExists(permission *domain.UserRESTPermission) (int, error) {

	switch permission.UserID {

	case 2: // Relationships Exist
		return 1, nil

	case 5: // Relationships Exist. Case to return err later on to the caller
		return 1, nil

	case 6: // Dal err
		return 0, errors.New("err user 6")

	case 7: // UserID does not exists
		return -1, nil
	case 8: // ResourceID does not exists
		return -2, nil
	case 9: // HttpMethodID does not exists
		return -3, nil

	case 10: // Unexpected value 0 when err == nil
		return 0, nil
	}

	panic("incorrect permission.UserID received on RelationshipsExists mocked func")
}
