package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/maxiancillotti/access-control/internal/controllers/internal"
	"github.com/maxiancillotti/access-control/internal/controllers/requestutil"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/mock"
	"github.com/maxiancillotti/access-control/internal/mock/writerloggermock"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestUsersRESTPermPOST(t *testing.T) {

	type testCase struct {
		name                 string
		input                *dto.UserRESTPermission
		expectedOutputStatus int
		expectedOutputResp   *dto.UserRESTPermission
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
		debugEnabled         bool
	}

	// ADD TEST CASES

	uintZero := uint(0)
	uintOne := uint(1)
	uintTwo := uint(2)

	table := make([]testCase, 0)

	table = append(table, testCase{
		name: "Success",
		input: &dto.UserRESTPermission{
			UserID: &uintOne,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.UserRESTPermission{
			UserID: &uintOne,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	table = append(table, testCase{
		name: "Success. Debug data enabled",
		input: &dto.UserRESTPermission{
			UserID: &uintOne,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.UserRESTPermission{
			UserID: &uintOne,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
		debugEnabled:         true,
	})

	table = append(table, testCase{
		name:                 "Error unmarshaling request body",
		input:                nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
		expectedLoggedErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
	})

	table = append(table, testCase{
		name: "Error empty userID",
		input: &dto.UserRESTPermission{
			//UserID: &uintOne,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
		expectedLoggedErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtCreating, "user REST permissions")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name: "Error creating permission invalid input ID. Combination of IDs already exists.",
		input: &dto.UserRESTPermission{
			UserID: &uintZero,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputStatus: http.StatusConflict,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name: "Error creating permission internal",
		input: &dto.UserRESTPermission{
			UserID: &uintTwo,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization

			var buf bytes.Buffer

			if test.input != nil {
				err := json.NewEncoder(&buf).Encode(*test.input)
				assert.Nil(t, err)
			}

			//t.Log("Input:", *test.input)
			//t.Log("Input buf:", buf.String())

			req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/api/UserRestPermissions", &buf)
			rwr := httptest.NewRecorder()

			if test.debugEnabled {
				req.Header.Set("X-Debug-Enabled", "true")
				ctx := requestutil.AssignRequestBodyToCtx(req.Context(), buf.Bytes())
				req = req.WithContext(ctx)
			}

			// Execution
			testUserRESTPermController.POST(rwr, req)

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
				var urpResp dto.UserRESTPermission
				err := json.NewDecoder(rwr.Result().Body).Decode(&urpResp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp.UserID, urpResp.UserID)
				assert.Equal(t, test.expectedOutputResp.Permission.ResourceID, urpResp.Permission.ResourceID)
				assert.Equal(t, test.expectedOutputResp.Permission.MethodID, urpResp.Permission.MethodID)
			}
		})

	}
}

func TestUsersRESTPermDELETE(t *testing.T) {

	type testCase struct {
		name                 string
		input                *dto.UserRESTPermission
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

	table := make([]testCase, 0)

	table = append(table, testCase{
		name: "Success",
		input: &dto.UserRESTPermission{
			UserID: &uintOne,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp:   "permissions deleted OK",
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	table = append(table, testCase{
		name: "Success. Debug data enabled",
		input: &dto.UserRESTPermission{
			UserID: &uintOne,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp:   "permissions deleted OK",
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
		debugEnabled:         true,
	})

	table = append(table, testCase{
		name:                 "Error unmarshaling request body",
		input:                nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   "",
		expectedOutputErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
		expectedLoggedErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
	})

	table = append(table, testCase{
		name: "Error empty userID",
		input: &dto.UserRESTPermission{
			//UserID: &uintOne,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   "",
		expectedOutputErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
		expectedLoggedErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtDeleting, "user REST permissions")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name: "Error deleting permission invalid input ID. UserID zero value.",
		input: &dto.UserRESTPermission{
			UserID: &uintZero,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   "",
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name: "Error deleting permission internal",
		input: &dto.UserRESTPermission{
			UserID: &uintTwo,
			Permission: &dto.Permission{
				ResourceID: &uintOne,
				MethodID:   &uintOne,
			},
		},
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   "",
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization

			var buf bytes.Buffer

			if test.input != nil {
				err := json.NewEncoder(&buf).Encode(*test.input)
				assert.Nil(t, err)
			}

			//t.Log("Input:", *test.input)
			//t.Log("Input buf:", buf.String())

			req := httptest.NewRequest(http.MethodDelete, "http://localhost:8001/api/UserRestPermissions", &buf)
			rwr := httptest.NewRecorder()

			if test.debugEnabled {
				req.Header.Set("X-Debug-Enabled", "true")
				ctx := requestutil.AssignRequestBodyToCtx(req.Context(), buf.Bytes())
				req = req.WithContext(ctx)
			}

			// Execution
			testUserRESTPermController.DELETE(rwr, req)

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

func TestUsersRESTPermGETCollectionByUserID(t *testing.T) {

	type testCase struct {
		name                 string
		inputUserID          *uint
		expectedOutputStatus int
		expectedOutputResp   *dto.UserRESTPermissionsCollection
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
		expectedOutputResp: &dto.UserRESTPermissionsCollection{
			UserID: 1,
			Permissions: []dto.PermissionsIDs{
				{
					ResourceID: 10,
					MethodsIDs: []uint{1, 2, 3},
				},
			},
		},
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

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "user REST permissions")
	privateErr := errors.Wrap(errors.New("empty_result"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting URP collection not found",
		inputUserID:          &uintZero,
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting URP collection internal",
		inputUserID:          &uintTwo,
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/Users/{userID}/RESTPermissions"
			if test.inputUserID != nil {
				url = fmt.Sprintf("http://localhost:8001/api/Users/%d/RESTPermissions", *test.inputUserID)
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rwr := httptest.NewRecorder()

			//Faking gorilla/mux vars
			var varsValue string

			if test.inputUserID != nil {
				varsValue = strconv.FormatUint(uint64(*test.inputUserID), 10)
			}

			vars := map[string]string{
				"userID": varsValue,
			}
			req = mux.SetURLVars(req, vars)

			// Execution
			testUserRESTPermController.GETCollectionByUserID(rwr, req)

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
				var resp *dto.UserRESTPermissionsCollection
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp.UserID, resp.UserID)
				assert.Equal(t, test.expectedOutputResp.Permissions[0].ResourceID, resp.Permissions[0].ResourceID)
				assert.Equal(t, test.expectedOutputResp.Permissions[0].MethodsIDs[0], resp.Permissions[0].MethodsIDs[0])
			}
		})

	}
}

func TestUsersRESTPermGETCollectionWithDescriptionsByUserID(t *testing.T) {

	type testCase struct {
		name                 string
		inputUserID          *uint
		expectedOutputStatus int
		expectedOutputResp   *dto.UserRESTPermissionsDescriptionsCollection
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
	}

	// ADD TEST CASES

	uintZero := uint(0)
	uintOne := uint(1)
	uintTwo := uint(2)

	uintTen := uint(10)
	resourcePathSuccess := "path"
	methodNameSuccess := "method"

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "Success",
		inputUserID:          &uintOne,
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.UserRESTPermissionsDescriptionsCollection{
			UserID: 1,
			Permissions: []dto.PermissionsWithDescriptions{
				{
					Resource: dto.Resource{
						ID:   &uintTen,
						Path: &resourcePathSuccess,
					},
					Methods: []dto.HttpMethod{
						{
							ID:   &uintTwo,
							Name: &methodNameSuccess,
						},
					},
				},
			},
		},
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

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "user REST permissions")
	privateErr := errors.Wrap(errors.New("empty_result"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting URP collection with descriptions not found",
		inputUserID:          &uintZero,
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting URP collection with descriptions internal",
		inputUserID:          &uintTwo,
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/Users/{userID}/RESTPermissionsWithDescriptions"
			if test.inputUserID != nil {
				url = fmt.Sprintf("http://localhost:8001/api/Users/%d/RESTPermissionsWithDescriptions", *test.inputUserID)
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rwr := httptest.NewRecorder()

			//Faking gorilla/mux vars
			var varsValue string

			if test.inputUserID != nil {
				varsValue = strconv.FormatUint(uint64(*test.inputUserID), 10)
			}

			vars := map[string]string{
				"userID": varsValue,
			}
			req = mux.SetURLVars(req, vars)

			// Execution
			testUserRESTPermController.GETCollectionWithDescriptionsByUserID(rwr, req)

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
				var resp *dto.UserRESTPermissionsDescriptionsCollection
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp.UserID, resp.UserID)
				assert.Equal(t, test.expectedOutputResp.Permissions[0].Resource.ID, resp.Permissions[0].Resource.ID)
				assert.Equal(t, test.expectedOutputResp.Permissions[0].Resource.Path, resp.Permissions[0].Resource.Path)
				assert.Equal(t, test.expectedOutputResp.Permissions[0].Methods[0].ID, resp.Permissions[0].Methods[0].ID)
				assert.Equal(t, test.expectedOutputResp.Permissions[0].Methods[0].Name, resp.Permissions[0].Methods[0].Name)
			}
		})

	}
}
