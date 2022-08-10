package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/maxiancillotti/access-control/internal/controllers/internal"
	"github.com/maxiancillotti/access-control/internal/controllers/requestutil"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/mock"
	"github.com/maxiancillotti/access-control/internal/mock/writerloggermock"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUsersPOST(t *testing.T) {

	type testCase struct {
		name                 string
		inputUser            *dto.User
		expectedOutputStatus int
		expectedOutputResp   *dto.User
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
		debugEnabled         bool
	}

	// ADD TEST CASES

	uintOne := uint(1)

	usernameStr1 := "username1"
	usernameStr0 := "username0"
	usernameStr2 := "username2"
	usernameStrInvalid := "usernameInvalid$$$"

	passwordStr := "password"

	table := make([]testCase, 0)

	table = append(table, testCase{
		name: "Success",
		inputUser: &dto.User{
			Username: &usernameStr1,
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.User{
			ID:       &uintOne,
			Username: &usernameStr1,
			Password: &passwordStr,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	table = append(table, testCase{
		name: "Success. Debug data enabled.",
		inputUser: &dto.User{
			Username: &usernameStr1,
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.User{
			ID:       &uintOne,
			Username: &usernameStr1,
			Password: &passwordStr,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
		debugEnabled:         true,
	})

	table = append(table, testCase{
		name:                 "Error empty body",
		inputUser:            nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
		expectedLoggedErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
	})

	table = append(table, testCase{
		name: "Error invalid username",
		inputUser: &dto.User{
			Username: &usernameStrInvalid,
		},
		expectedOutputStatus: http.StatusUnprocessableEntity,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
		expectedLoggedErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtCreating, "user")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name: "Error creating user conflict",
		inputUser: &dto.User{
			Username: &usernameStr0,
		},
		expectedOutputStatus: http.StatusConflict,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name: "Error creating user internal",
		inputUser: &dto.User{
			Username: &usernameStr2,
		},
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/users"

			buf := new(bytes.Buffer)

			if test.inputUser != nil {
				err := json.NewEncoder(buf).Encode(test.inputUser)
				assert.Nil(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, url, buf)
			rwr := httptest.NewRecorder()

			if test.debugEnabled {
				req.Header.Set("X-Debug-Enabled", "true")
				ctx := requestutil.AssignRequestBodyToCtx(req.Context(), buf.Bytes())
				req = req.WithContext(ctx)
			}

			// Execution
			testUserController.POST(rwr, req)

			// Check
			statusCode := rwr.Result().StatusCode
			assert.Equal(t, test.expectedOutputStatus, statusCode)

			if test.expectedOutputStatus > 299 {
				// error case
				var errResp mock.ResponseMockError
				err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
				assert.Nil(t, err)

				assert.Contains(t, errResp.Msg, test.expectedOutputErrMsg)

				errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
				assert.Contains(t, errLogged, test.expectedLoggedErrMsg)

			} else {
				// success case
				var resp dto.User
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp.ID, resp.ID)
				assert.Equal(t, test.expectedOutputResp.Username, resp.Username)
				assert.Equal(t, test.expectedOutputResp.Password, resp.Password)
			}
		})

	}
}

func TestUsersPATCHPassword(t *testing.T) {

	type testCase struct {
		name                 string
		inputUserID          *uint
		expectedOutputStatus int
		expectedOutputResp   *dto.User
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
	}

	// ADD TEST CASES

	uintZero := uint(0)
	uintOne := uint(1)
	uintTwo := uint(2)
	passwordStr := "password"

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "Success",
		inputUserID:          &uintOne,
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp:   &dto.User{ID: &uintOne, Password: &passwordStr},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	errParsingStr := errors.Errorf(internal.ErrMsgFmtParsingIDFromURL, "user").Error()
	table = append(table, testCase{
		name:                 "Error parsing userID from URL",
		inputUserID:          nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errParsingStr,
		expectedLoggedErrMsg: errParsingStr,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtUpdating, "user")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error updating user password not found",
		inputUserID:          &uintZero,
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error updating user password internal",
		inputUserID:          &uintTwo,
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/users/{id}/Password"
			if test.inputUserID != nil {
				url = fmt.Sprintf("http://localhost:8001/api/users/%d/Password", *test.inputUserID)
			}

			req := httptest.NewRequest(http.MethodPatch, url, nil)
			rwr := httptest.NewRecorder()

			//Faking gorilla/mux vars
			var varsValue string

			if test.inputUserID != nil {
				varsValue = strconv.FormatUint(uint64(*test.inputUserID), 10)
			}

			vars := map[string]string{
				"id": varsValue,
			}
			req = mux.SetURLVars(req, vars)

			// Execution
			testUserController.PATCHPassword(rwr, req)

			// Check
			statusCode := rwr.Result().StatusCode
			assert.Equal(t, test.expectedOutputStatus, statusCode)

			if test.expectedOutputStatus > 299 {
				// error case
				var errResp mock.ResponseMockError
				err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
				assert.Nil(t, err)

				assert.Contains(t, errResp.Msg, test.expectedOutputErrMsg)

				errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
				assert.Contains(t, errLogged, test.expectedLoggedErrMsg)

			} else {
				// success case
				var resp dto.User
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp.ID, resp.ID)
				assert.Equal(t, test.expectedOutputResp.Password, resp.Password)
			}
		})

	}
}

func TestUsersPATCHEnabledState(t *testing.T) {

	type testCase struct {
		name                 string
		inputUserID          *uint
		inputEnableState     *dto.User
		expectedOutputStatus int
		expectedOutputResp   string
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
		debugEnabled         bool
	}

	// ADD TEST CASES

	uintZero := uint(0)
	uintOne := uint(1)
	uintTwo := uint(2)
	boolTrue := true

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:        "Success",
		inputUserID: &uintOne,
		inputEnableState: &dto.User{
			EnabledState: &boolTrue,
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp:   "user enabled state updated OK",
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	table = append(table, testCase{
		name:        "Success. Debug data enabled.",
		inputUserID: &uintOne,
		inputEnableState: &dto.User{
			EnabledState: &boolTrue,
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp:   "user enabled state updated OK",
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
		debugEnabled:         true,
	})

	table = append(table, testCase{
		name:                 "Error empty body",
		inputUserID:          nil,
		inputEnableState:     nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   "",
		expectedOutputErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
		expectedLoggedErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
	})

	errEmptyEnabledStateStr := errors.Wrap(errors.New("state field cannot be empty"), internal.ErrMsgBlockReqBodyDoesntHaveProperFmt).Error()
	table = append(table, testCase{
		name:        "Error empty enabled state",
		inputUserID: nil,
		inputEnableState: &dto.User{
			EnabledState: nil,
		},
		expectedOutputStatus: http.StatusUnprocessableEntity,
		expectedOutputResp:   "",
		expectedOutputErrMsg: errEmptyEnabledStateStr,
		expectedLoggedErrMsg: errEmptyEnabledStateStr,
	})

	errParsingStr := errors.Errorf(internal.ErrMsgFmtParsingIDFromURL, "user").Error()
	table = append(table, testCase{
		name:        "Error parsing userID from URL",
		inputUserID: nil,
		inputEnableState: &dto.User{
			EnabledState: &boolTrue,
		},
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   "",
		expectedOutputErrMsg: errParsingStr,
		expectedLoggedErrMsg: errParsingStr,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtUpdating, "user")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name:        "Error updating user enabled state not found",
		inputUserID: &uintZero,
		inputEnableState: &dto.User{
			EnabledState: &boolTrue,
		},
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   "",
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:        "Error updating user enabled state internal",
		inputUserID: &uintTwo,
		inputEnableState: &dto.User{
			EnabledState: &boolTrue,
		},
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   "",
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/users/{id}/EnabledState"
			if test.inputUserID != nil {
				url = fmt.Sprintf("http://localhost:8001/api/users/%d/EnabledState", *test.inputUserID)
			}

			buf := new(bytes.Buffer)

			if test.inputEnableState != nil {
				err := json.NewEncoder(buf).Encode(test.inputEnableState)
				assert.Nil(t, err)
			}

			req := httptest.NewRequest(http.MethodPatch, url, buf)
			rwr := httptest.NewRecorder()

			if test.debugEnabled {
				req.Header.Set("X-Debug-Enabled", "true")
				ctx := requestutil.AssignRequestBodyToCtx(req.Context(), buf.Bytes())
				req = req.WithContext(ctx)
			}

			//Faking gorilla/mux vars
			var varsValue string

			if test.inputUserID != nil {
				varsValue = strconv.FormatUint(uint64(*test.inputUserID), 10)
			}

			vars := map[string]string{
				"id": varsValue,
			}
			req = mux.SetURLVars(req, vars)

			// Execution
			testUserController.PATCHEnabledState(rwr, req)

			// Check
			statusCode := rwr.Result().StatusCode
			assert.Equal(t, test.expectedOutputStatus, statusCode)

			if test.expectedOutputStatus > 299 {
				// error case
				var errResp mock.ResponseMockError
				err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
				assert.Nil(t, err)

				assert.Contains(t, errResp.Msg, test.expectedOutputErrMsg)

				errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
				assert.Contains(t, errLogged, test.expectedLoggedErrMsg)

			} else {
				// success case
				var resp mock.ResponseMockSuccess
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp, resp.Msg)
			}
		})

	}
}

