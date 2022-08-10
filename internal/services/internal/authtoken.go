package internal

import (
	"github.com/maxiancillotti/access-control/internal/domain"
	"github.com/maxiancillotti/access-control/internal/services/svcerr"
)

type AuthTokenServices interface {
	// Creates authorization token with the given permissions associated to the userID.
	GenerateToken(userID uint, userPermissions domain.UserPermissions) ([]byte, *svcerr.ServiceError)

	// Checks the data contained into the token is safe and valid to authorize a
	// request from the holder.
	// Also returns the users permissions originally set, so the caller can check them
	// after parsing/asserting and go on with the authorization.
	ValidateToken(token string) (permissions interface{}, svcErr *svcerr.ServiceError)
}
