package internal

import (
	"github.com/pkg/errors"
)

var (
	/******************************************************************************************/
	// General

	ErrRespInternalUnexpected = errors.New("There's been an unexpected error, please try again, we will have this solved soon")

	ErrMsgFmtParsingIDFromURL = "error parsing %s id from URL"

	ErrMsgBlockUnmarshalingReqBody        = "error unmarshaling request body"
	ErrMsgBlockReqBodyDoesntHaveProperFmt = "request body doesn't have the proper format"

	/******************************************************************************************/
	// CRUD ERRORS

	// Client Input Errors

	ErrMsgFmtCreating = "error creating %s"
	ErrMsgFmtUpdating = "error updating %s"
	ErrMsgFmtDeleting = "error deleting %s"
	ErrMsgFmtGetting  = "error getting %s"

	////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////

	/******************************************************************************************/
	// AuthorizationBasic & AuthorizationXAdmin

	ErrMsgFmtAuthorization = "error on authentication (Authorization %s)" // add Authorization header type: "Basic", "X-Admin", etc.

	ErrLogAuthorizationHeaderIsEmpty = errors.New("authorization header is empty")

	ErrLogFmtAuthorizationMustBeOfType                   = "authorization must be of type %s"
	ErrLogFmtAuthorizationCannotFindCredentialsSeparator = "authorization %s failed, cannot find credentials separator"
	ErrLogFmtAuthorizationCannotParseCredentials         = "authorization %s failed, cannot parse credentials"
	ErrLogFmtAuthorizationCredentialsCannotBeEmpty       = "authorization %s failed, credentials cannot be empty"

	ErrLogBlockFmtAuthorizationCannotBase64Decode            = "authorization %s failed, cannot base64 decode credentials"
	ErrLogBlockFmtAuthorizationInvalidCredentialsCauseFormat = "authorization %s failed, invalid credentials format"
	ErrLogBlockFmtAuthorizationInvalidCredentials            = "authorization %s failed, invalid credentials"
	ErrLogBlockFmtAuthorizationService                       = "authorization %s failed, error returned from authentication service"

	/******************************************************************************************/
	// PostAuthToken. Token Creation.

	ErrLogRetrievingUserIDFromContext = errors.New("error retrieving User ID from request context")

	ErrMsgBlockCreatingToken = "error returned from token creation service"

	/******************************************************************************************/
	// PostAuthTokenValidate. Token Validation.

	ErrRespCannotProcessTokenCorrectly = errors.New("cannot process token correctly")
	ErrRespNotEnoughPermissions        = errors.New("user doesn't have enough permissions for the requested resource")
	ErrRespInvalidToken                = errors.New("token is invalid")

	ErrMsgBlockAuthorizingAccess = "error authorizing access"
)
