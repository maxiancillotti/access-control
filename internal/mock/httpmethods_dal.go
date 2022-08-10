package mock

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"
	"github.com/pkg/errors"
)

var (
	HttpmethodsDALMock_SelectAll_ReturnSwitch = 0
)

type HttpmethodsDALMock struct{}

func (d *HttpmethodsDALMock) SelectByName(name string) (*domain.HttpMethod, error) {
	switch name {
	case "methodExists":
		return &domain.HttpMethod{
			ID:   1,
			Name: name,
		}, nil
	case "methodDoesNotExist":
		return nil, dal.ErrorDataAccessEmptyResult
	case "methodDALerr":
		return nil, errors.New("DAL error")
	}
	panic("incorrect name received on SelectByName mocked func")
}

func (d *HttpmethodsDALMock) SelectAll() ([]dto.HttpMethod, error) {

	switch HttpmethodsDALMock_SelectAll_ReturnSwitch {
	case 1:
		uintOne := uint(1)
		methodGET := "GET"

		return []dto.HttpMethod{
			{
				ID:   &uintOne,
				Name: &methodGET,
			},
		}, nil
	case 2:
		return nil, dal.ErrorDataAccessEmptyResult
	case 3:
		return nil, errors.New("another error msg")
	}
	panic("incorrect id received on Select mocked func")
}
