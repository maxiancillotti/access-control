package services

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
)

type usersAuthInteractorMock struct{}

// Returns UserID > 0 and error == nil when succesful
func (us *usersAuthInteractorMock) validateUserCredentials(userCredentials *dto.UserCredentials) (userID uint, svcErr *svcerr.ServiceError) {

	switch userCredentials.Username {
	case "userExists":
		userID = 1
		return
	default:
		svcErr = svcerr.New(
			errors.New("internalErr"),
			internal.ErrorCategoryInvalidCredentials,
		)
	}
	return
}

func (d *usersAuthInteractorMock) getUserPermissionsByUserID(userID uint) (usrPerm domain.UserPermissions, svcErr *svcerr.ServiceError) {

	switch userID {
	case 1, 3: // 1 Success - 3 Success but next step on caller fails
		usrPerm = make(domain.UserPermissions)

		userRESTPermissions := make(domain.RESTPermissionsPathsMethods)
		userRESTPermissions["/customers"] = append(userRESTPermissions["/customers"], "POST")

		// Add new permissions categories when needed
		usrPerm[domain.RESTpermissionCategory] = userRESTPermissions

		return

	default:
		usrPerm, svcErr = nil, svcerr.New(
			errors.Wrap(
				errors.New("sql: error"),
				internal.ErrMsgRetrievingPermissions,
			),
			internal.ErrorCategoryInternal,
		)
		return
	}
}

// Fully functional.
// Production copy.
// Can be mocked though but it wasn't necessary to achieve the needed result.
func (d *usersAuthInteractorMock) validatePermissions(perm interface{}, resourceRequested string, methodRequested string, permCat domain.PermissionCategory) (svcErr *svcerr.ServiceError) {

	permStr := perm.(string)

	switch permStr {
	case "permissionsOK":
		return
	default:
		svcErr = svcerr.New(
			errors.Wrap(
				errors.New("sql: error"),
				internal.ErrMsgRetrievingPermissions,
			),
			internal.ErrorCategoryInternal,
		)
	}
	return
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
