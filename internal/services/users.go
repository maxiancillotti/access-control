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

func NewUsersServices(dal dal.UserDAL) UsersServices {
	return &usersInteractor{dal: dal}
}

type UsersServices interface {
	// Creates a new user and returns it with the password raw so the user knows
	// what value to use when authenticating.
	// May return InvalidInputID or Internal error categories.
	Create(userDTO dto.User) (*dto.User, error)

	// Updates password and returns it raw so the user knows what value to use
	// when authenticating.
	// May return InvalidInputID or Internal error categories.
	UpdatePassword(id uint) (string, error)

	// May return InvalidInputID or Internal error categories.
	UpdateEnabledState(userDTO dto.User) error

	// May return InvalidInputID or Internal error categories.
	Delete(id uint) error

	// Not necessary at the moment
	//Retrieve(id uint) (*domain.User, error)

	// May return InvalidInputID or Internal error categories.
	RetrieveByUsername(username string) (*dto.User, error)

	// Internal use in the Service Package.
	retrieveByUsername(username string) (user *domain.User, svcErr *svcerr.ServiceError)

	// Checks if userID already exists. Internal use in the Service Package.
	existsOrErr(id uint) (svcErr *svcerr.ServiceError)
}

type usersInteractor struct {
	dal dal.UserDAL
}

const (
	userPasswordLenght     = 64
	userPasswordSaltLenght = 32

	defaultUserEnabledState = true
)

func (ui *usersInteractor) Create(userDTO dto.User) (*dto.User, error) {

	exists, err := ui.dal.UsernameExists(*userDTO.Username)
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

	randpw, randsalt := ui.getRandPasswordAndSalt()
	hashedPw, err := passwords.HashPw(randpw, randsalt)
	if err != nil {
		return nil, svcerr.New(
			errors.Wrap(err, internal.ErrMsgPwHashingFailed),
			internal.ErrorCategoryInternal,
		)
	}

	user := &domain.User{
		Username:     *userDTO.Username,
		PasswordHash: hashedPw,
		PasswordSalt: randsalt,
		Enabled:      defaultUserEnabledState,
	}

	user.ID, err = ui.dal.Insert(user)
	if err != nil {
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtInsertFailed, "user"),
			internal.ErrorCategoryInternal,
		)
	}

	userDTO.ID = &user.ID
	userDTO.Password = &randpw // setting password raw so the user knows what value to use when authenticating
	return &userDTO, nil
}

func (ui *usersInteractor) UpdatePassword(id uint) (string, error) {

	svcErr := ui.existsOrErr(id)
	if svcErr != nil {
		return "", svcErr
	}

	randpw, randsalt := ui.getRandPasswordAndSalt()
	hashedPw, err := passwords.HashPw(randpw, randsalt)
	if err != nil {
		return "", svcerr.New(
			errors.Wrap(err, internal.ErrMsgPwHashingFailed),
			internal.ErrorCategoryInternal,
		)
	}

	user := domain.User{
		ID:           id,
		PasswordHash: hashedPw,
		PasswordSalt: randsalt,
	}

	err = ui.dal.UpdatePassword(&user)
	if err != nil {
		return "", svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtUpdateFailed, "user's password"),
			internal.ErrorCategoryInternal,
		)
	}
	// returning password raw so the user knows what value to use when authenticating
	return randpw, nil
}

// Returns a random generated password and salt
func (ui *usersInteractor) getRandPasswordAndSalt() (string, string) {
	randpw := passwords.RandASCIIString(userPasswordLenght)
	randsalt := passwords.RandASCIIString(userPasswordSaltLenght)
	return randpw, randsalt
}

func (ui *usersInteractor) UpdateEnabledState(userDTO dto.User) error {

	svcErr := ui.existsOrErr(*userDTO.ID)
	if svcErr != nil {
		return svcErr
	}

	user := domain.User{ID: *userDTO.ID, Enabled: *userDTO.EnabledState}

	err := ui.dal.UpdateEnabledState(&user)
	if err != nil {
		return svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtUpdateFailed, "user's enabled state"),
			internal.ErrorCategoryInternal,
		)
	}
	return nil
}

func (ui *usersInteractor) Delete(id uint) error {

	svcErr := ui.existsOrErr(id)
	if svcErr != nil {
		return svcErr
	}

	err := ui.dal.Delete(id)
	if err != nil {
		return svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtDeleteFailed, "user"),
			internal.ErrorCategoryInternal,
		)
	}
	return nil
}

func (ui *usersInteractor) existsOrErr(id uint) (svcErr *svcerr.ServiceError) {
	exists, err := ui.dal.Exists(id)
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtFailedToCheckIfExists, "user"),
			internal.ErrorCategoryInternal,
		)
	} else if !exists {
		svcErr = svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "user"),
			internal.ErrorCategoryInvalidInputID,
		)
	}
	return
}

/*
func (ui *usersInteractor) Retrieve(id uint) (*dto.User, error) {
	user, err := ui.dal.Select(id)
	if err != nil {
		if errors.Is(err, ErrorDataAccessEmptyResult) {
			return nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtDoesNotExist, "user"),
				internal.ErrorCategoryInvalidInputID,
			)
		}
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "user"),
			internal.ErrorCategoryInternal,
		)
	}
	return &dto.User{
		ID:       &user.ID,
		Username: &user.Username,
	}, nil
}
*/

func (ui *usersInteractor) RetrieveByUsername(username string) (*dto.User, error) {

	user, err := ui.retrieveByUsername(username)
	if err != nil {
		return nil, err
	}

	return &dto.User{
		ID:       &user.ID,
		Username: &user.Username,
	}, nil
}

func (ui *usersInteractor) retrieveByUsername(username string) (user *domain.User, svcErr *svcerr.ServiceError) {
	user, err := ui.dal.SelectByUsername(username)
	if err != nil {
		if errors.Is(err, dal.ErrorDataAccessEmptyResult) {
			user, svcErr = nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtDoesNotExist, "username"),
				internal.ErrorCategoryInvalidInputID,
			)
			return
		}
		user, svcErr = nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "user"),
			internal.ErrorCategoryInternal,
		)
		return
	}
	return
}
