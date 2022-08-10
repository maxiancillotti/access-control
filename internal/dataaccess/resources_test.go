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

func TestResourcesDelete(t *testing.T) {

	type testCase struct {
		name            string
		inputResourceID uint
		returnsErr      bool
	}
	table := make([]testCase, 0)

	table = append(table, testCase{
		name:            "Success",
		inputResourceID: 1,
		returnsErr:      false,
	})

	table = append(table, testCase{
		name:            "Error running exec",
		inputResourceID: 0,
		returnsErr:      true,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			expectedQuery := sqlMock.ExpectExec(
				regexp.QuoteMeta("ResourcesDelete"),
			).
				WithArgs(
					sql.Named("id", test.inputResourceID),
				)

			errReturnedMock := errors.New("sql: couldn't connect to database")
			if test.returnsErr {
				expectedQuery.WillReturnError(errReturnedMock)
			} else {
				expectedQuery.WillReturnResult(sqlmock.NewResult(0, 1))
			}

			err := resourcesDALMock.Delete(test.inputResourceID)

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

func TestResourcesSelect(t *testing.T) {

	// SQL BASE CONFIG

	inputResourceID := uint(1)

	resultCollumnNames := sqlMock.NewRows([]string{"ID", "Path"})

	var errReturnedMock = errors.New("error mock")

	newExpectedQuery := func() *sqlmock.ExpectedQuery {
		return sqlMock.ExpectQuery(
			regexp.QuoteMeta("ResourcesSelect"),
		).
			WithArgs(sql.Named("id", inputResourceID))
	}

	//var expectedQuery *sqlmock.ExpectedQuery

	// CASES

	type testCase struct {
		name           string
		outputErr      error
		outputResource *domain.Resource
	}
	table := make([]testCase, 0)

	resourceID1 := uint(1)
	resourceName1 := "resource1"

	collumnNames1 := *resultCollumnNames
	resultOneRowOK := collumnNames1.AddRow(resourceID1, resourceName1)
	newExpectedQuery().
		WillReturnRows(resultOneRowOK)
	table = append(table, testCase{
		name: "Success. Result one row.",

		outputResource: &domain.Resource{ID: resourceID1, Path: resourceName1},
	})

	newExpectedQuery().
		WillReturnRows(sqlMock.NewRows([]string{"Collumn"}).AddRow("rowValue")) // invalid collumn qty
	table = append(table, testCase{
		name:      "Error SP Result Scan failed",
		outputErr: errors.New(errMsgSPResultScanFailed),
	})

	newExpectedQuery().
		WillReturnRows(resultCollumnNames) // result collums without values, only names
	table = append(table, testCase{
		name:      "Error Empty Result",
		outputErr: dal.ErrorDataAccessEmptyResult,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {

			// EXEC
			resource, err := resourcesDALMock.Select(inputResourceID)

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

					if assert.NotNil(t, resource) {
						assert.Equal(t, test.outputResource, resource)
					}
				}
			}
		})
	}
	//_ = expectedQuery
	sqlExpectationsErr := sqlMock.ExpectationsWereMet()
	assert.Nil(t, sqlExpectationsErr)
}

func TestResourcesSelectByPath(t *testing.T) {

	// SQL BASE CONFIG

	inputResourcePath := "resourcePath"

	resultCollumnNames := sqlMock.NewRows([]string{"ID", "Path"})

	var errReturnedMock = errors.New("error mock")

	newExpectedQuery := func() *sqlmock.ExpectedQuery {
		return sqlMock.ExpectQuery(
			regexp.QuoteMeta("ResourcesSelectByPath"),
		).
			WithArgs(sql.Named("path", inputResourcePath))
	}

	//var expectedQuery *sqlmock.ExpectedQuery

	// CASES

	type testCase struct {
		name           string
		outputErr      error
		outputResource *domain.Resource
	}
	table := make([]testCase, 0)

	resourceID1 := uint(1)
	resourceName1 := "resource1"
	collumnNames1 := *resultCollumnNames
	resultOneRowOK := collumnNames1.AddRow(resourceID1, resourceName1)

	newExpectedQuery().
		WillReturnRows(resultOneRowOK)
	table = append(table, testCase{
		name: "Success. Result one row.",

		outputResource: &domain.Resource{ID: resourceID1, Path: resourceName1},
	})

	newExpectedQuery().
		WillReturnRows(sqlMock.NewRows([]string{"Collumn"}).AddRow("rowValue")) // invalid collumn qty
	table = append(table, testCase{
		name:      "Error SP Result Scan failed",
		outputErr: errors.New(errMsgSPResultScanFailed),
	})

	newExpectedQuery().
		WillReturnRows(resultCollumnNames) // result collums without values, only names
	table = append(table, testCase{
		name:      "Error Empty Result",
		outputErr: dal.ErrorDataAccessEmptyResult,
	})

	////////////////////////////////////////////////////////////////

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {

			// EXEC
			resource, err := resourcesDALMock.SelectByPath(inputResourcePath)

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

					if assert.NotNil(t, resource) {
						assert.Equal(t, test.outputResource, resource)
					}
				}
			}
		})
	}
	//_ = expectedQuery
	sqlExpectationsErr := sqlMock.ExpectationsWereMet()
	assert.Nil(t, sqlExpectationsErr)
}

