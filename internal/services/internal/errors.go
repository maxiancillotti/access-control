package internal

import (
	"github.com/maxiancillotti/access-control/internal/services/svcerr"
	"github.com/pkg/errors"
)

var (
	ErrorCategoryInternal svcerr.ServiceErrorCategory = "Internal error"

	/******************************************************************************************/
	// CRUD ERRORS

	// Client Input Errors

	// When retrieving or writting using an ID.
	// Can return the full error message safely. Does not containg private info.
	ErrorCategoryInvalidInputID svcerr.ServiceErrorCategory = "Invalid Input Identifier"

	ErrMsgFmtAlreadyExists = "%s already exists"
	ErrMsgFmtDoesNotExist  = "%s does not exist"

	//

	// When retrieving using something that is not an ID as an input.
	// Can return the full error message safely. Does not containg private info.
	ErrorCategoryEmptyResult svcerr.ServiceErrorCategory = "Empty result"

	ErrMsgFmtEmptyResult              = "there are no %s registered"
	ErrMsgFmtEmptyResultForTheGivenID = "there are no %s registered for the given %s"

	// Internal Server Errors
	ErrMsgFmtFailedToCheckIfExists = "failed to check if %s already exists"
	ErrMsgFmtInsertFailed          = "%s insert failed"
	ErrMsgFmtUpdateFailed          = "%s update failed"
	ErrMsgFmtDeleteFailed          = "%s delete failed"
	ErrMsgFmtRetrievalFailed       = "%s retrieval failed"

	/******************************************************************************************/
	// User errors

	// Internal Server Errors

	ErrMsgPwHashingFailed = "password hashing failed"

	////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////////////////

	/******************************************************************************************/
	// Authenticate errors

	// Client Input Errors

	ErrorCategoryInvalidCredentials svcerr.ServiceErrorCategory = "Invalid credentials"

	ErrMsgInvalidUsername = "invalid username"
	ErrMsgInvalidPassword = "invalid password"

	// no need to use another err category
	//ErrorCategoryUserDisabled svcerr.ServiceErrorCategory = "User disabled"
	ErrMsgUserDisabled = "user is disabled"

	/******************************************************************************************/
	// CreateToken errors

	ErrMsgRetrievingPermissions        = "error retrieving user permissions"
	ErrMsgUserDoesntHaveAnyPermissions = "user doesn't have any permissions"

	// Internal Server Errors

	ErrMsgGeneratingPayload      = errors.New("error generating payload")
	ErrMsgSignFailed             = errors.New("token sign failed")
	ErrMsgFailedToEncryptPayload = errors.New("failed to encrypt payload")

	/******************************************************************************************/
	// Authorize errors

	ErrorCategoryInvalidToken             svcerr.ServiceErrorCategory = "Invalid token"
	ErrorCategoryNotEnoughPermissions     svcerr.ServiceErrorCategory = "Not enough permissions"
	ErrorCategorySemanticallyUnprocesable svcerr.ServiceErrorCategory = "Semantically unprocesable"

	// Internal Server Errors

	ErrMsgFailedToDecryptToken           = errors.New("failed to decrypt token")
	ErrMsgInvalidToken                   = errors.New("invalid token")
	ErrMsgTokenDoesntClaimAnyPermissions = errors.New("token doesn't claim any permissions")

	ErrMsgValidatingCredencials = errors.New("error validating credentials")
)
