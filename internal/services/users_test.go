package services

import (
	"testing"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUsersGetRandPasswordAndSalt(t *testing.T) {

	testUsrSvcs := testUsersServices.(*usersInteractor)
	pw, salt := testUsrSvcs.getRandPasswordAndSalt()

	pwLen := len(pw)
	saltLen := len(salt)

	assert.Equal(t, 64, pwLen)
	assert.Equal(t, 32, saltLen)
}

func TestUsersCreate(t *testing.T) {

	type testCase struct {
		name                      string
		usernameInput             string
		expectedUserOutput        *dto.User
		expectedPasswordLenOutput int
		expectedErrOutput         *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	outputUserIDCase1 := uint(2)
	outputUserUsernameCase1 := "userDoesNotExist"

	table = append(table, testCase{
		name:                      "Success",
		usernameInput:             "userDoesNotExist",
		expectedUserOutput:        &dto.User{ID: &outputUserIDCase1, Username: &outputUserUsernameCase1},
		expectedPasswordLenOutput: 64,
		expectedErrOutput:         nil,
	})

	table = append(table, testCase{
		name:               "Error. Username already exists. Invalid input ID.",
		usernameInput:      "userExists",
		expectedUserOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtAlreadyExists, "username"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:               "Error. Failed to check if username exists. Internal.",
		usernameInput:      "userDALerr",
		expectedUserOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("DAL error"), internal.ErrMsgFmtFailedToCheckIfExists, "username"),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name:               "Error. Insert failed. Internal.",
		usernameInput:      "userDoesNotExistsReturnInsertErr",
		expectedUserOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg"), internal.ErrMsgFmtInsertFailed, "user"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			userDTO, err := testUsersServices.Create(dto.User{Username: &test.usernameInput})

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.NotNil(t, userDTO)

				assert.Equal(t, test.expectedUserOutput.ID, userDTO.ID)
				assert.Equal(t, test.expectedUserOutput.Username, userDTO.Username)

				pwLen := len(*userDTO.Password)
				assert.Equal(t, test.expectedPasswordLenOutput, pwLen)
			} else {
				assert.NotNil(t, err)

				svcErr, ok := err.(*svcerr.ServiceError)
				assert.True(t, ok)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}
		})
	}
}

func TestUsersUpdatePassword(t *testing.T) {

	type testCase struct {
		name                      string
		input                     uint
		expectedPasswordLenOutput int
		expectedErrOutput         *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                      "Success",
		input:                     1,
		expectedPasswordLenOutput: 64,
		expectedErrOutput:         nil,
	})

	table = append(table, testCase{
		name:                      "Error from existsOrErr. User does not exist",
		input:                     2,
		expectedPasswordLenOutput: 0,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "user"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Internal. Update failed",
		input:                     4,
		expectedPasswordLenOutput: 0,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg user 4"), internal.ErrMsgFmtUpdateFailed, "user's password"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			pw, err := testUsersServices.UpdatePassword(test.input)

			pwLen := len(pw)
			assert.Equal(t, test.expectedPasswordLenOutput, pwLen)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)

				svcErr, ok := err.(*svcerr.ServiceError)
				assert.True(t, ok)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}
		})
	}
}

func TestUsersExistsOrErr(t *testing.T) {

	type testCase struct {
		name           string
		input          uint
		expectedOutput *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:           "Success",
		input:          1,
		expectedOutput: nil,
	})

	table = append(table, testCase{
		name:  "Error. User does not exist.",
		input: 2,
		expectedOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "user"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:  "Error. Failed to check if user exist",
		input: 3,
		expectedOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg"), internal.ErrMsgFmtFailedToCheckIfExists, "user"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			testUsrSvcs := testUsersServices.(*usersInteractor)
			svcErr := testUsrSvcs.existsOrErr(test.input)

			if test.expectedOutput == nil {
				assert.Nil(t, svcErr)
			} else {
				assert.NotNil(t, svcErr)

				assert.Equal(t, test.expectedOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedOutput.Category(), svcErr.Category())
			}
		})
	}
}

