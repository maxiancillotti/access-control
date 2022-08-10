package controllers

import (
	"github.com/maxiancillotti/access-control/internal/dto"

	"github.com/pkg/errors"
)

var (
	HttpMethodsSVCMock_RetrieveAll_ReturnSwitch int
)

type httpMethodsServiceMock struct{}

// May return EmptyResult or Internal error categories.
func (*httpMethodsServiceMock) RetrieveAll() ([]dto.HttpMethod, error) {
	switch ResourcesSVCMock_RetrieveAll_ReturnSwitch {
	case 0:

		uintOne := uint(1)
		methodNameStr := "methodName"

		return []dto.HttpMethod{
			{
				ID:   &uintOne,
				Name: &methodNameStr,
			},
		}, nil

	case 1:
		return nil, errors.New("empty_result")
	case 2:
		return nil, errors.New("internal_error")
	}
	panic("invalid HttpMethodsSVCMock_RetrieveAll_ReturnSwitch set in mocked RetrieveAll method")
}

// May return InvalidInput or Internal error categories.
func (*httpMethodsServiceMock) RetrieveByName(name string) (*dto.HttpMethod, error) {
	switch name {
	case "methodName1":
		uintOne := uint(1)

		return &dto.HttpMethod{
			ID:   &uintOne,
			Name: &name,
		}, nil

	case "methodName0":
		return nil, errors.New("invalid_input")
	case "methodName2":
		return nil, errors.New("internal_error")
	}
	panic("invalid name received in mocked RetrieveByName method")
}
