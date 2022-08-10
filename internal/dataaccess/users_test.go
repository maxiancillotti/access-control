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

func TestUsersUpdatePassword(t *testing.T) {

	type testCase struct {
		name       string
		inputUser  domain.User
		returnsErr bool
	}
	table := make([]testCase, 0)

	table = append(table, testCase{
		name:       "Success",
		inputUser:  domain.User{ID: 1, PasswordHash: "passwordHash", PasswordSalt: "pwsalt"},
		returnsErr: false,
	})

	table = append(table, testCase{
		name:       "Error running exec",
		inputUser:  domain.User{ID: 2, PasswordHash: "passwordHash", PasswordSalt: "pwsalt"},
		returnsErr: true,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectExec(
				regexp.QuoteMeta("UsersUpdatePassword"),
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

			err := userDALMock.UpdatePassword(&test.inputUser)

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

func TestUsersUpdateEnabledState(t *testing.T) {

	type testCase struct {
		name       string
		inputUser  domain.User
		returnsErr bool
	}
	table := make([]testCase, 0)

	table = append(table, testCase{
		name:       "Success",
		inputUser:  domain.User{ID: 1, Enabled: true},
		returnsErr: false,
	})

	table = append(table, testCase{
		name:       "Error running exec",
		inputUser:  domain.User{ID: 2, Enabled: true},
		returnsErr: true,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectExec(
				regexp.QuoteMeta("UsersUpdateEnableState"),
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

			err := userDALMock.UpdateEnabledState(&test.inputUser)

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

func TestUsersDelete(t *testing.T) {

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
				regexp.QuoteMeta("UsersDelete"),
			).
				WithArgs(test.inputUserID)

			errReturnedMock := errors.New("sql: couldn't connect to database")
			if test.returnsErr {
				expectedQuery.WillReturnError(errReturnedMock)
			} else {
				expectedQuery.WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err := userDALMock.Delete(test.inputUserID)

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

func TestUsersSelect(t *testing.T) {

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
				regexp.QuoteMeta("UsersSelect"),
			).
				WithArgs(test.inputUserID)

			if test.resultRows != nil {
				expectedQuery.WillReturnRows(resultCollums)
			} else {
				expectedQuery.WillReturnError(errors.New("sql: couldn't connect to database"))
			}

			resultUser, err := userDALMock.Select(test.inputUserID)

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

func TestUsersSelectByUsername(t *testing.T) {

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
				regexp.QuoteMeta("UsersSelectByUsername"),
			).
				WithArgs(username)

			if test.resultRows != nil {
				expectedQuery.WillReturnRows(resultCollums)
			} else {
				expectedQuery.WillReturnError(errors.New("sql: couldn't connect to database"))
			}

			resultUser, err := userDALMock.SelectByUsername(username)

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

/*
func TestUsersExists(t *testing.T) {

	type testCase struct {
		name          string
		inputUserID   uint
		returnsExists bool
		returnsErr    bool
	}
	table := make([]testCase, 0)

	table = append(table, testCase{
		name:          "Success. Exists",
		inputUserID:   1,
		returnsExists: true,
		returnsErr:    false,
	})

	// table = append(table, testCase{
	// 	name:          "Success. Does not exist",
	// 	inputUserID:   2,
	// 	returnsExists: false,
	// 	returnsErr:    false,
	// })

	// table = append(table, testCase{
	// 	name:          "Error running exec",
	// 	inputUserID:   0,
	// 	returnsExists: false,
	// 	returnsErr:    true,
	// })

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectExec(
				regexp.QuoteMeta("UsersExists"),
			).
				//WithArgs(test.inputUserID, test.returnsExists)
				//WithArgs(test.inputUserID, sql.Out{Dest: &test.returnsExists})
				//WithArgs(test.inputUserID, &test.returnsExists)
				WithArgs(
					sql.Named("id", test.inputUserID),
					sql.Named("exists", sql.Out{Dest: &test.returnsExists}),
				)

			errReturnedMock := errors.New("sql: couldn't connect to database")
			if test.returnsErr {
				expectedQuery.WillReturnError(errReturnedMock)
			} else {
				expectedQuery.WillReturnResult(sqlmock.NewResult(0, 0))
			}

			exists, err := userDALMock.Exists(test.inputUserID)

			if test.returnsErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), errMsgSPExecFailed)
				assert.Contains(t, err.Error(), errReturnedMock.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.returnsExists, exists)
			}
		})
	}
}
*/
/*
func TestSelectUserPermissionsByUserID(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		var userID uint = 1

		rows := sqlMock.NewRows([]string{"Resource", "Permission"}).
			AddRow("/customers", "GET").
			AddRow("/customers", "POST")

		sqlMock.ExpectQuery(
			regexp.QuoteMeta("EXEC [UsersSelectAllResourcesPermissionsByUserID] ?"),
		).
			WithArgs(userID).
			WillReturnRows(rows).
			RowsWillBeClosed()

		resultPerm, err := userDALMock.SelectUserPermissionsByUserID(userID)

		assert.Nil(t, err)
		assert.NotNil(t, resultPerm)
	})

	t.Run("Error: Empty result", func(t *testing.T) {

		var userID uint = 2

		rows := sqlMock.NewRows([]string{"Resource", "Permission"})

		sqlMock.ExpectQuery(
			regexp.QuoteMeta("EXEC [UsersSelectAllResourcesPermissionsByUserID] ?"),
		).
			WithArgs(userID).
			WillReturnRows(rows)

		resultPerm, err := userDALMock.SelectUserPermissionsByUserID(userID)

		assert.NotNil(t, err)
		assert.Nil(t, resultPerm)

		//isErr := strings.Contains(err.Error(), interactors.ErrorDataAccessEmptyResult.Error())
		isErr := errors.Is(err, services.ErrorDataAccessEmptyResult)
		assert.True(t, isErr)

		//t.Log("err value:", err)
	})

	t.Run("Error running query", func(t *testing.T) {

		var userID uint = 1

		//rows := sqlMock.NewRows([]string{"Resource", "Permission"}

		sqlMock.ExpectQuery(
			regexp.QuoteMeta("EXEC [UsersSelectAllResourcesPermissionsByUserID] ?"),
		).
			WithArgs(userID).
			WillReturnError(errors.New("sql: couldn't connect to database"))

		resultPerm, err := userDALMock.SelectUserPermissionsByUserID(userID)

		assert.NotNil(t, err)
		assert.Nil(t, resultPerm)

		isErr := strings.Contains(err.Error(), errMsgSPExecFailed)
		assert.True(t, isErr)

		t.Log("err value:", err)
	})
}
*/
