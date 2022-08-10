package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxiancillotti/access-control/internal/controllers/internal"
	"github.com/maxiancillotti/access-control/internal/controllers/requestutil"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/mock/writerloggermock"

	"github.com/maxiancillotti/access-control/internal/mock"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizationXAdmin(t *testing.T) {

	errUnauthorized := errors.New("ERROR_UNAUTHORIZED")
	expectedAuthHeaderType := "X-Admin"

	// ADD TEST CASES

	type testCase struct {
		name                 string
		headerValueInput     string
		expectedOutputStatus int
		expectedOutputErr    error
		expectedLoggedErr    error
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "success",
		headerValueInput:     "X-Admin dXNlclN1Y2Nlc3M6cGFzc3dvcmQ=",
		expectedOutputStatus: http.StatusOK,
		expectedOutputErr:    nil,
		expectedLoggedErr:    nil,
	})

	table = append(table, testCase{
		name:                 "Error: authorization header is empty",
		headerValueInput:     "",
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    internal.ErrLogAuthorizationHeaderIsEmpty,
	})

	table = append(table, testCase{
		name:                 "Error: authorization must be Basic",
		headerValueInput:     "Token 7ok3n",
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogFmtAuthorizationMustBeOfType, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: cannot base64 decode credentials",
		headerValueInput:     "X-Admin user:pw",
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogBlockFmtAuthorizationCannotBase64Decode, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: cannot find credentials separator",
		headerValueInput:     "X-Admin dXNlcnBhc3N3b3Jk", // base64Enc(userpassword)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogFmtAuthorizationCannotFindCredentialsSeparator, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: credentials cannot be empty. Case user.",
		headerValueInput:     "X-Admin OnB3", // base64Enc(:pw)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogFmtAuthorizationCredentialsCannotBeEmpty, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: credentials cannot be empty. Case password.",
		headerValueInput:     "X-Admin dXNlcjo=", // base64Enc(user:)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogFmtAuthorizationCredentialsCannotBeEmpty, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: credentials cannot be empty. Case both.",
		headerValueInput:     "X-Admin Og==", // base64Enc(:)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogFmtAuthorizationCredentialsCannotBeEmpty, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: credentials invalid format. Case user.",
		headerValueInput:     "X-Admin aW52YWxpZFVzZXIkJCQkOnB3", // base64Enc(invalidUser$$$$:pw)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogBlockFmtAuthorizationInvalidCredentialsCauseFormat, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: credentials invalid format. Case password.",
		headerValueInput:     "X-Admin dXNlcjpJbnZhbGlkUFfOps6mzqbOpg==", // base64Enc(user:InvalidPWΦΦΦΦ)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogBlockFmtAuthorizationInvalidCredentialsCauseFormat, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: invalid credentials",
		headerValueInput:     "X-Admin dXNlckludmFsaWRDcmVkZW50aWFsczpwYXNzd29yZA==", // base64Enc(userInvalidCredentials:password)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogBlockFmtAuthorizationInvalidCredentials, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: unexpected internal error",
		headerValueInput:     "X-Admin dXNlclVuZXhwZWN0ZWRFcnJvcjpwYXNzd29yZA==", // base64Enc(userUnexpectedError:password)
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputErr:    internal.ErrRespInternalUnexpected,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogBlockFmtAuthorizationService, expectedAuthHeaderType),
	})

	//////////////////////////////////////////////
	// TESTING BEGINS

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
			rwr := httptest.NewRecorder()

			req.Header.Set("Authorization", test.headerValueInput)

			// Execution
			testAuthController.AuthorizationXAdmin(handlerMockforMDW).ServeHTTP(rwr, req)

			// Check
			statusCode := rwr.Result().StatusCode
			assert.Equal(t, test.expectedOutputStatus, statusCode)

			if test.expectedOutputStatus > 299 {

				var errResp mock.ResponseMockError
				err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
				assert.Nil(t, err)

				assert.Contains(t, errResp.Msg, test.expectedOutputErr.Error())

				if test.expectedOutputStatus == http.StatusUnauthorized {
					headerAuthenticate := rwr.Result().Header.Get("WWW-Authenticate")
					assert.Equal(t, expectedAuthHeaderType, headerAuthenticate)

					failingHeaderKeyLogged := rwr.Result().Header.Get(writerloggermock.TestFailingHKHeaderKey)
					assert.Contains(t, failingHeaderKeyLogged, "Authorization")

					failingHeaderValueLogged := rwr.Result().Header.Get(writerloggermock.TestFailingHVHeaderKey)
					assert.Contains(t, failingHeaderValueLogged, test.headerValueInput)
				}

				errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
				assert.Contains(t, errLogged, test.expectedLoggedErr.Error())
			}
		})
	}
}

