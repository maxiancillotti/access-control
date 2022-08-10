package services

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
)

func NewAuthServices(userSvc UsersServices, admSvc AdminsServices, userRESTPermSvc UsersRESTPermissionsServices, tokenSvc internal.AuthTokenServices) AuthServices {
	return &authInteractor{
		userAuthSvc:   newUsersAuthServices(userSvc, userRESTPermSvc),
		adminsAuthSvc: newAdminsAuthServices(admSvc),
		tokenSvc:      tokenSvc,
	}
}

// Method callers can check error category using the ServiceErrorChecker interface.
type AuthServices interface {

	// Checks if the given credentials are accepted as proof of identity.
	// When err != nil, returns the UserID associated to the user credentials.
	// Method callers can check error category using the ServiceErrorChecker interface.
	Authenticate(*dto.UserCredentials) (uint, error)

	// Special method for Admin Users.
	// Checks if the given credentials are accepted as proof of identity.
	// When err != nil, returns the AdminID associated to the admin credentials.
	// Method callers can check error category using the ServiceErrorChecker interface.
	AuthenticateAdmin(*dto.AdminCredentials) (uint, error)

	// Creates and returns an authorization token with the user permissions,
	// to use as proof of identity, during the authorization of a future
	// request from the holder.
	// Method callers can check error category using the ServiceErrorChecker interface.
	CreateToken(userID uint) ([]byte, error)

	// Verifies if an authorization token is valid and has the necessary
	// permissions to access a resource.
	// See accepted values for permissionCategory en domain package.
	// Method callers can check error category using the ServiceErrorChecker interface.
	Authorize(authorizationData dto.AuthorizationRequest, permissionCategory string) error
}

type authInteractor struct {
	userAuthSvc   userAuthServices
	adminsAuthSvc adminsAuthServices
	tokenSvc      internal.AuthTokenServices
}

func (s *authInteractor) Authenticate(userCredentials *dto.UserCredentials) (uint, error) {

	userID, err := s.userAuthSvc.validateUserCredentials(userCredentials)
	if err != nil {
		return 0, svcerr.New(
			errors.Wrap(err.ErrorValue(), internal.ErrMsgValidatingCredencials.Error()),
			err.Category(),
		)
	}
	return userID, nil
}

func (s *authInteractor) AuthenticateAdmin(admCredentials *dto.AdminCredentials) (uint, error) {

	admID, err := s.adminsAuthSvc.validateAdminCredentials(admCredentials)
	if err != nil {
		return 0, svcerr.New(
			errors.Wrap(err.ErrorValue(), internal.ErrMsgValidatingCredencials.Error()),
			err.Category(),
		)
	}
	return admID, nil
}

func (s *authInteractor) CreateToken(userID uint) ([]byte, error) {

	userPermissions, err := s.userAuthSvc.getUserPermissionsByUserID(userID)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenSvc.GenerateToken(userID, userPermissions)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *authInteractor) Authorize(authorizationData dto.AuthorizationRequest, permissionCategory string) (svcErr error) {

	permissions, errValid := s.tokenSvc.ValidateToken(authorizationData.Token)
	if errValid != nil {
		svcErr = errValid
		return
	}

	errPerm := s.userAuthSvc.validatePermissions(permissions, authorizationData.ResourceRequested, authorizationData.MethodRequested, domain.PermissionCategory(permissionCategory))
	if errPerm != nil {
		svcErr = errPerm
	}
	//svcErr must have nil value and type so the caller can evaluate it as nil (untyped nil)
	return
}
