package dataaccess

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/services/dal"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAdminsUpdatePassword(t *testing.T) {

	type testCase struct {
		name       string
		inputUser  domain.Admin
		returnsErr bool
	}
	table := make([]testCase, 0)

	table = append(table, testCase{
		name:       "Success",
		inputUser:  domain.Admin{ID: 1, PasswordHash: "passwordHash", PasswordSalt: "pwsalt"},
		returnsErr: false,
	})

	table = append(table, testCase{
		name:       "Error running exec",
		inputUser:  domain.Admin{ID: 2, PasswordHash: "passwordHash", PasswordSalt: "pwsalt"},
		returnsErr: true,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectExec(
				regexp.QuoteMeta("AdminsUpdatePassword"),
			).
				WithArgs(
					sql.Named("id", test.inputUser.ID),
					sql.Named("passwordHash", test.inputUser.PasswordHash),
					sql.Named("passwordSalt", test.inputUser.PasswordSalt),
				)

			errReturnedMock := errors.New("sql: couldn't connect to database")
			if test.returnsErr {
				expectedQuery.WillReturnError(errReturnedMock)
			} else {
				expectedQuery.WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err := adminsDALMock.UpdatePassword(&test.inputUser)

			if test.returnsErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), errMsgSPExecFailed)
				assert.Contains(t, err.Error(), errReturnedMock.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAdminsUpdateEnabledState(t *testing.T) {

	type testCase struct {
		name       string
		inputUser  domain.Admin
		returnsErr bool
	}
	table := make([]testCase, 0)

	table = append(table, testCase{
		name:       "Success",
		inputUser:  domain.Admin{ID: 1, Enabled: true},
		returnsErr: false,
	})

	table = append(table, testCase{
		name:       "Error running exec",
		inputUser:  domain.Admin{ID: 2, Enabled: true},
		returnsErr: true,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectExec(
				regexp.QuoteMeta("AdminsUpdateEnableState"),
			).
				WithArgs(
					sql.Named("id", test.inputUser.ID),
					sql.Named("enabled", test.inputUser.Enabled),
				)

			errReturnedMock := errors.New("sql: couldn't connect to database")
			if test.returnsErr {
				expectedQuery.WillReturnError(errReturnedMock)
			} else {
				expectedQuery.WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err := adminsDALMock.UpdateEnabledState(&test.inputUser)

			if test.returnsErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), errMsgSPExecFailed)
				assert.Contains(t, err.Error(), errReturnedMock.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAdminsDelete(t *testing.T) {

	type testCase struct {
		name        string
		inputUserID uint
		returnsErr  bool
	}
	table := make([]testCase, 0)

	table = append(table, testCase{
		name:        "Success",
		inputUserID: 1,
		returnsErr:  false,
	})

	table = append(table, testCase{
		name:        "Error running exec",
		inputUserID: 0,
		returnsErr:  true,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectExec(
				regexp.QuoteMeta("AdminsDelete"),
			).
				WithArgs(test.inputUserID)

			errReturnedMock := errors.New("sql: couldn't connect to database")
			if test.returnsErr {
				expectedQuery.WillReturnError(errReturnedMock)
			} else {
				expectedQuery.WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err := adminsDALMock.Delete(test.inputUserID)

			if test.returnsErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), errMsgSPExecFailed)
				assert.Contains(t, err.Error(), errReturnedMock.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAdminsSelect(t *testing.T) {

	resultCollums := sqlMock.NewRows([]string{"ID", "Username", "PasswordHash", "PasswordSalt", "Enabled"})

	type testCase struct {
		name        string
		inputUserID uint
		resultRows  *sqlmock.Rows
	}
	table := make([]testCase, 0)

	table = append(table, testCase{
		name:        "Success",
		inputUserID: 1,
		resultRows:  resultCollums.AddRow(1, "username", "passwordHash", "passwordSalt", 1),
	})

	table = append(table, testCase{
		name:        "Error: Empty result",
		inputUserID: 2,
		resultRows:  resultCollums, // No result. Only collums.
	})

	table = append(table, testCase{
		name:        "Error running query",
		inputUserID: 0,
		resultRows:  nil, // No result nor collums. Returns err.
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectQuery(
				regexp.QuoteMeta("AdminsSelect"),
			).
				WithArgs(test.inputUserID)

			if test.resultRows != nil {
				expectedQuery.WillReturnRows(resultCollums)
			} else {
				expectedQuery.WillReturnError(errors.New("sql: couldn't connect to database"))
			}

			resultUser, err := adminsDALMock.Select(test.inputUserID)

			if test.inputUserID == 1 {
				assert.Nil(t, err)
				assert.NotNil(t, resultUser)
				assert.Equal(t, 1, int(resultUser.ID))
				assert.Equal(t, "username", resultUser.Username)
				assert.Equal(t, "passwordHash", resultUser.PasswordHash)
				assert.Equal(t, "passwordSalt", resultUser.PasswordSalt)
				assert.True(t, resultUser.Enabled)

			} else {
				assert.NotNil(t, err)
				assert.Nil(t, resultUser)

				if test.inputUserID == 2 {
					//assert.Contains(t, err.Error(), dal.ErrorDataAccessEmptyResult.Error())
					assert.ErrorIs(t, err, dal.ErrorDataAccessEmptyResult)
				} else {
					//assert.NotContains(t, err.Error(), dal.ErrorDataAccessEmptyResult.Error())
					assert.NotErrorIs(t, err, dal.ErrorDataAccessEmptyResult)
				}
				//t.Log("err value:", err)
			}
		})
	}
}

func TestAdminsSelectByUsername(t *testing.T) {

	resultCollums := sqlMock.NewRows([]string{"ID", "Username", "PasswordHash", "PasswordSalt", "Enabled"})

	type testCase struct {
		name          string
		inputUsername string
		resultRows    *sqlmock.Rows
	}
	table := make([]testCase, 0)

	var username string

	username = "userExists"
	table = append(table, testCase{
		name:          "Success",
		inputUsername: username,
		resultRows:    resultCollums.AddRow(1, username, "passwordHash", "passwordSalt", 1),
	})

	username = "userDoesNotExists"
	table = append(table, testCase{
		name:          "Error: Empty result",
		inputUsername: username,
		resultRows:    resultCollums, // No result. Only collums.
	})

	username = "userReturnsErr"
	table = append(table, testCase{
		name:          "Error running query",
		inputUsername: username,
		resultRows:    nil, // No result nor collums. Returns err.
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectQuery(
				regexp.QuoteMeta("AdminsSelectByUsername"),
			).
				WithArgs(username)

			if test.resultRows != nil {
				expectedQuery.WillReturnRows(resultCollums)
			} else {
				expectedQuery.WillReturnError(errors.New("sql: couldn't connect to database"))
			}

			resultUser, err := adminsDALMock.SelectByUsername(username)

			if test.inputUsername == "userExists" {
				assert.Nil(t, err)
				assert.NotNil(t, resultUser)
				assert.Equal(t, 1, int(resultUser.ID))
				assert.Equal(t, "userExists", resultUser.Username)
				assert.Equal(t, "passwordHash", resultUser.PasswordHash)
				assert.Equal(t, "passwordSalt", resultUser.PasswordSalt)
				assert.True(t, resultUser.Enabled)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, resultUser)

				if test.inputUsername == "userDoesNotExists" {
					//assert.Contains(t, err.Error(), dal.ErrorDataAccessEmptyResult.Error())
					assert.ErrorIs(t, err, dal.ErrorDataAccessEmptyResult)
				} else {
					//assert.NotContains(t, err.Error(), dal.ErrorDataAccessEmptyResult.Error())
					assert.NotErrorIs(t, err, dal.ErrorDataAccessEmptyResult)
				}
				//t.Log("err value:", err)
			}
		})
	}
}
