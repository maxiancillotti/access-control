package services

import (
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/dal"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
)

func NewHttpMethodsServices(dal dal.HttpMethodsDAL) HttpMethodsServices {
	return &httpMethodsInteractor{dal: dal}
}

type HttpMethodsServices interface {

	// May return EmptyResult or Internal error categories.
	RetrieveAll() ([]dto.HttpMethod, error)

	// May return InvalidInput or Internal error categories.
	RetrieveByName(name string) (*dto.HttpMethod, error)
}

type httpMethodsInteractor struct {
	dal dal.HttpMethodsDAL
}

func (mi *httpMethodsInteractor) RetrieveAll() ([]dto.HttpMethod, error) {
	methods, err := mi.dal.SelectAll()
	if err != nil {
		if errors.Is(err, dal.ErrorDataAccessEmptyResult) {
			return nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtEmptyResult, "http methods"),
				internal.ErrorCategoryEmptyResult,
			)
		}
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "http methods"),
			internal.ErrorCategoryInternal,
		)
	}
	return methods, nil
}

func (mi *httpMethodsInteractor) RetrieveByName(name string) (*dto.HttpMethod, error) {
	method, err := mi.dal.SelectByName(name)
	if err != nil {
		if errors.Is(err, dal.ErrorDataAccessEmptyResult) {
			return nil, svcerr.New(
				errors.Errorf(internal.ErrMsgFmtDoesNotExist, "http method"),
				internal.ErrorCategoryInvalidInputID,
			)
		}
		return nil, svcerr.New(
			errors.Wrapf(err, internal.ErrMsgFmtRetrievalFailed, "http method"),
			internal.ErrorCategoryInternal,
		)

	}
	return &dto.HttpMethod{ID: &method.ID, Name: &method.Name}, nil
}
