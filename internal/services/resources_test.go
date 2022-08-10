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

func TestResourcesCreate(t *testing.T) {

	type testCase struct {
		name                      string
		input                     dto.Resource
		expectedResourceOutput    *dto.Resource
		expectedPasswordLenOutput int
		expectedErrOutput         *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	inputResourceDoesNotExist := "resourceDoesNotExist"
	inputResourceExists := "resourceExists"
	inputResourceDALerr := "resourceDALerr"
	inputResourceDoesNotExistsReturnInsertErr := "resourceDoesNotExistsReturnInsertErr"

	uintTwo := uint(2)
	table = append(table, testCase{
		name:                      "Success",
		input:                     dto.Resource{Path: &inputResourceDoesNotExist},
		expectedResourceOutput:    &dto.Resource{ID: &uintTwo, Path: &inputResourceDoesNotExist},
		expectedPasswordLenOutput: 64,
		expectedErrOutput:         nil,
	})

	table = append(table, testCase{
		name:                   "Error. Resource already exists. Invalid input ID.",
		input:                  dto.Resource{Path: &inputResourceExists},
		expectedResourceOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtAlreadyExists, "resource"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                   "Error. Failed to check if resource exists. Internal.",
		input:                  dto.Resource{Path: &inputResourceDALerr},
		expectedResourceOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("DAL error"), internal.ErrMsgFmtFailedToCheckIfExists, "resource"),
			internal.ErrorCategoryInternal,
		),
	})

	table = append(table, testCase{
		name:                   "Error. Insert failed. Internal.",
		input:                  dto.Resource{Path: &inputResourceDoesNotExistsReturnInsertErr},
		expectedResourceOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg"), internal.ErrMsgFmtInsertFailed, "resource"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			resource, err := testResourcesSvc.Create(test.input)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedResourceOutput.ID, resource.ID)
				assert.Equal(t, test.expectedResourceOutput.Path, resource.Path)
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

func TestResourcesDelete(t *testing.T) {

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
		name:  "Error from existsOrErr. Resource does not exist",
		input: 2,
		expectedOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "resource"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:  "Error. Internal. Delete failed",
		input: 4,
		expectedOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg resource 4"), internal.ErrMsgFmtDeleteFailed, "resource"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			err := testResourcesSvc.Delete(test.input)

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

func TestResourcesExistsOrErr(t *testing.T) {

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
		name:  "Error. Resource does not exist.",
		input: 2,
		expectedOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "resource"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:  "Error. Failed to check if resource exist",
		input: 3,
		expectedOutput: svcerr.New(
			errors.Wrapf(errors.New("error msg"), internal.ErrMsgFmtFailedToCheckIfExists, "resource"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {
			testRscsSvcs := testResourcesSvc.(*resourcesInteractor)
			svcErr := testRscsSvcs.existsOrErr(test.input)

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

func TestResourcesRetrieve(t *testing.T) {

	type testCase struct {
		name                   string
		input                  uint
		expectedResourceOutput *dto.Resource
		expectedErrOutput      *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	uintOne := uint(1)
	pathCustomers := "/customers"
	table = append(table, testCase{
		name:                   "Success",
		input:                  1,
		expectedResourceOutput: &dto.Resource{ID: &uintOne, Path: &pathCustomers},
		expectedErrOutput:      nil,
	})

	table = append(table, testCase{
		name:                   "Error. Resource does not exist. Invalid input ID.",
		input:                  2,
		expectedResourceOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "resource"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                   "Error. Resource retrieval failed. Internal.",
		input:                  3,
		expectedResourceOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("another error msg"), internal.ErrMsgFmtRetrievalFailed, "resource"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			resource, err := testResourcesSvc.Retrieve(test.input)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedResourceOutput.Path, resource.Path)
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

func TestResourcesRetrieveByPath(t *testing.T) {

	type testCase struct {
		name                   string
		input                  string
		expectedResourceOutput *dto.Resource
		expectedErrOutput      *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	uintOne := uint(1)
	pathCustomers := "/customers"
	table = append(table, testCase{
		name:                   "Success",
		input:                  "/customers",
		expectedResourceOutput: &dto.Resource{ID: &uintOne, Path: &pathCustomers},
		expectedErrOutput:      nil,
	})

	table = append(table, testCase{
		name:                   "Error. Resource does not exist. Invalid input ID.",
		input:                  "/path/does/not/exist",
		expectedResourceOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "resource"),
			internal.ErrorCategoryInvalidInputID,
		),
	})

	table = append(table, testCase{
		name:                   "Error. Resource retrieval failed. Internal.",
		input:                  "/path/internal/err",
		expectedResourceOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("another error msg"), internal.ErrMsgFmtRetrievalFailed, "resource"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			resource, err := testResourcesSvc.RetrieveByPath(test.input)

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedResourceOutput.Path, resource.Path)
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

func TestResourcesRetrieveAll(t *testing.T) {

	type testCase struct {
		name                    string
		input                   int
		expectedResourcesOutput []dto.Resource
		expectedErrOutput       *svcerr.ServiceError
	}

	table := make([]testCase, 0)

	uintOne := uint(1)
	pathCustomers := "/customers"
	table = append(table, testCase{
		name:  "Success",
		input: 1,
		expectedResourcesOutput: []dto.Resource{
			{
				ID:   &uintOne,
				Path: &pathCustomers,
			},
		},
		expectedErrOutput: nil,
	})

	table = append(table, testCase{
		name:                    "Error. There aren't any resources.",
		input:                   2,
		expectedResourcesOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Errorf(internal.ErrMsgFmtEmptyResult, "resources"),
			internal.ErrorCategoryEmptyResult,
		),
	})

	table = append(table, testCase{
		name:                    "Error. Resources retrieval failed. Internal.",
		input:                   3,
		expectedResourcesOutput: nil,
		expectedErrOutput: svcerr.New(
			errors.Wrapf(errors.New("another error msg"), internal.ErrMsgFmtRetrievalFailed, "resources"),
			internal.ErrorCategoryInternal,
		),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			mock.ResourcesDALMock_SelectAll_ReturnSwitch = test.input

			resources, err := testResourcesSvc.RetrieveAll()

			if test.expectedErrOutput == nil {
				assert.Nil(t, err)
				assert.Equal(t, test.expectedResourcesOutput[0].Path, resources[0].Path)
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
