package services

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"
	"github.com/pkg/errors"
)

func NewResourcesServices(dal dal.ResourcesDAL) ResourcesServices {
	return &resourcesInteractor{dal: dal}
}

type ResourcesServices interface {

	// May return InvalidInputID or Internal error categories.
	Create(resourceDTO dto.Resource) (*dto.Resource, error)

	// May return InvalidInputID or Internal error categories.
	Delete(id uint) error

	// May return InvalidInputID or Internal error categories.
	Retrieve(id uint) (*dto.Resource, error)

	// May return InvalidInputID or Internal error categories.
	RetrieveByPath(path string) (*dto.Resource, error)

	// May return EmptyResult or Internal error categories.
	RetrieveAll() ([]dto.Resource, error)
}

type resourcesInteractor struct {
	dal dal.ResourcesDAL
}

func (ri *resourcesInteractor) Create(resourceDTO dto.Resource) (*dto.Resource, error) {

	exists, err := ri.dal.PathExists(*resourceDTO.Path)
	if err != nil {
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtFailedToCheckIfExists, "resource"),
			internal.ErrorCategoryInternal,
		)
	} else if exists {
		return nil, svcerr.New(
			errors.Errorf(internal.ErrMsgFmtAlreadyExists, "resource"),
			internal.ErrorCategoryInvalidInputID,
		)
	}

	resource := &domain.Resource{Path: *resourceDTO.Path}

	resource.ID, err = ri.dal.Insert(resource)
	if err != nil {
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtInsertFailed, "resource"),
			internal.ErrorCategoryInternal,
		)
	}
	return &dto.Resource{ID: &resource.ID, Path: &resource.Path}, nil
}

func (ri *resourcesInteractor) Delete(id uint) error {

	svcErr := ri.existsOrErr(id)
	if svcErr != nil {
		return svcErr
	}

	err := ri.dal.Delete(id)
	if err != nil {
		return svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtDeleteFailed, "resource"),
			internal.ErrorCategoryInternal,
		)
	}
	return nil
}

func (ri *resourcesInteractor) existsOrErr(id uint) (svcErr *svcerr.ServiceError) {
	exists, err := ri.dal.Exists(id)
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtFailedToCheckIfExists, "resource"),
			internal.ErrorCategoryInternal,
		)
	} else if !exists {
		svcErr = svcerr.New(
			errors.Errorf(internal.ErrMsgFmtDoesNotExist, "resource"),
			internal.ErrorCategoryInvalidInputID,
		)
	}
	return
}

func (ri *resourcesInteractor) Retrieve(id uint) (*dto.Resource, error) {
	resource, err := ri.dal.Select(id)
	if err != nil {
		if errors.Is(err, dal.ErrorDataAccessEmptyResult) {
			return nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtDoesNotExist, "resource"),
				internal.ErrorCategoryInvalidInputID,
			)
		}
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "resource"),
			internal.ErrorCategoryInternal,
		)
	}
	return &dto.Resource{
		ID:   &resource.ID,
		Path: &resource.Path,
	}, nil
}

func (ri *resourcesInteractor) RetrieveByPath(path string) (*dto.Resource, error) {
	resource, err := ri.dal.SelectByPath(path)
	if err != nil {
		if errors.Is(err, dal.ErrorDataAccessEmptyResult) {
			return nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtDoesNotExist, "resource"),
				internal.ErrorCategoryInvalidInputID,
			)
		}
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "resource"),
			internal.ErrorCategoryInternal,
		)
	}
	return &dto.Resource{
		ID:   &resource.ID,
		Path: &resource.Path,
	}, nil
}

func (ri *resourcesInteractor) RetrieveAll() ([]dto.Resource, error) {
	resources, err := ri.dal.SelectAll()
	if err != nil {
		if errors.Is(err, dal.ErrorDataAccessEmptyResult) {
			return nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtEmptyResult, "resources"),
				internal.ErrorCategoryEmptyResult,
			)
		}
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "resources"),
			internal.ErrorCategoryInternal,
		)
	}
	return resources, nil
}
