package services

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/maxiancillotti/passwords"

	"github.com/pkg/errors"
)

func NewAdminsServices(dal dal.AdminsDAL) AdminsServices {
	return &adminsInteractor{dal: dal}
}

type AdminsServices interface {
	// Creates a new admin and returns it with the password raw so the wielder knows
	// what value to use when authenticating.
	// May return InvalidInputID or Internal error categories.
	Create(adminDTO dto.Admin) (*dto.Admin, error)

	// Updates password and returns it raw so the wielder knows what value to use
	// when authenticating.
	// May return InvalidInputID or Internal error categories.
	UpdatePassword(id uint) (string, error)

	// May return InvalidInputID or Internal error categories.
	UpdateEnabledState(adminDTO dto.Admin) error

	// May return InvalidInputID or Internal error categories.
	Delete(id uint) error

	// Not necessary at the moment
	//Retrieve(id uint) (*domain.Admin, error)

	// May return InvalidInputID or Internal error categories.
	RetrieveByUsername(username string) (*dto.Admin, error)

	// Internal use in the Service Package.
	retrieveByUsername(username string) (admin *domain.Admin, svcErr *svcerr.ServiceError)

	// Checks if adminID already exists. Internal use in the Service Package.
	existsOrErr(id uint) (svcErr *svcerr.ServiceError)
}

type adminsInteractor struct {
	dal dal.AdminsDAL
}

const (
	adminPasswordLenght     = 64
	adminPasswordSaltLenght = 32

	defaultAdminEnabledState = true
)

func (ai *adminsInteractor) Create(adminDTO dto.Admin) (*dto.Admin, error) {

	exists, err := ai.dal.UsernameExists(*adminDTO.Username)
	if err != nil {
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtFailedToCheckIfExists, "username"),
			internal.ErrorCategoryInternal,
		)
	} else if exists {
		return nil, svcerr.New(
			errors.Errorf(internal.ErrMsgFmtAlreadyExists, "username"),
			internal.ErrorCategoryInvalidInputID,
		)
	}

	randpw, randsalt := ai.getRandPasswordAndSalt()
	hashedPw, err := passwords.HashPw(randpw, randsalt)
	if err != nil {
		return nil, svcerr.New(
			errors.Wrap(err, internal.ErrMsgPwHashingFailed),
			internal.ErrorCategoryInternal,
		)
	}

	admin := &domain.Admin{
		Username:     *adminDTO.Username,
		PasswordHash: hashedPw,
		PasswordSalt: randsalt,
		Enabled:      defaultAdminEnabledState,
	}

	admin.ID, err = ai.dal.Insert(admin)
	if err != nil {
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtInsertFailed, "admin"),
			internal.ErrorCategoryInternal,
		)
	}

	adminDTO.ID = &admin.ID
	adminDTO.Password = &randpw // setting password raw so the wielder knows what value to use when authenticating
	return &adminDTO, nil
}

func (ai *adminsInteractor) UpdatePassword(id uint) (string, error) {

	svcErr := ai.existsOrErr(id)
	if svcErr != nil {
		return "", svcErr
	}

	randpw, randsalt := ai.getRandPasswordAndSalt()
	hashedPw, err := passwords.HashPw(randpw, randsalt)
	if err != nil {
		return "", svcerr.New(
			errors.Wrap(err, internal.ErrMsgPwHashingFailed),
			internal.ErrorCategoryInternal,
		)
	}

	admin := domain.Admin{
		ID:           id,
		PasswordHash: hashedPw,
		PasswordSalt: randsalt,
	}

	err = ai.dal.UpdatePassword(&admin)
	if err != nil {
		return "", svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtUpdateFailed, "admin's password"),
			internal.ErrorCategoryInternal,
		)
	}
	// returning password raw so the wielder knows what value to use when authenticating
	return randpw, nil
}

// Returns a random generated password and salt
func (ai *adminsInteractor) getRandPasswordAndSalt() (string, string) {
	randpw := passwords.RandASCIIString(adminPasswordLenght)
	randsalt := passwords.RandASCIIString(adminPasswordSaltLenght)

	//fmt.Printf("randpw: '%s'\n", randpw)
	//fmt.Printf("randsalt: '%s'\n", randsalt)

	return randpw, randsalt
}

func (ai *adminsInteractor) UpdateEnabledState(adminDTO dto.Admin) error {

	svcErr := ai.existsOrErr(*adminDTO.ID)
	if svcErr != nil {
		return svcErr
	}

	admin := domain.Admin{ID: *adminDTO.ID, Enabled: *adminDTO.EnabledState}

	err := ai.dal.UpdateEnabledState(&admin)
	if err != nil {
		return svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtUpdateFailed, "admin's enabled state"),
			internal.ErrorCategoryInternal,
		)
	}
	return nil
}

func (ai *adminsInteractor) Delete(id uint) error {

	svcErr := ai.existsOrErr(id)
	if svcErr != nil {
		return svcErr
	}

	err := ai.dal.Delete(id)
	if err != nil {
		return svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtDeleteFailed, "admin"),
			internal.ErrorCategoryInternal,
		)
	}
	return nil
}

func (ai *adminsInteractor) existsOrErr(id uint) (svcErr *svcerr.ServiceError) {
	exists, err := ai.dal.Exists(id)
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtFailedToCheckIfExists, "admin"),
			internal.ErrorCategoryInternal,
		)
	} else if !exists {
		svcErr = svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "admin"),
			internal.ErrorCategoryInvalidInputID,
		)
	}
	return
}

/*
func (ai *adminsInteractor) Retrieve(id uint) (*dto.admin, error) {
	admin, err := ai.dal.Select(id)
	if err != nil {
		if errors.Is(err, ErrorDataAccessEmptyResult) {
			return nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtDoesNotExist, "admin"),
				internal.ErrorCategoryInvalidInputID,
			)
		}
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "admin"),
			internal.ErrorCategoryInternal,
		)
	}
	return &dto.Admin{
		ID:       &admin.ID,
		Username: &admin.Username,
	}, nil
}
*/

func (ai *adminsInteractor) RetrieveByUsername(username string) (*dto.Admin, error) {

	admin, err := ai.retrieveByUsername(username)
	if err != nil {
		return nil, err
	}

	return &dto.Admin{
		ID:       &admin.ID,
		Username: &admin.Username,
	}, nil
}

func (ai *adminsInteractor) retrieveByUsername(username string) (admin *domain.Admin, svcErr *svcerr.ServiceError) {
	admin, err := ai.dal.SelectByUsername(username)
	if err != nil {
		if errors.Is(err, dal.ErrorDataAccessEmptyResult) {
			admin, svcErr = nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtDoesNotExist, "username"),
				internal.ErrorCategoryInvalidInputID,
			)
			return
		}
		admin, svcErr = nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "admin"),
			internal.ErrorCategoryInternal,
		)
		return
	}
	return
}
