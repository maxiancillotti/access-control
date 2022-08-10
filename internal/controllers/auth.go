package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/maxiancillotti/access-control/internal/controllers/internal"
	"github.com/maxiancillotti/access-control/internal/controllers/presenter"
	"github.com/maxiancillotti/access-control/internal/controllers/requestutil"
	"github.com/maxiancillotti/access-control/internal/controllers/responseutil"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/services"

	"github.com/pkg/errors"
)

type AuthController interface {

	// Authenticates an admin vía http header "Authorization"
	// with the custom type "X-Admin". Credentials formating
	// is equal to "Basic" authentication.
	AuthorizationXAdmin(next http.HandlerFunc) http.HandlerFunc

	// Authenticates a user vía http header "Authorization"
	// with type "Basic"
	AuthorizationBasic(next http.HandlerFunc) http.HandlerFunc

	// Creates a new access authorization token for a user
	// authenticated vía Authorization Basic.
	PostAuthToken(resp http.ResponseWriter, req *http.Request)

	// Authorices access to a resource for a given token.
	// Receives dto.AuthorizationRequest body with fields:
	/*
		{
		    "token": string,
		    "resource_requested": string,
		    "method_requested": string
		}
	*/
	PostAuthTokenAuthorize(resp http.ResponseWriter, req *http.Request)
}

type authController struct {
	authServices        services.AuthServices
	serviceErrorChecker services.ServiceErrorChecker
	presenter           presenter.Presenter
	writerLogger        responseutil.WriterLogger
}

func NewAuthController(authService services.AuthServices, svcErrorChecker services.ServiceErrorChecker, presenter presenter.Presenter, writerLogger responseutil.WriterLogger) AuthController {
	return &authController{
		authServices:        authService,
		serviceErrorChecker: svcErrorChecker,
		presenter:           presenter,
		writerLogger:        writerLogger,
	}
}

func (c *authController) AuthorizationXAdmin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		expectedHeaderType := "X-Admin"

		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, internal.ErrLogAuthorizationHeaderIsEmpty, authHeader, expectedHeaderType)
			return
		}

		splitAuthHeader := strings.Split(authHeader, fmt.Sprintf("%s ", expectedHeaderType))
		if len(splitAuthHeader) < 2 {
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationMustBeOfType, expectedHeaderType), authHeader, expectedHeaderType)
			return
		}
		authHash := splitAuthHeader[1]

		decodedCredentialsBytes, err := base64.StdEncoding.DecodeString(authHash)
		if err != nil {
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Wrapf(err, internal.ErrLogBlockFmtAuthorizationCannotBase64Decode, expectedHeaderType), authHeader, expectedHeaderType)
			return
		}
		decodedCredentials := string(decodedCredentialsBytes)

		// OLD
		/*
			if !strings.ContainsAny(decodedCredentials, ":") {
				c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationCannotFindCredentialsSeparator, expectedHeaderType), authHeader, expectedHeaderType)
				return
			}
			splitCredentials := strings.Split(decodedCredentials, ":")

			// TODO: check if this is necessary. See test comments.
			if len(splitCredentials) < 2 {
				c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationCannotParseCredentials, expectedHeaderType), authHeader, expectedHeaderType)
				return
			}

			if splitCredentials[0] == "" || splitCredentials[1] == "" {
				c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationCredentialsCannotBeEmpty, expectedHeaderType), authHeader, expectedHeaderType)
				return
			}

			adminCredentials := dto.AdminCredentials{
				Username: splitCredentials[0],
				Password: splitCredentials[1],
			}
		*/

		// NEW
		colonIndex := strings.IndexRune(decodedCredentials, ':')
		if colonIndex == -1 {
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationCannotFindCredentialsSeparator, expectedHeaderType), authHeader, expectedHeaderType)
			return
		}

		username := decodedCredentials[0:colonIndex]
		password := decodedCredentials[colonIndex+1:]

		if username == "" || password == "" {
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationCredentialsCannotBeEmpty, expectedHeaderType), authHeader, expectedHeaderType)
			return
		}

		adminCredentials := dto.AdminCredentials{
			Username: username,
			Password: password,
		}

		//////////////////////////

		err = adminCredentials.ValidateFormat()
		if err != nil {
			// Don't give too much info to the caller, like why credentials are invalid.
			// Status 400 or 422 would fit the error, but this would be too much info.
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Wrapf(err, internal.ErrLogBlockFmtAuthorizationInvalidCredentialsCauseFormat, expectedHeaderType), authHeader, expectedHeaderType)
			return
		}

		//adminID
		_, err = c.authServices.AuthenticateAdmin(&adminCredentials)
		if err != nil {
			if c.serviceErrorChecker.ErrorIsInvalidCredentials(err) {
				c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Wrapf(err, internal.ErrLogBlockFmtAuthorizationInvalidCredentials, expectedHeaderType), authHeader, expectedHeaderType)
			} else {
				c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusInternalServerError, internal.ErrRespInternalUnexpected, errors.Wrapf(err, internal.ErrLogBlockFmtAuthorizationService, expectedHeaderType))
			}
			return
		}
		/*
			ctx := req.Context()
			ctx = requestutil.AssignAuthenticatedUserIDToCxt(ctx, adminID)
			//*req = *req.Clone(ctx)
			req = req.Clone(ctx)
		*/
		next.ServeHTTP(rw, req)
	})
}

