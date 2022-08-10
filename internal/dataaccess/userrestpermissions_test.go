package dataaccess

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUsersRESTPermissionsInsert(t *testing.T) {

	type testCase struct {
		name       string
		inputURP   domain.UserRESTPermission
		returnsErr bool
	}
	table := make([]testCase, 0)

	table = append(table, testCase{
		name:       "Success",
		inputURP:   domain.UserRESTPermission{UserID: 1, ResourceID: 1, MethodID: 1},
		returnsErr: false,
	})

	table = append(table, testCase{
		name:       "Error running exec",
		inputURP:   domain.UserRESTPermission{},
		returnsErr: true,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectExec(
				regexp.QuoteMeta("UsersRESTPermissions_Insert"),
			).
				WithArgs(
					sql.Named("userID", test.inputURP.UserID),
					sql.Named("resourceID", test.inputURP.ResourceID),
					sql.Named("methodID", test.inputURP.MethodID),
				)

			errReturnedMock := errors.New("sql: couldn't connect to database")
			if test.returnsErr {
				expectedQuery.WillReturnError(errReturnedMock)
			} else {
				expectedQuery.WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err := userRESTPermDALMock.Insert(&test.inputURP)

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

func TestUsersRESTPermissionsDelete(t *testing.T) {

	type testCase struct {
		name       string
		inputURP   domain.UserRESTPermission
		returnsErr bool
	}
	table := make([]testCase, 0)

	table = append(table, testCase{
		name:       "Success",
		inputURP:   domain.UserRESTPermission{UserID: 1, ResourceID: 1, MethodID: 1},
		returnsErr: false,
	})

	table = append(table, testCase{
		name:       "Error running exec",
		inputURP:   domain.UserRESTPermission{},
		returnsErr: true,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectExec(
				regexp.QuoteMeta("UsersRESTPermissions_Delete"),
			).
				WithArgs(
					sql.Named("userID", test.inputURP.UserID),
					sql.Named("resourceID", test.inputURP.ResourceID),
					sql.Named("methodID", test.inputURP.MethodID),
				)

			errReturnedMock := errors.New("sql: couldn't connect to database")
			if test.returnsErr {
				expectedQuery.WillReturnError(errReturnedMock)
			} else {
				expectedQuery.WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err := userRESTPermDALMock.Delete(&test.inputURP)

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

func TestUsersRESTPermissionsSelectAllWithDescriptionsByUserID(t *testing.T) {

	// SQL BASE CONFIG

	resultCollumnNames := sqlMock.NewRows([]string{"ResourceID", "ResourcePath", "MethodID", "MethodName"})

	var inputUserID uint = 1
	var errReturnedMock = errors.New("error mock")

	newExpectedQuery := func() *sqlmock.ExpectedQuery {
		return sqlMock.ExpectQuery(
			regexp.QuoteMeta("UsersRESTPermissions_SelectAllWithDescriptionsByUserID"),
		).
			WithArgs(
				sql.Named("userID", inputUserID),
			)
	}

	//var expectedQuery *sqlmock.ExpectedQuery

	// CASES

	type testCase struct {
		name              string
		outputErr         error
		outputPermissions []dto.PermissionsWithDescriptions
	}
	table := make([]testCase, 0)

	resourceID1 := uint(1)
	resourcePathStr1 := "resource1"
	methodID1 := uint(1)
	methodNameStr1 := "method1"

	resourceID2 := uint(2)
	resourcePathStr2 := "resource2"
	methodID2 := uint(2)
	methodNameStr2 := "method2"

	collumnNames1 := *resultCollumnNames
	resultOneRowOK := collumnNames1.AddRow(resourceID1, resourcePathStr1, methodID1, methodNameStr1)

	newExpectedQuery().
		WillReturnRows(resultOneRowOK).
		RowsWillBeClosed()
	table = append(table, testCase{
		name: "Success. Result one row.",

		outputPermissions: []dto.PermissionsWithDescriptions{
			{
				Resource: dto.Resource{
					ID:   &resourceID1,
					Path: &resourcePathStr1,
				},

				Methods: []dto.HttpMethod{
					{
						ID:   &methodID1,
						Name: &methodNameStr1,
					},
				},
			},
		},
	})

	collumnNames2 := *resultCollumnNames
	resultTwoRowsOK := collumnNames2.
		AddRow(resourceID1, resourcePathStr1, methodID1, methodNameStr1).
		AddRow(resourceID1, resourcePathStr1, methodID2, methodNameStr2)
	newExpectedQuery().
		WillReturnRows(resultTwoRowsOK).
		RowsWillBeClosed()
	table = append(table, testCase{
		name: "Success. Result multiple rows, one resource, more than one method.",

		outputPermissions: []dto.PermissionsWithDescriptions{
			{
				Resource: dto.Resource{
					ID:   &resourceID1,
					Path: &resourcePathStr1,
				},

				Methods: []dto.HttpMethod{
					{
						ID:   &methodID1,
						Name: &methodNameStr1,
					},
					{
						ID:   &methodID2,
						Name: &methodNameStr2,
					},
				},
			},
		},
	})

	collumnNames3 := *resultCollumnNames
	resultThreeRowsOK := collumnNames3.
		AddRow(resourceID1, resourcePathStr1, methodID1, methodNameStr1).
		AddRow(resourceID1, resourcePathStr1, methodID2, methodNameStr2).
		AddRow(resourceID2, resourcePathStr2, methodID2, methodNameStr2)
	newExpectedQuery().
		WillReturnRows(resultThreeRowsOK).
		RowsWillBeClosed()
	table = append(table, testCase{
		name: "Success. Result multiple rows, more than one resource, more than one method.",

		outputPermissions: []dto.PermissionsWithDescriptions{
			{
				Resource: dto.Resource{
					ID:   &resourceID1,
					Path: &resourcePathStr1,
				},

				Methods: []dto.HttpMethod{
					{
						ID:   &methodID1,
						Name: &methodNameStr1,
					},
					{
						ID:   &methodID2,
						Name: &methodNameStr2,
					},
				},
			},

			{
				Resource: dto.Resource{
					ID:   &resourceID2,
					Path: &resourcePathStr2,
				},

				Methods: []dto.HttpMethod{
					{
						ID:   &methodID2,
						Name: &methodNameStr2,
					},
				},
			},
		},
	})

	newExpectedQuery().
		WillReturnError(errReturnedMock)
	table = append(table, testCase{
		name:      "Error SP Exec failed",
		outputErr: errors.New(errMsgSPExecFailed),
	})

	newExpectedQuery().
		WillReturnRows(sqlMock.NewRows([]string{"Collumn"}).AddRow("rowValue")). // invalid collumn qty
		RowsWillBeClosed()
	table = append(table, testCase{
		name:      "Error SP Result Scan failed",
		outputErr: errors.New(errMsgSPResultScanFailed),
	})

	collumnNames4 := *resultCollumnNames
	// must be > 1 row because otherwise it doesn't go through the loop, and the returned err is empty result.
	// also error row must be > 3rd because it must be higher than the longest successful result, for some strange
	// reason. Otherwise, an error occurs in this cases too.
	resultRowErr := collumnNames4.
		AddRow(resourceID1, resourcePathStr1, methodID1, methodNameStr1).
		AddRow(resourceID1, resourcePathStr1, methodID2, methodNameStr2).
		AddRow(resourceID1, resourcePathStr1, methodID2, methodNameStr2).
		AddRow(resourceID1, resourcePathStr1, methodID2, methodNameStr2).
		RowError(3, errReturnedMock)
	newExpectedQuery().
		WillReturnRows(resultRowErr).
		RowsWillBeClosed()
	table = append(table, testCase{
		name:      "Error Rows Scan returned an err",
		outputErr: errors.New(errMsgRowsScanReturnedAnError),
	})

	newExpectedQuery().
		WillReturnRows(resultCollumnNames). // result collums without values, only names
		RowsWillBeClosed()
	table = append(table, testCase{
		name:      "Error Empty Result",
		outputErr: dal.ErrorDataAccessEmptyResult,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {

			// EXEC
			permissions, err := userRESTPermDALMock.SelectAllWithDescriptionsByUserID(inputUserID)

			// CHECK
			if test.outputErr != nil {
				if assert.NotNil(t, err) {

					switch test.outputErr.Error() {

					case errMsgSPExecFailed:
						assert.Contains(t, err.Error(), errReturnedMock.Error())
						assert.Contains(t, err.Error(), errMsgSPExecFailed)

					case errMsgSPResultScanFailed:
						assert.Contains(t, err.Error(), errMsgSPResultScanFailed)

					case errMsgRowsScanReturnedAnError:
						assert.Contains(t, err.Error(), errReturnedMock.Error())
						assert.Contains(t, err.Error(), errMsgRowsScanReturnedAnError)

					case dal.ErrorDataAccessEmptyResult.Error():
						assert.ErrorIs(t, dal.ErrorDataAccessEmptyResult, err)
					}
				}

			} else {
				if assert.Nil(t, err) {
					/*
						assert.Equal(t, test.outputPermissions[0].Resource.ID, permissions[0].Resource.ID)
						assert.Equal(t, test.outputPermissions[0].Resource.Path, permissions[0].Resource.Path)
						assert.Equal(t, test.outputPermissions[0].Methods[0].ID, permissions[0].Methods[0].ID)
						assert.Equal(t, test.outputPermissions[0].Methods[0].Name, permissions[0].Methods[0].Name)

					*/
					for i, expectedR := range test.outputPermissions {

						assert.Equal(t, expectedR.Resource.ID, permissions[i].Resource.ID)
						assert.Equal(t, expectedR.Resource.Path, permissions[i].Resource.Path)

						for j, expectedM := range expectedR.Methods {

							assert.Equal(t, expectedM.ID, permissions[i].Methods[j].ID)
							assert.Equal(t, expectedM.Name, permissions[i].Methods[j].Name)
						}
					}
				}
			}
		})
	}
	//_ = expectedQuery
	sqlExpectationsErr := sqlMock.ExpectationsWereMet()
	assert.Nil(t, sqlExpectationsErr)
}

func TestUsersRESTPermissionsSelectAllPathMethodsByUserID(t *testing.T) {

	// SQL BASE CONFIG

	resultCollumnNames := sqlMock.NewRows([]string{"Resource", "Method"})

	var inputUserID uint = 1
	var errReturnedMock = errors.New("error mock")

	newExpectedQuery := func() *sqlmock.ExpectedQuery {
		return sqlMock.ExpectQuery(
			regexp.QuoteMeta("UsersRESTPermissions_SelectAllPathsMethodsByUserID"),
		).
			WithArgs(
				sql.Named("userID", inputUserID),
			)
	}

	//var expectedQuery *sqlmock.ExpectedQuery

	// CASES

	type testCase struct {
		name              string
		outputErr         error
		outputPermissions domain.RESTPermissionsPathsMethods
	}
	table := make([]testCase, 0)

	resourcePathStr1 := "resource1"
	methodNameStr1 := "method1"

	resourcePathStr2 := "resource2"
	methodNameStr2 := "method2"

	collumnNames1 := *resultCollumnNames
	resultOneRowOK := collumnNames1.AddRow(resourcePathStr1, methodNameStr1)

	newExpectedQuery().
		WillReturnRows(resultOneRowOK).
		RowsWillBeClosed()
	table = append(table, testCase{
		name: "Success. Result one row.",

		outputPermissions: domain.RESTPermissionsPathsMethods{
			resourcePathStr1: {
				methodNameStr1,
			},
		},
	})

	collumnNames2 := *resultCollumnNames
	resultTwoRowsOK := collumnNames2.
		AddRow(resourcePathStr1, methodNameStr1).
		AddRow(resourcePathStr1, methodNameStr2)
	newExpectedQuery().
		WillReturnRows(resultTwoRowsOK).
		RowsWillBeClosed()
	table = append(table, testCase{
		name: "Success. Result multiple rows, one resource, more than one method.",

		outputPermissions: domain.RESTPermissionsPathsMethods{
			resourcePathStr1: {
				methodNameStr1,
				methodNameStr2,
			},
		},
	})

	collumnNames3 := *resultCollumnNames
	resultThreeRowsOK := collumnNames3.
		AddRow(resourcePathStr1, methodNameStr1).
		AddRow(resourcePathStr1, methodNameStr2).
		AddRow(resourcePathStr2, methodNameStr2)
	newExpectedQuery().
		WillReturnRows(resultThreeRowsOK).
		RowsWillBeClosed()
	table = append(table, testCase{
		name: "Success. Result multiple rows, more than one resource, more than one method.",

		outputPermissions: domain.RESTPermissionsPathsMethods{
			resourcePathStr1: {
				methodNameStr1,
				methodNameStr2,
			},
			resourcePathStr2: {
				methodNameStr2,
			},
		},
	})

	newExpectedQuery().
		WillReturnError(errReturnedMock)
	table = append(table, testCase{
		name:      "Error SP Exec failed",
		outputErr: errors.New(errMsgSPExecFailed),
	})

	newExpectedQuery().
		WillReturnRows(sqlMock.NewRows([]string{"Collumn"}).AddRow("rowValue")). // invalid collumn qty
		RowsWillBeClosed()
	table = append(table, testCase{
		name:      "Error SP Result Scan failed",
		outputErr: errors.New(errMsgSPResultScanFailed),
	})

	collumnNames4 := *resultCollumnNames
	// must be > 1 row because otherwise it doesn't go through the loop, and the returned err is empty result.
	// also error row must be > 3rd because it must be higher than the longest successful result, for some strange
	// reason. Otherwise, an error occurs in this cases too.
	resultRowErr := collumnNames4.
		AddRow(resourcePathStr1, methodNameStr1).
		AddRow(resourcePathStr1, methodNameStr2).
		AddRow(resourcePathStr1, methodNameStr2).
		AddRow(resourcePathStr1, methodNameStr2).
		RowError(3, errReturnedMock)
	newExpectedQuery().
		WillReturnRows(resultRowErr).
		RowsWillBeClosed()
	table = append(table, testCase{
		name:      "Error Rows Scan returned an err",
		outputErr: errors.New(errMsgRowsScanReturnedAnError),
	})

	newExpectedQuery().
		WillReturnRows(resultCollumnNames). // result collums without values, only names
		RowsWillBeClosed()
	table = append(table, testCase{
		name:      "Error Empty Result",
		outputErr: dal.ErrorDataAccessEmptyResult,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {

			// EXEC
			permissions, err := userRESTPermDALMock.SelectAllPathMethodsByUserID(inputUserID)

			// CHECK
			if test.outputErr != nil {
				if assert.NotNil(t, err) {

					switch test.outputErr.Error() {

					case errMsgSPExecFailed:
						assert.Contains(t, err.Error(), errReturnedMock.Error())
						assert.Contains(t, err.Error(), errMsgSPExecFailed)

					case errMsgSPResultScanFailed:
						assert.Contains(t, err.Error(), errMsgSPResultScanFailed)

					case errMsgRowsScanReturnedAnError:
						assert.Contains(t, err.Error(), errReturnedMock.Error())
						assert.Contains(t, err.Error(), errMsgRowsScanReturnedAnError)

					case dal.ErrorDataAccessEmptyResult.Error():
						assert.ErrorIs(t, dal.ErrorDataAccessEmptyResult, err)
					}
				}

			} else {
				if assert.Nil(t, err) {

					for key, expectedMethods := range test.outputPermissions {

						for i, m := range expectedMethods {

							if assert.NotNil(t, permissions[key][i]) {

								assert.Equal(t, m, permissions[key][i])
							}
						}
					}
				}
			}
		})
	}
	//_ = expectedQuery
	sqlExpectationsErr := sqlMock.ExpectationsWereMet()
	assert.Nil(t, sqlExpectationsErr)
}