func TestAuthorizationBasic(t *testing.T) {

	errUnauthorized := errors.New("ERROR_UNAUTHORIZED")
	expectedAuthHeaderType := "Basic"

	// ADD TEST CASES

	type testCase struct {
		name                 string
		headerValueInput     string
		expectedOutputStatus int
		expectedOutputErr    error
		expectedLoggedErr    error
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "success",
		headerValueInput:     "Basic dXNlclN1Y2Nlc3M6cGFzc3dvcmQ=", // base64Enc(userSuccess:password)
		expectedOutputStatus: http.StatusOK,
		expectedOutputErr:    nil,
		expectedLoggedErr:    nil,
	})

	table = append(table, testCase{
		name:                 "Error: authorization header is empty",
		headerValueInput:     "",
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    internal.ErrLogAuthorizationHeaderIsEmpty,
	})

	table = append(table, testCase{
		name:                 "Error: authorization must be Basic",
		headerValueInput:     "Token 7ok3n",
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogFmtAuthorizationMustBeOfType, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: cannot base64 decode credentials",
		headerValueInput:     "Basic user:pw",
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogBlockFmtAuthorizationCannotBase64Decode, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: cannot find credentials separator",
		headerValueInput:     "Basic dXNlcnBhc3N3b3Jk", // base64Enc(userpassword)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogFmtAuthorizationCannotFindCredentialsSeparator, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: credentials cannot be empty. Case user.",
		headerValueInput:     "Basic OnB3", // base64Enc(:pw)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogFmtAuthorizationCredentialsCannotBeEmpty, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: credentials cannot be empty. Case password.",
		headerValueInput:     "Basic dXNlcjo=", // base64Enc(user:)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogFmtAuthorizationCredentialsCannotBeEmpty, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: credentials cannot be empty. Case both.",
		headerValueInput:     "Basic Og==", // base64Enc(:)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogFmtAuthorizationCredentialsCannotBeEmpty, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: credentials invalid format. Case user.",
		headerValueInput:     "Basic aW52YWxpZFVzZXIkJCQkOnB3", // base64Enc(invalidUser$$$$:pw)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogBlockFmtAuthorizationInvalidCredentialsCauseFormat, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: credentials invalid format. Case password.",
		headerValueInput:     "Basic dXNlcjpJbnZhbGlkUFfOps6mzqbOpg==", // base64Enc(user:InvalidPWΦΦΦΦ)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogBlockFmtAuthorizationInvalidCredentialsCauseFormat, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: invalid credentials",
		headerValueInput:     "Basic dXNlckludmFsaWRDcmVkZW50aWFsczpwYXNzd29yZA==", // base64Enc(userInvalidCredentials:password)
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputErr:    errUnauthorized,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogBlockFmtAuthorizationInvalidCredentials, expectedAuthHeaderType),
	})

	table = append(table, testCase{
		name:                 "Error: unexpected internal error",
		headerValueInput:     "Basic dXNlclVuZXhwZWN0ZWRFcnJvcjpwYXNzd29yZA==", // base64Enc(userUnexpectedError:password)
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputErr:    internal.ErrRespInternalUnexpected,
		expectedLoggedErr:    errors.Errorf(internal.ErrLogBlockFmtAuthorizationService, expectedAuthHeaderType),
	})

	//////////////////////////////////////////////
	// TESTING BEGINS

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
			rwr := httptest.NewRecorder()

			req.Header.Set("Authorization", test.headerValueInput)

			// Execution
			testAuthController.AuthorizationBasic(handlerMockforMDW).ServeHTTP(rwr, req)

			// Check
			statusCode := rwr.Result().StatusCode
			assert.Equal(t, test.expectedOutputStatus, statusCode)

			if test.expectedOutputStatus > 299 {

				var errResp mock.ResponseMockError
				err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
				assert.Nil(t, err)

				assert.Contains(t, errResp.Msg, test.expectedOutputErr.Error())

				if test.expectedOutputStatus == http.StatusUnauthorized {
					headerAuthenticate := rwr.Result().Header.Get("WWW-Authenticate")
					assert.Equal(t, expectedAuthHeaderType, headerAuthenticate)

					failingHeaderKeyLogged := rwr.Result().Header.Get(writerloggermock.TestFailingHKHeaderKey)
					assert.Contains(t, failingHeaderKeyLogged, "Authorization")

					failingHeaderValueLogged := rwr.Result().Header.Get(writerloggermock.TestFailingHVHeaderKey)
					assert.Contains(t, failingHeaderValueLogged, test.headerValueInput)
				}

				errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
				assert.Contains(t, errLogged, test.expectedLoggedErr.Error())
			}
		})
	}
}