func TestUsersUpdateEnabledState(t *testing.T) {

	type testCase struct {
		name              string
		inputUserID       uint
		inputEnabledState bool
		expectedOutput    *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:              "Success",
		inputUserID:       1,
		inputEnabledState: false,
		expectedOutput:    nil,
	})

	table = append(table, testCase{
		name:              "Error from existsOrErr. User does not exist",
		inputUserID:       2,
		inputEnabledState: false,
		expectedOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "user"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:              "Error. Internal. Update failed",
		inputUserID:       4,
		inputEnabledState: false,
		expectedOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg user 4"), internal.ErrMsgFmtUpdateFailed, "user's enabled state"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			userDTO := dto.User{ID: &test.inputUserID, EnabledState: &test.inputEnabledState}

			err := testUsersServices.UpdateEnabledState(userDTO)

			if test.expectedOutput == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)

				svcErr, ok := err.(*svcerr.ServiceError)
				assert.True(t, ok)

				assert.Equal(t, test.expectedOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedOutput.Category(), svcErr.Category())
			}
		})
	}
}

func TestUsersDelete(t *testing.T) {

	type testCase struct {
		name           string
		input          uint
		expectedOutput *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:           "Success",
		input:          1,
		expectedOutput: nil,
	})

	table = append(table, testCase{
		name:  "Error from existsOrErr. User does not exist",
		input: 2,
		expectedOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "user"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:  "Error. Internal. Delete failed",
		input: 4,
		expectedOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg user 4"), internal.ErrMsgFmtDeleteFailed, "user"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			err := testUsersServices.Delete(test.input)

			if test.expectedOutput == nil {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)

				svcErr, ok := err.(*svcerr.ServiceError)
				assert.True(t, ok)

				assert.Equal(t, test.expectedOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedOutput.Category(), svcErr.Category())
			}
		})
	}
}

/*
func TestUsersRetrieve(t *testing.T) {

	type testCase struct {
		name                   string
		input                  uint
		expectedUsernameOutput *domain.User
		expectedErrOutput      *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                   "Success",
		input:                  1,
		expectedUsernameOutput: &domain.User{ID: 1, Username: "username1", PasswordHash: "passwordHash", PasswordSalt: "salt"},
		expectedErrOutput:      nil,
	})

	table = append(table, testCase{
		name:                   "Error. User does not exist. Invalid input ID.",
		input:                  2,
		expectedUsernameOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "user"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                   "Error. User retrieval failed. Internal.",
		input:                  3,
		expectedUsernameOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("another error msg"), internal.ErrMsgFmtRetrievalFailed, "user"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			user, err := testUsersServices.Retrieve(test.input)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedUsernameOutput.Username, user.Username)
			} else {
				assert.NotNil(t, err)

				svcErr, ok := err.(*svcerr.ServiceError)
				assert.True(t, ok)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}
		})
	}

}
*/

func TestUsersRetrieveByUsername(t *testing.T) {

	type testCase struct {
		name                   string
		input                  string
		expectedUsernameOutput *domain.User
		expectedErrOutput      *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                   "Success",
		input:                  "userExists",
		expectedUsernameOutput: &domain.User{ID: 1, Username: "userExists", PasswordHash: "passwordHash", PasswordSalt: "salt"},
		expectedErrOutput:      nil,
	})

	table = append(table, testCase{
		name:                   "Error. Username does not exist. Invalid input ID.",
		input:                  "userDoesNotExist",
		expectedUsernameOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "username"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                   "Error. User retrieval failed. Internal.",
		input:                  "userDALerr",
		expectedUsernameOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("DAL error"), internal.ErrMsgFmtRetrievalFailed, "user"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			user, err := testUsersServices.RetrieveByUsername(test.input)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedUsernameOutput.Username, *user.Username)
			} else {

				svcErr, ok := err.(*svcerr.ServiceError)
				assert.True(t, ok)

				assert.NotNil(t, svcErr)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}
		})
	}
}

func TestUsersRetrieveByUsernameInternal(t *testing.T) {

	type testCase struct {
		name                   string
		input                  string
		expectedUsernameOutput *domain.User
		expectedErrOutput      *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                   "Success",
		input:                  "userExists",
		expectedUsernameOutput: &domain.User{ID: 1, Username: "userExists", PasswordHash: "passwordHash", PasswordSalt: "salt"},
		expectedErrOutput:      nil,
	})

	table = append(table, testCase{
		name:                   "Error. Username does not exist. Invalid input ID.",
		input:                  "userDoesNotExist",
		expectedUsernameOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "username"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                   "Error. User retrieval failed. Internal.",
		input:                  "userDALerr",
		expectedUsernameOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("DAL error"), internal.ErrMsgFmtRetrievalFailed, "user"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			user, svcErr := testUsersServices.retrieveByUsername(test.input)

			if test.expectedErrOutput == nil {
				assert.Nil(t, svcErr)
				assert.Equal(t, test.expectedUsernameOutput.Username, user.Username)
			} else {
				assert.NotNil(t, svcErr)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}
		})
	}
}
