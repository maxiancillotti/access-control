package services

import (
	"testing"

	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/mock"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestHttpMethodsRetrieveAll(t *testing.T) {

	type testCase struct {
		name                  string
		input                 int
		expectedMethodsOutput []dto.HttpMethod
		expectedErrOutput     *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	uintOne := uint(1)
	methodGET := "GET"
	table = append(table, testCase{
		name:  "Success",
		input: 1,
		expectedMethodsOutput: []dto.HttpMethod{
			{
				ID:   &uintOne,
				Name: &methodGET,
			},
		},
		expectedErrOutput: nil,
	})

	table = append(table, testCase{
		name:                  "Error. There aren't any http methods.",
		input:                 2,
		expectedMethodsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtEmptyResult, "http methods"),
			internal.ErrorCategoryEmptyResult,
		),
	})

	table = append(table, testCase{
		name:                  "Error. Http methods retrieval failed. Internal.",
		input:                 3,
		expectedMethodsOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("another error msg"), internal.ErrMsgFmtRetrievalFailed, "http methods"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			mock.HttpmethodsDALMock_SelectAll_ReturnSwitch = test.input

			httpMethods, err := testHttpMethodsSvc.RetrieveAll()

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedMethodsOutput[0].Name, httpMethods[0].Name)
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

func TestHttpMethodsRetrieveByName(t *testing.T) {

	type testCase struct {
		name                 string
		input                string
		expectedMethodOutput *dto.HttpMethod
		expectedErrOutput    *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	uintOne := uint(1)
	methodNameExists := "methodExists"
	table = append(table, testCase{
		name:                 "Success",
		input:                "methodExists",
		expectedMethodOutput: &dto.HttpMethod{ID: &uintOne, Name: &methodNameExists},
		expectedErrOutput:    nil,
	})

	table = append(table, testCase{
		name:                 "Error. Name does not exist. Invalid input ID.",
		input:                "methodDoesNotExist",
		expectedMethodOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "http method"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                 "Error. Http method retrieval failed. Internal.",
		input:                "methodDALerr",
		expectedMethodOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("DAL error"), internal.ErrMsgFmtRetrievalFailed, "http method"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			method, err := testHttpMethodsSvc.RetrieveByName(test.input)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedMethodOutput.Name, method.Name)
			} else {
				svcErr, ok := err.(*svcerr.ServiceError)
				assert.True(t, ok)

				assert.NotNil(t, err)

				assert.Equal(t, test.expectedErrOutput.ErrorValue().Error(), svcErr.ErrorValue().Error())
				assert.Equal(t, test.expectedErrOutput.Category(), svcErr.Category())
			}
		})
	}
}
