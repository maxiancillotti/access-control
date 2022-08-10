package errchk

import (
	"github.com/maxiancillotti/access-control/internal/services"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"
)

func NewServiceErrorChecker() services.ServiceErrorChecker {
	return &errorChecker{}
}

type errorChecker struct{}

func (c *errorChecker) ErrorIsInvalidInputIdentifier(err error) bool {
	svcErr, ok := err.(*svcerr.ServiceError)
	return ok && svcErr.Category() == internal.ErrorCategoryInvalidInputID
}

func (c *errorChecker) ErrorIsEmptyResult(err error) bool {
	svcErr, ok := err.(*svcerr.ServiceError)
	return ok && svcErr.Category() == internal.ErrorCategoryEmptyResult
}

func (c *errorChecker) ErrorIsInternal(err error) bool {
	svcErr, ok := err.(*svcerr.ServiceError)
	return ok && svcErr.Category() == internal.ErrorCategoryInternal
}

func (c *errorChecker) ErrorIsInvalidCredentials(err error) bool {
	svcErr, ok := err.(*svcerr.ServiceError)
	return ok && svcErr.Category() == internal.ErrorCategoryInvalidCredentials
}

func (c *errorChecker) ErrorIsInvalidToken(err error) bool {
	svcErr, ok := err.(*svcerr.ServiceError)
	return ok && svcErr.Category() == internal.ErrorCategoryInvalidToken
}

func (c *errorChecker) ErrorIsNotEnoughPermissions(err error) bool {
	svcErr, ok := err.(*svcerr.ServiceError)
	return ok && svcErr.Category() == internal.ErrorCategoryNotEnoughPermissions
}

func (c *errorChecker) ErrorIsSemanticallyUnprocesable(err error) bool {
	svcErr, ok := err.(*svcerr.ServiceError)
	return ok && svcErr.Category() == internal.ErrorCategorySemanticallyUnprocesable
}