func TestResourcesSelectAll(t *testing.T) {

	// SQL BASE CONFIG

	resultCollumnNames := sqlMock.NewRows([]string{"ID", "Path"})

	var errReturnedMock = errors.New("error mock")

	newExpectedQuery := func() *sqlmock.ExpectedQuery {
		return sqlMock.ExpectQuery(
			regexp.QuoteMeta("ResourcesSelectAll"),
		)
	}

	//var expectedQuery *sqlmock.ExpectedQuery

	// CASES

	type testCase struct {
		name            string
		outputErr       error
		outputResources []dto.Resource
	}
	table := make([]testCase, 0)

	resourceID1 := uint(1)
	resourceID2 := uint(2)

	resourceName1 := "resourceName1"
	resourceName2 := "resourceName2"

	collumnNames1 := *resultCollumnNames
	resultOneRowOK := collumnNames1.AddRow(resourceID1, resourceName1)
	newExpectedQuery().
		WillReturnRows(resultOneRowOK).
		RowsWillBeClosed()
	table = append(table, testCase{
		name: "Success. Result one row.",

		outputResources: []dto.Resource{
			{ID: &resourceID1, Path: &resourceName1},
		},
	})

	collumnNames2 := *resultCollumnNames
	resultTwoRowsOK := collumnNames2.
		AddRow(resourceID1, resourceName1).
		AddRow(resourceID2, resourceName2)
	newExpectedQuery().
		WillReturnRows(resultTwoRowsOK).
		RowsWillBeClosed()
	table = append(table, testCase{
		name: "Success. Result two rows.",

		outputResources: []dto.Resource{
			{ID: &resourceID1, Path: &resourceName1},
			{ID: &resourceID2, Path: &resourceName2},
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
	// also error row must be > 2nd because it must be higher than the longest successful result, for some strange
	// reason. Otherwise, an error occurs in this cases too.
	resultRowErr := collumnNames4.
		AddRow(resourceID1, resourceName1).
		AddRow(resourceID1, resourceName1).
		AddRow(resourceID1, resourceName1).
		RowError(2, errReturnedMock)
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
			resources, err := resourcesDALMock.SelectAll()

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

					for i, m := range test.outputResources {

						if assert.NotNil(t, resources[i].ID) {

							assert.Equal(t, m.ID, resources[i].ID)
						}
						if assert.NotNil(t, resources[i].Path) {

							assert.Equal(t, m.Path, resources[i].Path)
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
