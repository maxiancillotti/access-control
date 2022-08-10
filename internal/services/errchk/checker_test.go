package errchk

import (
	"errors"
	"testing"

	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/stretchr/testify/assert"
)

const (
	testNameTrue                         = "True"
	testNameFalseNotServiceError         = "False. Not ServiceError concrete type."
	testNameFalseNotTheSameErrorCategory = "False. Not the same error category."
)

var (
	testErrchkr = NewServiceErrorChecker()
	errMsgTest  = errors.New("error msg")
)

type testTableErrorChecker struct {
	name           string
	input          error
	expectedOutput bool
	f              func(err error) bool
}

func TestErrorIsInvalidInputIdentifier(t *testing.T) {

	f := testErrchkr.ErrorIsInvalidInputIdentifier
	table := make([]testTableErrorChecker, 0)

	table = append(table, testTableErrorChecker{
		name: testNameTrue,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryInvalidInputID,
		),
		expectedOutput: true,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name:           testNameFalseNotServiceError,
		input:          errMsgTest,
		expectedOutput: false,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name: testNameFalseNotTheSameErrorCategory,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryInvalidCredentials,
		),
		expectedOutput: false,
		f:              f,
	})

	for _, table := range table {

		t.Run(table.name, func(t *testing.T) {
			output := table.f(table.input)
			assert.Equal(t, table.expectedOutput, output)
		})
	}
}

func TestErrorIsEmptyResult(t *testing.T) {

	f := testErrchkr.ErrorIsEmptyResult
	table := make([]testTableErrorChecker, 0)

	table = append(table, testTableErrorChecker{
		name: testNameTrue,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryEmptyResult,
		),
		expectedOutput: true,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name:           testNameFalseNotServiceError,
		input:          errMsgTest,
		expectedOutput: false,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name: testNameFalseNotTheSameErrorCategory,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryInvalidCredentials,
		),
		expectedOutput: false,
		f:              f,
	})

	for _, table := range table {

		t.Run(table.name, func(t *testing.T) {
			output := table.f(table.input)
			assert.Equal(t, table.expectedOutput, output)
		})
	}
}

func TestErrorIsInternal(t *testing.T) {

	f := testErrchkr.ErrorIsInternal
	table := make([]testTableErrorChecker, 0)

	table = append(table, testTableErrorChecker{
		name: testNameTrue,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryInternal,
		),
		expectedOutput: true,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name:           testNameFalseNotServiceError,
		input:          errMsgTest,
		expectedOutput: false,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name: testNameFalseNotTheSameErrorCategory,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryInvalidCredentials,
		),
		expectedOutput: false,
		f:              f,
	})

	for _, table := range table {

		t.Run(table.name, func(t *testing.T) {
			output := table.f(table.input)
			assert.Equal(t, table.expectedOutput, output)
		})
	}
}

func TestErrorIsInvalidCredentials(t *testing.T) {

	f := testErrchkr.ErrorIsInvalidCredentials
	table := make([]testTableErrorChecker, 0)

	table = append(table, testTableErrorChecker{
		name: testNameTrue,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryInvalidCredentials,
		),
		expectedOutput: true,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name:           testNameFalseNotServiceError,
		input:          errMsgTest,
		expectedOutput: false,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name: testNameFalseNotTheSameErrorCategory,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryInvalidToken,
		),
		expectedOutput: false,
		f:              f,
	})

	for _, table := range table {

		t.Run(table.name, func(t *testing.T) {
			output := table.f(table.input)
			assert.Equal(t, table.expectedOutput, output)
		})
	}
}

func TestErrorIsInvalidToken(t *testing.T) {

	f := testErrchkr.ErrorIsInvalidToken
	table := make([]testTableErrorChecker, 0)

	table = append(table, testTableErrorChecker{
		name: testNameTrue,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryInvalidToken,
		),
		expectedOutput: true,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name:           testNameFalseNotServiceError,
		input:          errMsgTest,
		expectedOutput: false,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name: testNameFalseNotTheSameErrorCategory,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryNotEnoughPermissions,
		),
		expectedOutput: false,
		f:              f,
	})

	for _, table := range table {

		t.Run(table.name, func(t *testing.T) {
			output := table.f(table.input)
			assert.Equal(t, table.expectedOutput, output)
		})
	}
}

func TestErrorIsNotEnoughPermissions(t *testing.T) {

	f := testErrchkr.ErrorIsNotEnoughPermissions
	table := make([]testTableErrorChecker, 0)

	table = append(table, testTableErrorChecker{
		name: testNameTrue,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryNotEnoughPermissions,
		),
		expectedOutput: true,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name:           testNameFalseNotServiceError,
		input:          errMsgTest,
		expectedOutput: false,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name: testNameFalseNotTheSameErrorCategory,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategorySemanticallyUnprocesable,
		),
		expectedOutput: false,
		f:              f,
	})

	for _, table := range table {

		t.Run(table.name, func(t *testing.T) {
			output := table.f(table.input)
			assert.Equal(t, table.expectedOutput, output)
		})
	}
}

func TestErrorIsSemanticallyUnprocesable(t *testing.T) {

	f := testErrchkr.ErrorIsSemanticallyUnprocesable
	table := make([]testTableErrorChecker, 0)

	table = append(table, testTableErrorChecker{
		name: testNameTrue,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategorySemanticallyUnprocesable,
		),
		expectedOutput: true,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name:           testNameFalseNotServiceError,
		input:          errMsgTest,
		expectedOutput: false,
		f:              f,
	})

	table = append(table, testTableErrorChecker{
		name: testNameFalseNotTheSameErrorCategory,
		input: svcerr.New(
			errMsgTest,
			internal.ErrorCategoryInternal,
		),
		expectedOutput: false,
		f:              f,
	})

	for _, table := range table {

		t.Run(table.name, func(t *testing.T) {
			output := table.f(table.input)
			assert.Equal(t, table.expectedOutput, output)
		})
	}
}
