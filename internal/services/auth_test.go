package services

import (
	"testing"

	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		usrCred := &dto.UserCredentials{
			Username: "userExists",
			Password: "password",
		}

		userID, err := testAuthServices.Authenticate(usrCred)
		assert.Nil(t, err)
		assert.NotNil(t, userID)
	})
	t.Run("ErrMsgValidatingCredencials", func(t *testing.T) {
		usrCred := &dto.UserCredentials{
			Username: "userDoesNotExist",
			Password: "password",
		}

		userID, err := testAuthServices.Authenticate(usrCred)
		assert.NotNil(t, err)
		assert.Zero(t, userID)

		svcErr, ok := err.(*svcerr.ServiceError)
		assert.True(t, ok)

		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgValidatingCredencials.Error())
		assert.Equal(t, svcErr.Category(), internal.ErrorCategoryInvalidCredentials)
	})
}

func TestAuthenticateAdmin(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		usrCred := &dto.AdminCredentials{
			Username: "userExists",
			Password: "password",
		}

		userID, err := testAuthServices.AuthenticateAdmin(usrCred)
		assert.Nil(t, err)
		assert.NotNil(t, userID)
	})
	t.Run("ErrMsgValidatingCredencials", func(t *testing.T) {
		usrCred := &dto.AdminCredentials{
			Username: "userDoesNotExist",
			Password: "password",
		}

		userID, err := testAuthServices.AuthenticateAdmin(usrCred)
		assert.NotNil(t, err)
		assert.Zero(t, userID)

		svcErr, ok := err.(*svcerr.ServiceError)
		assert.True(t, ok)

		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgValidatingCredencials.Error())
		assert.Equal(t, svcErr.Category(), internal.ErrorCategoryInvalidCredentials)
	})
}

func TestCreate(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		var userID uint = 1

		token, err := testAuthServices.CreateToken(userID)

		assert.Nil(t, err)
		assert.NotNil(t, token)

		t.Log("Token created:", string(token))
	})
	t.Run("ErrMsgRetrievingPermissions", func(t *testing.T) {

		var userID uint = 2

		token, err := testAuthServices.CreateToken(userID)

		assert.NotNil(t, err)
		assert.Nil(t, token)

		svcErr, ok := err.(*svcerr.ServiceError)
		assert.True(t, ok)

		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgRetrievingPermissions)
		assert.Equal(t, svcErr.Category(), internal.ErrorCategoryInternal)
	})
	t.Run("ErrMsgGeneratingPayload", func(t *testing.T) {

		var userID uint = 3

		token, err := testAuthServices.CreateToken(userID)

		assert.NotNil(t, err)
		assert.Nil(t, token)

		svcErr, ok := err.(*svcerr.ServiceError)
		assert.True(t, ok)

		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgGeneratingPayload.Error())
		assert.Equal(t, svcErr.Category(), internal.ErrorCategoryInternal)
	})

}

func TestAuthorize(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		validationData := dto.AuthorizationRequest{
			Token:             "tokenOK",
			ResourceRequested: "resource",
			MethodRequested:   "permission",
		}

		err := testAuthServices.Authorize(validationData, "REST")
		assert.Nil(t, err)
	})

	t.Run("ErrMsgInvalidToken", func(t *testing.T) {
		validationData := dto.AuthorizationRequest{
			Token:             "tokenInvalid",
			ResourceRequested: "/resource",
			MethodRequested:   "method",
		}

		// CALL
		err := testAuthServices.Authorize(validationData, "REST")

		// TEST
		assert.NotNil(t, err)

		svcErr, ok := err.(*svcerr.ServiceError)
		assert.True(t, ok)

		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgInvalidToken.Error())
		assert.Equal(t, internal.ErrorCategoryInvalidToken, svcErr.Category())
	})

	t.Run("ErrMsgRetrievingPermissions", func(t *testing.T) {
		validationData := dto.AuthorizationRequest{
			Token:             "tokenOKButPermissionsInvalid",
			ResourceRequested: "/customers",
			MethodRequested:   "POST",
		}

		// CALL
		err := testAuthServices.Authorize(validationData, "REST")

		// TEST
		assert.NotNil(t, err)

		svcErr, ok := err.(*svcerr.ServiceError)
		assert.True(t, ok)

		assert.Contains(t, svcErr.ErrorValue().Error(), internal.ErrMsgRetrievingPermissions)
		assert.Equal(t, internal.ErrorCategoryInternal, svcErr.Category())
	})
}
