package services

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"
	"github.com/pkg/errors"
)

type usersRESTPermissionsInteractorMock struct{}

// Not implemented
func (pi *usersRESTPermissionsInteractorMock) Create(pDTO dto.UserRESTPermission) error {
	return nil
}

// Not implemented
func (pi *usersRESTPermissionsInteractorMock) Delete(pDTO dto.UserRESTPermission) error {
	return nil
}

// Not implemented
func (pi *usersRESTPermissionsInteractorMock) RetrieveAllByUserID(userID uint) (*dto.UserRESTPermissionsCollection, error) {
	return nil, nil
}

// Not implemented
func (pi *usersRESTPermissionsInteractorMock) RetrieveAllWithDescriptionsByUserID(userID uint) (*dto.UserRESTPermissionsDescriptionsCollection, error) {
	return nil, nil
}

// Not implemented
func (pi *usersRESTPermissionsInteractorMock) retrieveAllPathMethodsByUserID(userID uint) (domain.RESTPermissionsPathsMethods, error) {
	switch userID {
	case 1:
		return domain.RESTPermissionsPathsMethods{
			"/customers": {"POST"},
		}, nil
	case 2:
		return nil, svcerr.New(
			errors.New("err msg empty result"),
			internal.ErrorCategoryEmptyResult,
		)
	case 3:
		return nil, svcerr.New(
			errors.New("err msg internal"),
			internal.ErrorCategoryInternal,
		)
	case 4:
		return nil, svcerr.New(
			errors.New("err msg invalid input id"),
			internal.ErrorCategoryInvalidInputID,
		)
	}
	panic("incorrect userID received on RetrieveAllPathMethodsByUserID mocked func")
}
