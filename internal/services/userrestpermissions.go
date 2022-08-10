package services

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
)

func NewUsersRESTPermissionsServices(dal dal.UsersRESTPermissionsDAL, userSvc UsersServices) UsersRESTPermissionsServices {
	return &usersRESTPermissionsInteractor{dal: dal, userSvc: userSvc}
}

type UsersRESTPermissionsServices interface {

	/************************************************/
	// Domain methods

	// May return InvalidInputID or Internal error categories.
	Create(pDTO dto.UserRESTPermission) error
	// May return InvalidInputID or Internal error categories.
	Delete(pDTO dto.UserRESTPermission) error

	/************************************************/
	// Intersection methods

	// May return EmptyResult or Internal error categories.
	RetrieveAllByUserID(userID uint) (*dto.UserRESTPermissionsCollection, error)
	// May return EmptyResult or Internal error categories.
	RetrieveAllWithDescriptionsByUserID(userID uint) (*dto.UserRESTPermissionsDescriptionsCollection, error)

	// Intended for internal calling in service layer.
	// May return EmptyResult or Internal error categories.
	retrieveAllPathMethodsByUserID(userID uint) (domain.RESTPermissionsPathsMethods, error)
}

type usersRESTPermissionsInteractor struct {
	dal     dal.UsersRESTPermissionsDAL
	userSvc UsersServices
}

func (pi *usersRESTPermissionsInteractor) Create(pDTO dto.UserRESTPermission) error {

	permission := domain.UserRESTPermission{UserID: *pDTO.UserID, ResourceID: *pDTO.Permission.ResourceID, MethodID: *pDTO.Permission.MethodID}

	exists, err := pi.dal.Exists(&permission)
	if err != nil {
		return svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtFailedToCheckIfExists, "permission"),
			internal.ErrorCategoryInternal,
		)
	} else if exists {
		return svcerr.New(
			errors.Errorf(internal.ErrMsgFmtAlreadyExists, "permission"),
			internal.ErrorCategoryInvalidInputID,
		)
	}

	svcErr := pi.relationshipsExistsOrErr(&permission)
	if svcErr != nil {
		return svcErr
	}

	err = pi.dal.Insert(&permission)
	if err != nil {
		return svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtInsertFailed, "permission"),
			internal.ErrorCategoryInternal,
		)
	}
	return nil
}

func (pi *usersRESTPermissionsInteractor) Delete(pDTO dto.UserRESTPermission) error {

	permission := domain.UserRESTPermission{UserID: *pDTO.UserID, ResourceID: *pDTO.Permission.ResourceID, MethodID: *pDTO.Permission.MethodID}

	svcErr := pi.existsOrErr(&permission)
	if svcErr != nil {
		return svcErr
	}

	err := pi.dal.Delete(&permission)
	if err != nil {
		return svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtDeleteFailed, "permission"),
			internal.ErrorCategoryInternal,
		)
	}
	return nil
}

func (pi *usersRESTPermissionsInteractor) relationshipsExistsOrErr(permission *domain.UserRESTPermission) (svcErr *svcerr.ServiceError) {

	exists, err := pi.dal.RelationshipsExists(permission)
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtFailedToCheckIfExists, "permission's relationships"),
			internal.ErrorCategoryInternal,
		)
		return
	}

	switch exists {

	// Success
	case 1:

	// Invalid Inputs
	case -1:
		svcErr = svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "UserID"),
			internal.ErrorCategoryInvalidInputID,
		)
	case -2:
		svcErr = svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "ResourceID"),
			internal.ErrorCategoryInvalidInputID,
		)
	case -3:
		svcErr = svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "HttpMethodID"),
			internal.ErrorCategoryInvalidInputID,
		)

	// Returned value unexpected
	default:
		svcErr = svcerr.New(
			errors.Wrapf(errors.Errorf("unexpected value received from DAL for 'exists' = %d", exists),
				internal.ErrMsgFmtFailedToCheckIfExists, "permission's relationships"),
			internal.ErrorCategoryInternal,
		)
	}
	return
}

func (pi *usersRESTPermissionsInteractor) existsOrErr(permission *domain.UserRESTPermission) (svcErr *svcerr.ServiceError) {

	exists, err := pi.dal.Exists(permission)
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtFailedToCheckIfExists, "permission"),
			internal.ErrorCategoryInternal,
		)
	} else if !exists {
		svcErr = svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "permission"),
			internal.ErrorCategoryInvalidInputID,
		)
	}
	return
}

func (pi *usersRESTPermissionsInteractor) RetrieveAllByUserID(userID uint) (*dto.UserRESTPermissionsCollection, error) {

	svcErr := pi.userSvc.existsOrErr(userID)
	if svcErr != nil {
		return nil, svcErr
	}

	permissions, err := pi.dal.SelectAllByUserID(userID)
	if err != nil {
		if errors.Is(err, dal.ErrorDataAccessEmptyResult) {
			return nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtEmptyResultForTheGivenID, "permissions", "userID"),
				internal.ErrorCategoryEmptyResult,
			)
		}
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "permissions"),
			internal.ErrorCategoryInternal,
		)
	}

	return &dto.UserRESTPermissionsCollection{UserID: userID, Permissions: permissions}, nil
}

func (pi *usersRESTPermissionsInteractor) RetrieveAllWithDescriptionsByUserID(userID uint) (*dto.UserRESTPermissionsDescriptionsCollection, error) {

	svcErr := pi.userSvc.existsOrErr(userID)
	if svcErr != nil {
		return nil, svcErr
	}

	permissions, err := pi.dal.SelectAllWithDescriptionsByUserID(userID)
	if err != nil {
		if errors.Is(err, dal.ErrorDataAccessEmptyResult) {
			return nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtEmptyResultForTheGivenID, "permissions", "userID"),
				internal.ErrorCategoryEmptyResult,
			)
		}
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "permissions"),
			internal.ErrorCategoryInternal,
		)
	}

	return &dto.UserRESTPermissionsDescriptionsCollection{UserID: userID, Permissions: permissions}, nil
}

func (pi *usersRESTPermissionsInteractor) retrieveAllPathMethodsByUserID(userID uint) (domain.RESTPermissionsPathsMethods, error) {

	svcErr := pi.userSvc.existsOrErr(userID)
	if svcErr != nil {
		return nil, svcErr
	}

	permissions, err := pi.dal.SelectAllPathMethodsByUserID(userID)
	if err != nil {
		if errors.Is(err, dal.ErrorDataAccessEmptyResult) {
			return nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtEmptyResultForTheGivenID, "permissions", "userID"),
				internal.ErrorCategoryEmptyResult,
			)
		}
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "permissions"),
			internal.ErrorCategoryInternal,
		)
	}

	return permissions, nil
}
