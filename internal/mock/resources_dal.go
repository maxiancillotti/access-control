package mock

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"
	"github.com/pkg/errors"
)

var (
	ResourcesDALMock_SelectAll_ReturnSwitch = 0
)

type ResourcesDALMock struct{}

func (d *ResourcesDALMock) Insert(resource *domain.Resource) (uint, error) {

	switch resource.Path {
	case "resourceDoesNotExist":
		return 2, nil
	case "resourceDoesNotExistsReturnInsertErr":
		return 0, errors.New("error msg")
	}
	panic("incorrect resource received on Insert mocked func")
}

func (d *ResourcesDALMock) Delete(id uint) error {

	switch id {
	case 1:
		return nil
	case 4:
		return errors.New("error msg resource 4")
	}
	panic("incorrect id received on Delete mocked func")
}

func (d *ResourcesDALMock) Select(id uint) (*domain.Resource, error) {

	switch id {
	case 1:
		return &domain.Resource{ID: 1, Path: "/customers"}, nil
	case 2:
		return nil, dal.ErrorDataAccessEmptyResult
	case 3:
		return nil, errors.New("another error msg")
	}
	panic("incorrect id received on Select mocked func")
}

func (d *ResourcesDALMock) SelectByPath(path string) (*domain.Resource, error) {

	switch path {
	case "/customers":
		return &domain.Resource{ID: 1, Path: "/customers"}, nil
	case "/path/does/not/exist":
		return nil, dal.ErrorDataAccessEmptyResult
	case "/path/internal/err":
		return nil, errors.New("another error msg")
	}
	panic("incorrect path received on SelectByPath mocked func")
}

func (d *ResourcesDALMock) SelectAll() ([]dto.Resource, error) {

	switch ResourcesDALMock_SelectAll_ReturnSwitch {
	case 1:
		uintOne := uint(1)
		pathStr := "/customers"

		return []dto.Resource{
			{
				ID:   &uintOne,
				Path: &pathStr,
			},
		}, nil
	case 2:
		return nil, dal.ErrorDataAccessEmptyResult
	case 3:
		return nil, errors.New("another error msg")
	}
	panic("incorrect id received on Select mocked func")
}

func (d *ResourcesDALMock) Exists(id uint) (bool, error) {

	switch id {
	case 1:
		return true, nil
	case 2:
		return false, nil
	case 3:
		return false, errors.New("error msg")
	// returns ok now but you can use it to return err later on
	case 4:
		return true, nil
	// returns ok now but you can use it to return err later on
	case 5:
		return true, nil
	}

	panic("incorrect id received on Exists mocked func")
}

func (d *ResourcesDALMock) PathExists(path string) (bool, error) {

	switch path {
	case "resourceExists":
		return true, nil
	case "resourceDoesNotExist":
		return false, nil
	case "resourceDALerr":
		return false, errors.New("DAL error")
	case "resourceDoesNotExistsReturnInsertErr":
		return false, nil
	}
	panic("incorrect path received on PathExists mocked func")
}
