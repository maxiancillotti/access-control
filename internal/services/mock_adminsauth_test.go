package services

import (
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/pkg/errors"
)

type adminsAuthInteractorMock struct{}

// Returns AdminID > 0 and error == nil when succesful
func (as *adminsAuthInteractorMock) validateAdminCredentials(adminCredentials *dto.AdminCredentials) (adminID uint, svcErr *svcerr.ServiceError) {

	switch adminCredentials.Username {
	case "userExists":
		adminID = 1
		return
	default:
		svcErr = svcerr.New(
			errors.New("internalErr"),
			internal.ErrorCategoryInvalidCredentials,
		)
	}
	return
}