func TestUsersDELETE(t *testing.T) {

	type testCase struct {
		name                 string
		inputUserID          *uint
		expectedOutputStatus int
		expectedOutputResp   string
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
	}

	// ADD TEST CASES

	uintZero := uint(0)
	uintOne := uint(1)
	uintTwo := uint(2)

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "Success",
		inputUserID:          &uintOne,
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp:   "user deleted OK",
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	errParsingStr := errors.Errorf(internal.ErrMsgFmtParsingIDFromURL, "user").Error()
	table = append(table, testCase{
		name:                 "Error parsing userID from URL",
		inputUserID:          nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   "",
		expectedOutputErrMsg: errParsingStr,
		expectedLoggedErrMsg: errParsingStr,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtDeleting, "user")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error deleting user not found",
		inputUserID:          &uintZero,
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   "",
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error deleting user internal",
		inputUserID:          &uintTwo,
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   "",
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/users/{id}"
			if test.inputUserID != nil {
				url = fmt.Sprintf("http://localhost:8001/api/users/%d", *test.inputUserID)
			}
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			rwr := httptest.NewRecorder()

			//Faking gorilla/mux vars
			var varsValue string

			if test.inputUserID != nil {
				varsValue = strconv.FormatUint(uint64(*test.inputUserID), 10)
			}

			vars := map[string]string{
				"id": varsValue,
			}
			req = mux.SetURLVars(req, vars)

			// Execution
			testUserController.DELETE(rwr, req)

			// Check
			statusCode := rwr.Result().StatusCode
			assert.Equal(t, test.expectedOutputStatus, statusCode)

			if test.expectedOutputStatus > 299 {
				// error case
				var errResp mock.ResponseMockError
				err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
				assert.Nil(t, err)

				assert.Contains(t, errResp.Msg, test.expectedOutputErrMsg)

				errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
				assert.Contains(t, errLogged, test.expectedLoggedErrMsg)

			} else {
				// success case
				var resp mock.ResponseMockSuccess
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp, resp.Msg)
			}
		})

	}
}

