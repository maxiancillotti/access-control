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

func TestResourcesPOST(t *testing.T) {

	type testCase struct {
		name                 string
		inputResourcePath    *dto.Resource
		expectedOutputStatus int
		expectedOutputResp   *dto.Resource
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
		debugEnabled         bool
	}

	// ADD TEST CASES

	uintOne := uint(1)

	pathStr1 := "path1"
	pathStr0 := "path0"
	pathStr2 := "path2"

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "Success",
		inputResourcePath:    &dto.Resource{Path: &pathStr1},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.Resource{
			ID:   &uintOne,
			Path: &pathStr1,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	table = append(table, testCase{
		name:                 "Success. Debug data enabled",
		inputResourcePath:    &dto.Resource{Path: &pathStr1},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.Resource{
			ID:   &uintOne,
			Path: &pathStr1,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
		debugEnabled:         true,
	})

	table = append(table, testCase{
		name:                 "Error empty body",
		inputResourcePath:    nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
		expectedLoggedErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
	})

	table = append(table, testCase{
		name:                 "Error empty path",
		inputResourcePath:    &dto.Resource{Path: nil},
		expectedOutputStatus: http.StatusUnprocessableEntity,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
		expectedLoggedErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtCreating, "resource")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error creating resource invalid path",
		inputResourcePath:    &dto.Resource{Path: &pathStr0},
		expectedOutputStatus: http.StatusConflict,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error creating resource internal",
		inputResourcePath:    &dto.Resource{Path: &pathStr2},
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/resources"

			buf := new(bytes.Buffer)

			if test.inputResourcePath != nil {
				err := json.NewEncoder(buf).Encode(test.inputResourcePath)
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
			testResourcesController.POST(rwr, req)

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
				var resp dto.Resource
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp.ID, resp.ID)
				assert.Equal(t, test.expectedOutputResp.Path, resp.Path)
			}
		})

	}
}

func TestResourcesDELETE(t *testing.T) {

	type testCase struct {
		name                 string
		inputResourceID      *uint
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
		inputResourceID:      &uintOne,
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp:   "resource deleted OK",
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	errParsingStr := errors.Errorf(internal.ErrMsgFmtParsingIDFromURL, "resource").Error()
	table = append(table, testCase{
		name:                 "Error parsing resourceID from URL",
		inputResourceID:      nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   "",
		expectedOutputErrMsg: errParsingStr,
		expectedLoggedErrMsg: errParsingStr,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtDeleting, "resource")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error deleting resource not found",
		inputResourceID:      &uintZero,
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   "",
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error deleting resource internal",
		inputResourceID:      &uintTwo,
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   "",
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/resources/{id}"
			if test.inputResourceID != nil {
				url = fmt.Sprintf("http://localhost:8001/api/resources/%d", *test.inputResourceID)
			}
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			rwr := httptest.NewRecorder()

			//Faking gorilla/mux vars
			var varsValue string

			if test.inputResourceID != nil {
				varsValue = strconv.FormatUint(uint64(*test.inputResourceID), 10)
			}

			vars := map[string]string{
				"id": varsValue,
			}
			req = mux.SetURLVars(req, vars)

			// Execution
			testResourcesController.DELETE(rwr, req)

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

func TestResourcesGET(t *testing.T) {

	type testCase struct {
		name                 string
		inputResourceID      *uint
		expectedOutputStatus int
		expectedOutputResp   *dto.Resource
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
	}

	// ADD TEST CASES

	uintZero := uint(0)
	uintOne := uint(1)
	uintTwo := uint(2)

	pathStr := "path"

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "Success",
		inputResourceID:      &uintOne,
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.Resource{
			ID:   &uintOne,
			Path: &pathStr,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	errParsingStr := errors.Errorf(internal.ErrMsgFmtParsingIDFromURL, "resource").Error()
	table = append(table, testCase{
		name:                 "Error parsing resourceID from URL",
		inputResourceID:      nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errParsingStr,
		expectedLoggedErrMsg: errParsingStr,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "resource")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting resource not found",
		inputResourceID:      &uintZero,
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting resource internal",
		inputResourceID:      &uintTwo,
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/resources/{id}"
			if test.inputResourceID != nil {
				url = fmt.Sprintf("http://localhost:8001/api/resources/%d", *test.inputResourceID)
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			rwr := httptest.NewRecorder()

			//Faking gorilla/mux vars
			var varsValue string

			if test.inputResourceID != nil {
				varsValue = strconv.FormatUint(uint64(*test.inputResourceID), 10)
			}

			vars := map[string]string{
				"id": varsValue,
			}
			req = mux.SetURLVars(req, vars)

			// Execution
			testResourcesController.GET(rwr, req)

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
				var resp dto.Resource
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp.ID, resp.ID)
				assert.Equal(t, test.expectedOutputResp.Path, resp.Path)
			}
		})

	}
}

func TestResourcesGETByPath(t *testing.T) {

	type testCase struct {
		name                 string
		inputResourcePath    *dto.Resource
		expectedOutputStatus int
		expectedOutputResp   *dto.Resource
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
		debugEnabled         bool
	}

	// ADD TEST CASES

	uintOne := uint(1)

	pathStr1 := "path1"
	pathStr0 := "path0"
	pathStr2 := "path2"

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "Success",
		inputResourcePath:    &dto.Resource{Path: &pathStr1},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.Resource{
			ID:   &uintOne,
			Path: &pathStr1,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	table = append(table, testCase{
		name:                 "Success. Debug data enabled",
		inputResourcePath:    &dto.Resource{Path: &pathStr1},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.Resource{
			ID:   &uintOne,
			Path: &pathStr1,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
		debugEnabled:         true,
	})

	table = append(table, testCase{
		name:                 "Error empty body",
		inputResourcePath:    nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
		expectedLoggedErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
	})

	table = append(table, testCase{
		name:                 "Error empty path",
		inputResourcePath:    &dto.Resource{Path: nil},
		expectedOutputStatus: http.StatusUnprocessableEntity,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
		expectedLoggedErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "resource")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting resource not found",
		inputResourcePath:    &dto.Resource{Path: &pathStr0},
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting resource internal",
		inputResourcePath:    &dto.Resource{Path: &pathStr2},
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/resources"

			buf := new(bytes.Buffer)

			if test.inputResourcePath != nil {
				err := json.NewEncoder(buf).Encode(test.inputResourcePath)
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
			testResourcesController.GETByPath(rwr, req)

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
				var resp dto.Resource
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp.ID, resp.ID)
				assert.Equal(t, test.expectedOutputResp.Path, resp.Path)
			}
		})

	}
}

func TestResourcesGETCollection(t *testing.T) {

	type testCase struct {
		name                 string
		svcErrReturnTrigger  int
		expectedOutputStatus int
		expectedOutputResp   []dto.Resource
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
	}

	// ADD TEST CASES

	uintOne := uint(1)
	pathStr := "path"

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "Success",
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: []dto.Resource{
			{
				ID:   &uintOne,
				Path: &pathStr,
			},
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "resources")
	privateErr := errors.Wrap(errors.New("empty_result"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting resources empty result",
		svcErrReturnTrigger:  1,
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting resource internal",
		svcErrReturnTrigger:  2,
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/resources"

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rwr := httptest.NewRecorder()

			// Execution
			ResourcesSVCMock_RetrieveAll_ReturnSwitch = test.svcErrReturnTrigger

			testResourcesController.GETCollection(rwr, req)

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
				var resp []dto.Resource
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp[0].ID, resp[0].ID)
				assert.Equal(t, test.expectedOutputResp[0].Path, resp[0].Path)
			}
		})

	}
}
