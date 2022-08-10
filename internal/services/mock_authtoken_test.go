package services

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/services/internal"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"
	"github.com/pkg/errors"
)

type authTokenMock struct{}

func (s *authTokenMock) GenerateToken(userID uint, userPermissions domain.UserPermissions) ([]byte, *svcerr.ServiceError) {

	switch userID {
	case 1:
		return []byte("encryptedToken"), nil
	case 3:
		return nil, svcerr.New(
			errors.Wrap(errors.New("PayloadErr"), internal.ErrMsgGeneratingPayload.Error()),
			internal.ErrorCategoryInternal,
		)
	}

	return nil, nil
}

func (s *authTokenMock) ValidateToken(token string) (permissions interface{}, svcErr *svcerr.ServiceError) {

	switch token {
	case "tokenOK":
		permissions = string("permissionsOK")
		return
	case "tokenInvalid":
		svcErr = svcerr.New(
			errors.Wrap(errors.New("InvalidToken"), internal.ErrMsgInvalidToken.Error()),
			internal.ErrorCategoryInvalidToken,
		)
		return
	case "tokenOKButPermissionsInvalid":
		permissions = string("permissionsInvalid")
		return
	}

	return nil, nil
}
