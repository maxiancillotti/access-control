package controllers

import (
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/pkg/errors"
)

var (
	ResourcesSVCMock_RetrieveAll_ReturnSwitch int
)

type resourcesServiceMock struct{}

// May return InvalidInputID or Internal error categories.
func (*resourcesServiceMock) Create(resourceDTO dto.Resource) (*dto.Resource, error) {

	switch *resourceDTO.Path {
	case "path1":
		uintOne := uint(1)
		resourceDTO.ID = &uintOne
		return &resourceDTO, nil
	case "path0":
		return nil, errors.New("invalid_input")
	case "path2":
		return nil, errors.New("internal_error")
	}
	panic("invalid path received in mocked Create method")
}

// May return InvalidInputID or Internal error categories.
func (*resourcesServiceMock) Delete(id uint) error {
	switch id {
	case 1:
		return nil
	case 0:
		return errors.New("invalid_input")
	case 2:
		return errors.New("internal_error")
	}
	panic("invalid id received in mocked Delete method")
}

// May return InvalidInputID or Internal error categories.
func (*resourcesServiceMock) Retrieve(id uint) (*dto.Resource, error) {
	switch id {
	case 1:
		uintOne := uint(1)
		pathStr := "path"

		return &dto.Resource{
			ID:   &uintOne,
			Path: &pathStr,
		}, nil
	case 0:
		return nil, errors.New("invalid_input")
	case 2:
		return nil, errors.New("internal_error")
	}
	panic("invalid id received in mocked Retrieve method")
}

// May return InvalidInputID or Internal error categories.
func (*resourcesServiceMock) RetrieveByPath(path string) (*dto.Resource, error) {
	switch path {
	case "path1":
		uintOne := uint(1)

		return &dto.Resource{
			ID:   &uintOne,
			Path: &path,
		}, nil

	case "path0":
		return nil, errors.New("invalid_input")
	case "path2":
		return nil, errors.New("internal_error")
	}
	panic("invalid path received in mocked RetrieveByPath method")
}

// May return EmptyResult or Internal error categories.
func (*resourcesServiceMock) RetrieveAll() ([]dto.Resource, error) {
	switch ResourcesSVCMock_RetrieveAll_ReturnSwitch {
	case 0:

		uintOne := uint(1)
		pathStr := "path"

		return []dto.Resource{
			{
				ID:   &uintOne,
				Path: &pathStr,
			},
		}, nil

	case 1:
		return nil, errors.New("empty_result")
	case 2:
		return nil, errors.New("internal_error")
	}
	panic("invalid ResourcesSVCMock_RetrieveAll_ReturnSwitch set in mocked RetrieveAll method")
}
