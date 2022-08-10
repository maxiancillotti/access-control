package dal

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
)

type UsersRESTPermissionsDAL interface {
	Insert(*domain.UserRESTPermission) error
	Delete(*domain.UserRESTPermission) error

	Exists(*domain.UserRESTPermission) (bool, error)

	// When err == nil, exists will return one of this values:
	// 1 == All relationships exist.
	// -1 == UserID does not exist.
	// -2 == ResourceID does not exist.
	// -3 == HttpMethodID does not exist.
	RelationshipsExists(*domain.UserRESTPermission) (exists int, err error)

	SelectAllByUserID(uint) ([]dto.PermissionsIDs, error)
	SelectAllWithDescriptionsByUserID(userID uint) ([]dto.PermissionsWithDescriptions, error)

	SelectAllPathMethodsByUserID(uint) (domain.RESTPermissionsPathsMethods, error)
}
