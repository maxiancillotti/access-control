package dal

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
)

type HttpMethodsDAL interface {
	SelectAll() ([]dto.HttpMethod, error)
	SelectByName(string) (*domain.HttpMethod, error)
}