func TestPostAuthTokenAuthorize(t *testing.T) {

	// ADD TEST CASES

	type testCase struct {
		name                 string
		input                interface{}
		expectedOutputStatus int
		expectedOutputMsg    string
		expectedOutputMsg2   string
		expectedLoggedErrMsg string
		debugModeEnabled     bool
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name: "Success",
		input: dto.AuthorizationRequest{
			Token:             "valid",
			ResourceRequested: "/resource",
			MethodRequested:   "METHOD",
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputMsg:    "Token authorization OK",
	})

	table = append(table, testCase{
		name: "Success debug enabled",
		input: dto.AuthorizationRequest{
			Token:             "valid",
			ResourceRequested: "/resource",
			MethodRequested:   "METHOD",
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputMsg:    "Token authorization OK",

		debugModeEnabled: true,
	})

	table = append(table, testCase{
		name:                 "Error unmarshaling body",
		input:                `{"key":"value"}`, // Error trigger: invalid body
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputMsg:    internal.ErrMsgBlockUnmarshalingReqBody,

		expectedLoggedErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
	})

	table = append(table, testCase{
		name: "Error invalid body format",
		input: dto.AuthorizationRequest{
			Token:             "valid",
			ResourceRequested: "/resource",
			MethodRequested:   "METHOD123", // Error trigger: invalid numeric chars
		},
		expectedOutputStatus: http.StatusUnprocessableEntity,
		expectedOutputMsg:    internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,

		expectedLoggedErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
	})

	table = append(table, testCase{
		name: "Error invalid token",
		input: dto.AuthorizationRequest{
			Token:             "invalid", // Error trigger
			ResourceRequested: "/resource",
			MethodRequested:   "METHOD",
		},
		expectedOutputStatus: http.StatusUnauthorized,
		expectedOutputMsg:    internal.ErrRespInvalidToken.Error(),
		expectedOutputMsg2:   internal.ErrMsgBlockAuthorizingAccess,

		expectedLoggedErrMsg: internal.ErrMsgBlockAuthorizingAccess,
	})

	table = append(table, testCase{
		name: "Error permissions",
		input: dto.AuthorizationRequest{
			Token:             "not_enough_permissions", // Error trigger
			ResourceRequested: "/resource",
			MethodRequested:   "METHOD",
		},
		expectedOutputStatus: http.StatusForbidden,
		expectedOutputMsg:    internal.ErrRespNotEnoughPermissions.Error(),
		expectedOutputMsg2:   internal.ErrMsgBlockAuthorizingAccess,

		expectedLoggedErrMsg: internal.ErrMsgBlockAuthorizingAccess,
	})

	table = append(table, testCase{
		name: "Error unprocessable token",
		input: dto.AuthorizationRequest{
			Token:             "unprocessable", // Error trigger
			ResourceRequested: "/resource",
			MethodRequested:   "METHOD",
		},
		expectedOutputStatus: http.StatusUnprocessableEntity,
		expectedOutputMsg:    internal.ErrRespCannotProcessTokenCorrectly.Error(),
		expectedOutputMsg2:   internal.ErrMsgBlockAuthorizingAccess,

		expectedLoggedErrMsg: internal.ErrMsgBlockAuthorizingAccess,
	})

	table = append(table, testCase{
		name: "Error unexpected on validation",
		input: dto.AuthorizationRequest{
			Token:             "unexpected_internal", // Error trigger
			ResourceRequested: "/resource",
			MethodRequested:   "METHOD",
		},
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputMsg:    internal.ErrRespInternalUnexpected.Error(),
		expectedOutputMsg2:   internal.ErrMsgBlockAuthorizingAccess,

		expectedLoggedErrMsg: internal.ErrMsgBlockAuthorizingAccess,
	})

	/*
		table = append(table, testCase{
			name: "",
			input: ,
			expectedOutputStatus: ,
			expectedOutputMsg:    ,
		})
	*/
	//////////////////////////////////////////////
	// TESTING BEGINS

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var bufBodyBytes []byte
			bufReqBody := bytes.NewBuffer(bufBodyBytes)
			err := json.NewEncoder(bufReqBody).Encode(&test.input)
			if err != nil {
				t.Fatal("cannot initialize request body")
			}

			req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken/authorize", bufReqBody)
			rwr := httptest.NewRecorder()

			// Debug Mode Enabled: body is set into ctx
			if test.debugModeEnabled {
				reqBodyBytes, err := ioutil.ReadAll(bufReqBody)
				if err != nil && err != io.EOF {
					t.Fatal("cannot initialize debug mode")
				}

				ctx := requestutil.AssignRequestBodyToCtx(req.Context(), reqBodyBytes)
				// req = req.Clone(ctx) //deep copy
				req = req.WithContext(ctx) // shallow copy
				req.Header.Set("X-Debug-Enabled", "true")
			}

			// Execution
			testAuthController.PostAuthTokenAuthorize(rwr, req)

			// Check
			statusCode := rwr.Result().StatusCode
			assert.Equal(t, test.expectedOutputStatus, statusCode)

			if test.expectedOutputStatus < 300 {
				// success case
				var successResp mock.ResponseMockSuccess
				err = json.NewDecoder(rwr.Result().Body).Decode(&successResp)
				assert.Nil(t, err)

				assert.Equal(t, test.expectedOutputMsg, successResp.Msg)
			} else {
				// error case
				var errResp mock.ResponseMockError
				err = json.NewDecoder(rwr.Result().Body).Decode(&errResp)
				assert.Nil(t, err)

				assert.Contains(t, errResp.Msg, test.expectedOutputMsg)

				if test.expectedOutputMsg2 != "" {
					assert.Contains(t, errResp.Msg, test.expectedOutputMsg2)
				}

				errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
				assert.Contains(t, errLogged, test.expectedLoggedErrMsg)
			}
		})
	}
}

