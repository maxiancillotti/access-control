package services

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/maxiancillotti/passwords"
	"github.com/pkg/errors"
)

func newUsersAuthServices(userSvc UsersServices, userRESTPermSvc UsersRESTPermissionsServices) userAuthServices {
	return &usersAuthInteractor{userSvc: userSvc, userRESTPermSvc: userRESTPermSvc}
}

type userAuthServices interface {
	validateUserCredentials(userCredentials *dto.UserCredentials) (userID uint, svcErr *svcerr.ServiceError)

	getUserPermissionsByUserID(userID uint) (usrPerm domain.UserPermissions, svcErr *svcerr.ServiceError)

	validatePermissions(perm interface{}, resourceRequested string, methodRequested string, permCat domain.PermissionCategory) (svcErr *svcerr.ServiceError)
}

type usersAuthInteractor struct {
	userSvc         UsersServices
	userRESTPermSvc UsersRESTPermissionsServices
}

// Returns UserID > 0 and error == nil when succesful
func (uai *usersAuthInteractor) validateUserCredentials(userCredentials *dto.UserCredentials) (userID uint, svcErr *svcerr.ServiceError) {

	user, svcErr := uai.userSvc.retrieveByUsername(userCredentials.Username)
	if svcErr != nil {
		if svcErr.Category() == internal.ErrorCategoryInvalidInputID {
			svcErr = svcerr.New(
				errors.Wrap(svcErr.ErrorValue(), internal.ErrMsgInvalidUsername),
				internal.ErrorCategoryInvalidCredentials,
			)
		}
		return
	}

	if !user.Enabled {
		svcErr = svcerr.New(
			errors.New(internal.ErrMsgUserDisabled),
			//internal.ErrorCategoryUserDisabled,
			internal.ErrorCategoryInvalidCredentials,
		)
		return
	}

	err := passwords.ValidateHashPw(userCredentials.Password, user.PasswordSalt, user.PasswordHash)
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrap(err, internal.ErrMsgInvalidPassword),
			internal.ErrorCategoryInvalidCredentials,
		)
		return
	}
	userID = user.ID
	return
}

func (uai *usersAuthInteractor) getUserPermissionsByUserID(userID uint) (usrPerm domain.UserPermissions, svcErr *svcerr.ServiceError) {

	userRESTPermissions, err := uai.userRESTPermSvc.retrieveAllPathMethodsByUserID(userID)
	if err != nil {
		// I'm sort of depending on implementation details here, but take into account that this
		// struct is made to function around the two that are its fields.
		restPermSvcErr, ok := err.(*svcerr.ServiceError)
		if !ok || restPermSvcErr.Category() == internal.ErrorCategoryInvalidInputID {
			svcErr = svcerr.New(
				errors.Wrap(err, internal.ErrMsgRetrievingPermissions),
				internal.ErrorCategoryInternal,
			)
		} else if restPermSvcErr.Category() == internal.ErrorCategoryEmptyResult {
			svcErr = svcerr.New(
				errors.Wrap(err, internal.ErrMsgUserDoesntHaveAnyPermissions),
				internal.ErrorCategoryInternal,
			)
		} else {
			svcErr = restPermSvcErr
		}
		return
	}

	usrPerm = make(domain.UserPermissions, 0)

	// Add new permissions categories when needed
	usrPerm[domain.RESTpermissionCategory] = userRESTPermissions

	return
}

func (uai *usersAuthInteractor) validatePermissions(permRaw interface{}, resourceRequested string, accessRequested string, permCat domain.PermissionCategory) (svcErr *svcerr.ServiceError) {

	// ASSERT
	//userPermissionsCategories, ok := permRaw.(map[domain.PermissionCategory]interface{})
	userPermissionsCategories, ok := permRaw.(map[string]interface{}) // domain.PermissionCategory's underlying type
	if !ok {
		svcErr = svcerr.New(
			errors.New("permission categories assertion failed"),
			internal.ErrorCategorySemanticallyUnprocesable,
		)
		return
	}

	userPermissions, ok := userPermissionsCategories[string(permCat)].(map[string]interface{})
	if !ok {
		svcErr = svcerr.New(
			errors.New("permissions assertion failed"),
			internal.ErrorCategorySemanticallyUnprocesable,
		)
		return
	}

	resourcePermissions := userPermissions[resourceRequested]
	if resourcePermissions == nil {
		svcErr = svcerr.New(
			errors.New("token doesn't have any permissions for this resource"),
			internal.ErrorCategoryNotEnoughPermissions,
		)
		return
	}

	accessPermissions, ok := resourcePermissions.([]interface{})
	if !ok {
		svcErr = svcerr.New(
			errors.New("access types permissions assertion failed"),
			internal.ErrorCategorySemanticallyUnprocesable,
		)
		return
	}

	// CHECK METHODS
	sliceLen := len(accessPermissions)
	for i, access := range accessPermissions {

		accessStr, ok := access.(string)
		if !ok {
			svcErr = svcerr.New(
				errors.New("access type assertion failed"),
				internal.ErrorCategorySemanticallyUnprocesable,
			)
			return
		}

		if accessStr == accessRequested {
			break
		}
		if i+1 == sliceLen {
			svcErr = svcerr.New(
				errors.New("token doesn't have the requested permissions for this resource"),
				internal.ErrorCategoryNotEnoughPermissions,
			)
			return
		}
	}
	return
}
