package dataaccess

import (
	"regexp"
	"testing"

	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestHttpMethodsSelectByName(t *testing.T) {

	// SQL BASE CONFIG

	inputMethodName := "methodName"

	resultCollumnNames := sqlMock.NewRows([]string{"ID", "Name"})

	var errReturnedMock = errors.New("error mock")

	newExpectedQuery := func() *sqlmock.ExpectedQuery {
		return sqlMock.ExpectQuery(
			regexp.QuoteMeta("HttpMethodsSelectByName"),
		)
	}

	//var expectedQuery *sqlmock.ExpectedQuery

	// CASES

	type testCase struct {
		name         string
		outputErr    error
		outputMethod *domain.HttpMethod
	}
	table := make([]testCase, 0)

	methodID1 := uint(1)
	methodStr1 := "method1"

	collumnNames1 := *resultCollumnNames
	resultOneRowOK := collumnNames1.AddRow(methodID1, methodStr1)
	newExpectedQuery().
		WillReturnRows(resultOneRowOK)
	table = append(table, testCase{
		name: "Success. Result one row.",

		outputMethod: &domain.HttpMethod{ID: methodID1, Name: methodStr1},
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
			method, err := httpMethodsDALMock.SelectByName(inputMethodName)

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

					if assert.NotNil(t, method) {

						assert.Equal(t, test.outputMethod, method)
					}
				}
			}
		})
	}
	//_ = expectedQuery
	sqlExpectationsErr := sqlMock.ExpectationsWereMet()
	assert.Nil(t, sqlExpectationsErr)
}

func TestHttpMethodsSelectAll(t *testing.T) {

	// SQL BASE CONFIG

	resultCollumnNames := sqlMock.NewRows([]string{"ID", "Name"})

	var errReturnedMock = errors.New("error mock")

	newExpectedQuery := func() *sqlmock.ExpectedQuery {
		return sqlMock.ExpectQuery(
			regexp.QuoteMeta("HttpMethodsSelectAll"),
		)
	}

	//var expectedQuery *sqlmock.ExpectedQuery

	// CASES

	type testCase struct {
		name          string
		outputErr     error
		outputMethods []dto.HttpMethod
	}
	table := make([]testCase, 0)

	methodID1 := uint(1)
	methodID2 := uint(2)

	methodStr1 := "method1"
	methodStr2 := "method2"

	collumnNames1 := *resultCollumnNames
	resultOneRowOK := collumnNames1.AddRow(methodID1, methodStr1)
	newExpectedQuery().
		WillReturnRows(resultOneRowOK).
		RowsWillBeClosed()
	table = append(table, testCase{
		name: "Success. Result one row.",

		outputMethods: []dto.HttpMethod{
			{ID: &methodID1, Name: &methodStr1},
		},
	})

	collumnNames2 := *resultCollumnNames
	resultTwoRowsOK := collumnNames2.
		AddRow(methodID1, methodStr1).
		AddRow(methodID2, methodStr2)
	newExpectedQuery().
		WillReturnRows(resultTwoRowsOK).
		RowsWillBeClosed()
	table = append(table, testCase{
		name: "Success. Result two rows.",

		outputMethods: []dto.HttpMethod{
			{ID: &methodID1, Name: &methodStr1},
			{ID: &methodID2, Name: &methodStr2},
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
		AddRow(methodID1, methodStr1).
		AddRow(methodID1, methodStr1).
		AddRow(methodID1, methodStr1).
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
			methods, err := httpMethodsDALMock.SelectAll()

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

					for i, m := range test.outputMethods {

						if assert.NotNil(t, methods[i].ID) {

							assert.Equal(t, m.ID, methods[i].ID)
						}
						if assert.NotNil(t, methods[i].Name) {

							assert.Equal(t, m.Name, methods[i].Name)
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