type testPostAuthTokenCase struct {
	name                 string
	inputUserID          uint
	expectedOutputStatus int
	expectedOutputMsg    string
	expectedLoggedErrMsg string
}

func TestPostAuthToken(t *testing.T) {

	// ADD TEST CASES

	table := make([]testPostAuthTokenCase, 0)

	table = append(table, testPostAuthTokenCase{
		name:                 "Success",
		inputUserID:          1,
		expectedOutputStatus: http.StatusOK,
		expectedOutputMsg:    "header.payload.sign",
		expectedLoggedErrMsg: "",
	})

	table = append(table, testPostAuthTokenCase{
		name: "Error credentials in ctx not found",
		//input:                0, // error trigger: no userID
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputMsg:    internal.ErrRespInternalUnexpected.Error(),
		expectedLoggedErrMsg: internal.ErrLogRetrievingUserIDFromContext.Error(),
	})

	table = append(table, testPostAuthTokenCase{
		name:                 "Error internal creating token",
		inputUserID:          2, // error trigger: userID mocked to fail on call to service
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputMsg:    internal.ErrRespInternalUnexpected.Error(),
		expectedLoggedErrMsg: internal.ErrMsgBlockCreatingToken,
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
			rwr := httptest.NewRecorder()

			if test.inputUserID != 0 {
				ctx := req.Context()
				ctx = requestutil.AssignAuthenticatedUserIDToCxt(ctx, test.inputUserID)
				// req = req.Clone(ctx) //deep copy
				req = req.WithContext(ctx) // shallow copy
			}

			// Execution
			testAuthController.PostAuthToken(rwr, req)

			// Check
			statusCode := rwr.Result().StatusCode
			assert.Equal(t, test.expectedOutputStatus, statusCode)

			if test.expectedOutputStatus > 299 {
				// error case
				var errResp mock.ResponseMockError
				err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
				assert.Nil(t, err)

				assert.Contains(t, errResp.Msg, test.expectedOutputMsg)

				errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
				assert.Contains(t, errLogged, test.expectedLoggedErrMsg)

			} else {
				// success case
				var token dto.TokenResp
				err := json.NewDecoder(rwr.Result().Body).Decode(&token)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputMsg, token.Token)
			}
		})

	}
}