func TestUsersGETByUsername(t *testing.T) {

	type testCase struct {
		name                 string
		inputUsername        *dto.User
		expectedOutputStatus int
		expectedOutputResp   *dto.User
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
		debugEnabled         bool
	}

	// ADD TEST CASES

	uintOne := uint(1)

	usernameStr1 := "username1"
	usernameStr0 := "username0"
	usernameStr2 := "username2"

	boolTrue := true

	table := make([]testCase, 0)

	table = append(table, testCase{
		name: "Success",
		inputUsername: &dto.User{
			Username: &usernameStr1,
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.User{
			ID:           &uintOne,
			Username:     &usernameStr1,
			EnabledState: &boolTrue,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	table = append(table, testCase{
		name: "Success. Debug data enabled.",
		inputUsername: &dto.User{
			Username: &usernameStr1,
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.User{
			ID:           &uintOne,
			Username:     &usernameStr1,
			EnabledState: &boolTrue,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
		debugEnabled:         true,
	})

	table = append(table, testCase{
		name:                 "Error empty body",
		inputUsername:        nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
		expectedLoggedErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
	})

	table = append(table, testCase{
		name: "Error empty username",
		inputUsername: &dto.User{
			Username: nil,
		},
		expectedOutputStatus: http.StatusUnprocessableEntity,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
		expectedLoggedErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "user")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name: "Error getting user by username not found",
		inputUsername: &dto.User{
			Username: &usernameStr0,
		},
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name: "Error getting user by username internal",
		inputUsername: &dto.User{
			Username: &usernameStr2,
		},
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/users/"

			buf := new(bytes.Buffer)

			if test.inputUsername != nil {
				err := json.NewEncoder(buf).Encode(test.inputUsername)
				assert.Nil(t, err)
			}

			req := httptest.NewRequest(http.MethodGet, url, buf)
			rwr := httptest.NewRecorder()

			if test.debugEnabled {
				req.Header.Set("X-Debug-Enabled", "true")
				ctx := requestutil.AssignRequestBodyToCtx(req.Context(), buf.Bytes())
				req = req.WithContext(ctx)
			}

			// Execution
			testUserController.GETByUsername(rwr, req)

			// Check
			statusCode := rwr.Result().StatusCode
			assert.Equal(t, test.expectedOutputStatus, statusCode)

			if test.expectedOutputStatus > 299 {
				// error case
				var errResp mock.ResponseMockError
				err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
				assert.Nil(t, err)

				assert.Contains(t, errResp.Msg, test.expectedOutputErrMsg)

				errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
				assert.Contains(t, errLogged, test.expectedLoggedErrMsg)

			} else {
				// success case
				var resp dto.User
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp.ID, resp.ID)
				assert.Equal(t, test.expectedOutputResp.Username, resp.Username)
				assert.Equal(t, test.expectedOutputResp.EnabledState, resp.EnabledState)
			}
		})

	}
}
