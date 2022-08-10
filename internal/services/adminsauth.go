package services

import (
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/maxiancillotti/passwords"
	"github.com/pkg/errors"
)

func newAdminsAuthServices(admSvc AdminsServices) adminsAuthServices {
	return &adminsAuthInteractor{adminsSvc: admSvc}
}

type adminsAuthServices interface {
	validateAdminCredentials(adminCredentials *dto.AdminCredentials) (adminID uint, svcErr *svcerr.ServiceError)
}

type adminsAuthInteractor struct {
	adminsSvc AdminsServices
}

// Returns AdminID > 0 and error == nil when succesful
func (aai *adminsAuthInteractor) validateAdminCredentials(adminCredentials *dto.AdminCredentials) (adminID uint, svcErr *svcerr.ServiceError) {

	admin, svcErr := aai.adminsSvc.retrieveByUsername(adminCredentials.Username)
	if svcErr != nil {
		if svcErr.Category() == internal.ErrorCategoryInvalidInputID {
			svcErr = svcerr.New(
				errors.Wrap(svcErr.ErrorValue(), internal.ErrMsgInvalidUsername),
				internal.ErrorCategoryInvalidCredentials,
			)
		}
		return
	}

	if !admin.Enabled {
		svcErr = svcerr.New(
			errors.New(internal.ErrMsgUserDisabled),
			//internal.ErrorCategoryUserDisabled,
			internal.ErrorCategoryInvalidCredentials,
		)
		return
	}

	//fmt.Printf("plain pw: '%s'\n", adminCredentials.Password)
	//fmt.Printf("pw salt: '%s'\n", admin.PasswordSalt)
	//fmt.Printf("pw stored hash: '%s'\n", admin.PasswordHash)

	err := passwords.ValidateHashPw(adminCredentials.Password, admin.PasswordSalt, admin.PasswordHash)
	if err != nil {
		svcErr = svcerr.New(
			errors.Wrap(err, internal.ErrMsgInvalidPassword),
			internal.ErrorCategoryInvalidCredentials,
		)
		return
	}
	adminID = admin.ID
	return
}
