package dal

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
)

type ResourcesDAL interface {
	// Inserts a resource and returns the created ID
	Insert(resource *domain.Resource) (uint, error)
	Delete(id uint) error
	Select(id uint) (*domain.Resource, error)
	SelectByPath(Path string) (*domain.Resource, error)
	SelectAll() ([]dto.Resource, error)
	Exists(id uint) (bool, error)
	PathExists(path string) (bool, error)
}