func (c *authController) AuthorizationBasic(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		expectedHeaderType := "Basic"

		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, internal.ErrLogAuthorizationHeaderIsEmpty, authHeader, expectedHeaderType)
			return
		}

		splitAuthHeader := strings.Split(authHeader, fmt.Sprintf("%s ", expectedHeaderType))
		if len(splitAuthHeader) < 2 {
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationMustBeOfType, expectedHeaderType), authHeader, expectedHeaderType)
			return
		}
		authHash := splitAuthHeader[1]

		decodedCredentialsBytes, err := base64.StdEncoding.DecodeString(authHash)
		if err != nil {
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Wrapf(err, internal.ErrLogBlockFmtAuthorizationCannotBase64Decode, expectedHeaderType), authHeader, expectedHeaderType)
			return
		}
		decodedCredentials := string(decodedCredentialsBytes)

		// OLD
		/*
			if !strings.ContainsAny(decodedCredentials, ":") {
				c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationCannotFindCredentialsSeparator, expectedHeaderType), authHeader, expectedHeaderType)
				return
			}
			splitCredentials := strings.Split(decodedCredentials, ":")
			// TODO: check if this is necessary. See test comments.
			if len(splitCredentials) < 2 {
				c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationCannotParseCredentials, expectedHeaderType), authHeader, expectedHeaderType)
				return
			}
			if splitCredentials[0] == "" || splitCredentials[1] == "" {
				c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationCredentialsCannotBeEmpty, expectedHeaderType), authHeader, expectedHeaderType)
				return
			}

			userCredentials := dto.UserCredentials{
				Username: splitCredentials[0],
				Password: splitCredentials[1],
			}
		*/

		// NEW
		colonIndex := strings.IndexRune(decodedCredentials, ':')
		if colonIndex == -1 {
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationCannotFindCredentialsSeparator, expectedHeaderType), authHeader, expectedHeaderType)
			return
		}

		username := decodedCredentials[0:colonIndex]
		password := decodedCredentials[colonIndex+1:]

		if username == "" || password == "" {
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Errorf(internal.ErrLogFmtAuthorizationCredentialsCannotBeEmpty, expectedHeaderType), authHeader, expectedHeaderType)
			return
		}

		userCredentials := dto.UserCredentials{
			Username: username,
			Password: password,
		}

		//////////////////////////

		err = userCredentials.ValidateFormat()
		if err != nil {
			// Don't give too much info to the caller, like why credentials are invalid.
			// Status 400 or 422 would fit the error, but this would be too much info.
			c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Wrapf(err, internal.ErrLogBlockFmtAuthorizationInvalidCredentialsCauseFormat, expectedHeaderType), authHeader, expectedHeaderType)
			return
		}

		userID, err := c.authServices.Authenticate(&userCredentials)
		if err != nil {
			if c.serviceErrorChecker.ErrorIsInvalidCredentials(err) {
				c.writerLogger.WriteUnauthorizedErrorRespAndLog(req.Context(), rw, errors.Wrapf(err, internal.ErrLogBlockFmtAuthorizationInvalidCredentials, expectedHeaderType), authHeader, expectedHeaderType)
			} else {
				c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusInternalServerError, internal.ErrRespInternalUnexpected, errors.Wrapf(err, internal.ErrLogBlockFmtAuthorizationService, expectedHeaderType))
			}
			return
		}

		ctx := req.Context()
		ctx = requestutil.AssignAuthenticatedUserIDToCxt(ctx, userID)
		//*req = *req.Clone(ctx)
		req = req.Clone(ctx)

		next.ServeHTTP(rw, req)
	})
}

