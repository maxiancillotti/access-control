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

func TestAdminsGetRandPasswordAndSalt(t *testing.T) {

	testAdmSvcs := testAdminsServices.(*adminsInteractor)
	pw, salt := testAdmSvcs.getRandPasswordAndSalt()

	pwLen := len(pw)
	saltLen := len(salt)

	assert.Equal(t, 64, pwLen)
	assert.Equal(t, 32, saltLen)
}

func TestAdminsCreate(t *testing.T) {

	type testCase struct {
		name                      string
		usernameInput             string
		expectedUserOutput        *dto.Admin
		expectedPasswordLenOutput int
		expectedErrOutput         *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	outputUserIDCase1 := uint(2)
	outputUserUsernameCase1 := "userDoesNotExist"

	table = append(table, testCase{
		name:                      "Success",
		usernameInput:             "userDoesNotExist",
		expectedUserOutput:        &dto.Admin{ID: &outputUserIDCase1, Username: &outputUserUsernameCase1},
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
			errors.Wrapf(errors.New("error msg"), internal.ErrMsgFmtInsertFailed, "admin"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			adminDTO, err := testAdminsServices.Create(dto.Admin{Username: &test.usernameInput})

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.NotNil(t, adminDTO)

				assert.Equal(t, test.expectedUserOutput.ID, adminDTO.ID)
				assert.Equal(t, test.expectedUserOutput.Username, adminDTO.Username)

				pwLen := len(*adminDTO.Password)
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

func TestAdminsUpdatePassword(t *testing.T) {

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
		name:                      "Error from existsOrErr. Admin does not exist",
		input:                     2,
		expectedPasswordLenOutput: 0,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "admin"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                      "Error. Internal. Update failed",
		input:                     4,
		expectedPasswordLenOutput: 0,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg admin 4"), internal.ErrMsgFmtUpdateFailed, "admin's password"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			pw, err := testAdminsServices.UpdatePassword(test.input)

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

func TestAdminsExistsOrErr(t *testing.T) {

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
		name:  "Error. Admin does not exist.",
		input: 2,
		expectedOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "admin"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:  "Error. Failed to check if admin exist",
		input: 3,
		expectedOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg"), internal.ErrMsgFmtFailedToCheckIfExists, "admin"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			testUsrSvcs := testAdminsServices.(*adminsInteractor)
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

func TestAdminsUpdateEnabledState(t *testing.T) {

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
		name:              "Error from existsOrErr. Admin does not exist",
		inputUserID:       2,
		inputEnabledState: false,
		expectedOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "admin"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:              "Error. Internal. Update failed",
		inputUserID:       4,
		inputEnabledState: false,
		expectedOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg admin 4"), internal.ErrMsgFmtUpdateFailed, "admin's enabled state"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			adminDTO := dto.Admin{ID: &test.inputUserID, EnabledState: &test.inputEnabledState}

			err := testAdminsServices.UpdateEnabledState(adminDTO)

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

func TestAdminsDelete(t *testing.T) {

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
		name:  "Error from existsOrErr. Admin does not exist",
		input: 2,
		expectedOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "admin"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:  "Error. Internal. Delete failed",
		input: 4,
		expectedOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg admin 4"), internal.ErrMsgFmtDeleteFailed, "admin"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			err := testAdminsServices.Delete(test.input)

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
func TestAdminsRetrieve(t *testing.T) {

	type testCase struct {
		name                   string
		input                  uint
		expectedUsernameOutput *domain.Admin
		expectedErrOutput      *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                   "Success",
		input:                  1,
		expectedUsernameOutput: &domain.Admin{ID: 1, Username: "username1", Password: "password", PasswordSalt: "salt"},
		expectedErrOutput:      nil,
	})

	table = append(table, testCase{
		name:                   "Error. Admin does not exist. Invalid input ID.",
		input:                  2,
		expectedUsernameOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "admin"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                   "Error. Admin retrieval failed. Internal.",
		input:                  3,
		expectedUsernameOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("another error msg"), internal.ErrMsgFmtRetrievalFailed, "admin"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			admin, err := testAdminsServices.Retrieve(test.input)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedUsernameOutput.Username, admin.Username)
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

func TestAdminsRetrieveByUsername(t *testing.T) {

	type testCase struct {
		name                   string
		input                  string
		expectedUsernameOutput *domain.Admin
		expectedErrOutput      *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                   "Success",
		input:                  "userExists",
		expectedUsernameOutput: &domain.Admin{ID: 1, Username: "userExists", PasswordHash: "password", PasswordSalt: "salt"},
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
		name:                   "Error. Admin retrieval failed. Internal.",
		input:                  "userDALerr",
		expectedUsernameOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("DAL error"), internal.ErrMsgFmtRetrievalFailed, "admin"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			admin, err := testAdminsServices.RetrieveByUsername(test.input)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedUsernameOutput.Username, *admin.Username)
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

func TestAdminsRetrieveByUsernameInternal(t *testing.T) {

	type testCase struct {
		name                   string
		input                  string
		expectedUsernameOutput *domain.Admin
		expectedErrOutput      *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                   "Success",
		input:                  "userExists",
		expectedUsernameOutput: &domain.Admin{ID: 1, Username: "userExists", PasswordHash: "password", PasswordSalt: "salt"},
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
		name:                   "Error. Admin retrieval failed. Internal.",
		input:                  "userDALerr",
		expectedUsernameOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("DAL error"), internal.ErrMsgFmtRetrievalFailed, "admin"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			admin, svcErr := testAdminsServices.retrieveByUsername(test.input)

			if test.expectedErrOutput == nil {
				assert.Nil(t, svcErr)
				assert.Equal(t, test.expectedUsernameOutput.Username, admin.Username)
			} else {
				assert.NotNil(t, svcErr)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}
		})
	}
}