func (c *authController) PostAuthToken(rw http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	userID := requestutil.GetAuthenticatedUserIDFromCtx(ctx)
	if userID == 0 {
		c.writerLogger.WriteErrorRespAndLog(ctx, rw, http.StatusInternalServerError, internal.ErrRespInternalUnexpected, internal.ErrLogRetrievingUserIDFromContext)
		return
	}

	token, err := c.authServices.CreateToken(userID)
	if err != nil {
		c.writerLogger.WriteErrorRespAndLog(ctx, rw, http.StatusInternalServerError, internal.ErrRespInternalUnexpected, errors.Wrap(err, internal.ErrMsgBlockCreatingToken))
		return
	}

	c.presenter.PresentResponse(ctx, rw, http.StatusOK, dto.TokenResp{Token: string(token)})
}

func (c *authController) PostAuthTokenAuthorize(rw http.ResponseWriter, req *http.Request) {

	var reqBody dto.AuthorizationRequest
	var err error

	if requestutil.IsDebugDataEnabled(req) {
		reqBodyBytes := requestutil.GetRequestBodyFromCtx(req.Context())
		err = json.NewDecoder(bytes.NewReader(reqBodyBytes)).Decode(&reqBody)
	} else {
		err = json.NewDecoder(req.Body).Decode(&reqBody)
	}
	if err != nil {
		respErr := errors.Wrap(err, internal.ErrMsgBlockUnmarshalingReqBody)
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusBadRequest, respErr, respErr) //respErr == logErr
		return
	}

	if err := reqBody.ValidateFormat(); err != nil {
		respErr := errors.Wrap(err, internal.ErrMsgBlockReqBodyDoesntHaveProperFmt)
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, http.StatusUnprocessableEntity, respErr, respErr) //respErr == logErr
		return
	}

	if err := c.authServices.Authorize(reqBody, "REST"); err != nil {

		localErrMsg := internal.ErrMsgBlockAuthorizingAccess
		var status int
		var respErr, logErr error

		if c.serviceErrorChecker.ErrorIsInvalidToken(err) {
			// Header WWW-Authenticate: Bearer
			// Caller must return this header to the resource requester
			status, respErr, logErr = http.StatusUnauthorized, errors.Wrap(internal.ErrRespInvalidToken, localErrMsg), errors.Wrap(err, localErrMsg)

		} else if c.serviceErrorChecker.ErrorIsNotEnoughPermissions(err) {
			status, respErr, logErr = http.StatusForbidden, errors.Wrap(internal.ErrRespNotEnoughPermissions, localErrMsg), errors.Wrap(err, localErrMsg)

		} else if c.serviceErrorChecker.ErrorIsSemanticallyUnprocesable(err) {
			status, respErr, logErr = http.StatusUnprocessableEntity, errors.Wrap(internal.ErrRespCannotProcessTokenCorrectly, localErrMsg), errors.Wrap(err, localErrMsg)

		} else {
			// Error category should match one of the three above
			// there is no internal error expected
			status, respErr, logErr = http.StatusInternalServerError, errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg), errors.Wrap(err, localErrMsg)
		}
		c.writerLogger.WriteErrorRespAndLog(req.Context(), rw, status, respErr, logErr)
		return
	}

	c.presenter.SuccessResp(req.Context(), rw, http.StatusOK, "Token authorization OK")
}
