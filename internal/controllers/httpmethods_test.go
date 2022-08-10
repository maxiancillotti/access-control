package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxiancillotti/access-control/internal/controllers/internal"
	"github.com/maxiancillotti/access-control/internal/controllers/requestutil"
	"github.com/maxiancillotti/access-control/internal/dto"
	"github.com/maxiancillotti/access-control/internal/mock"
	"github.com/maxiancillotti/access-control/internal/mock/writerloggermock"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestHttpMethodsGETByName(t *testing.T) {

	type testCase struct {
		name                 string
		inputMethodName      *dto.HttpMethod
		expectedOutputStatus int
		expectedOutputResp   *dto.HttpMethod
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
		debugEnabled         bool
	}

	// ADD TEST CASES

	uintOne := uint(1)

	methodNameStr1 := "methodName1"
	methodNameStr0 := "methodName0"
	methodNameStr2 := "methodName2"

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "Success",
		inputMethodName:      &dto.HttpMethod{Name: &methodNameStr1},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.HttpMethod{
			ID:   &uintOne,
			Name: &methodNameStr1,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	table = append(table, testCase{
		name:                 "Success. Debug data enabled",
		inputMethodName:      &dto.HttpMethod{Name: &methodNameStr1},
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: &dto.HttpMethod{
			ID:   &uintOne,
			Name: &methodNameStr1,
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
		debugEnabled:         true,
	})

	table = append(table, testCase{
		name:                 "Error empty body",
		inputMethodName:      nil,
		expectedOutputStatus: http.StatusBadRequest,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
		expectedLoggedErrMsg: internal.ErrMsgBlockUnmarshalingReqBody,
	})

	table = append(table, testCase{
		name:                 "Error empty name",
		inputMethodName:      &dto.HttpMethod{Name: nil},
		expectedOutputStatus: http.StatusUnprocessableEntity,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
		expectedLoggedErrMsg: internal.ErrMsgBlockReqBodyDoesntHaveProperFmt,
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "method")
	privateErr := errors.Wrap(errors.New("invalid_input"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting method not found",
		inputMethodName:      &dto.HttpMethod{Name: &methodNameStr0},
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting method internal",
		inputMethodName:      &dto.HttpMethod{Name: &methodNameStr2},
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/httpmethods"

			buf := new(bytes.Buffer)

			if test.inputMethodName != nil {
				err := json.NewEncoder(buf).Encode(test.inputMethodName)
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
			testHttpMethodsController.GETByName(rwr, req)

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
				var resp dto.HttpMethod
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp.ID, resp.ID)
				assert.Equal(t, test.expectedOutputResp.Name, resp.Name)
			}
		})

	}
}

func TestHttpMethodsGETCollection(t *testing.T) {

	type testCase struct {
		name                 string
		svcErrReturnTrigger  int
		expectedOutputStatus int
		expectedOutputResp   []dto.HttpMethod
		expectedOutputErrMsg string
		expectedLoggedErrMsg string
	}

	// ADD TEST CASES

	uintOne := uint(1)
	methodNameStr := "methodName"

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:                 "Success",
		expectedOutputStatus: http.StatusOK,
		expectedOutputResp: []dto.HttpMethod{
			{
				ID:   &uintOne,
				Name: &methodNameStr,
			},
		},
		expectedOutputErrMsg: "",
		expectedLoggedErrMsg: "",
	})

	localErrMsg := fmt.Sprintf(internal.ErrMsgFmtGetting, "methods")
	privateErr := errors.Wrap(errors.New("empty_result"), localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting methods empty result",
		svcErrReturnTrigger:  1,
		expectedOutputStatus: http.StatusNotFound,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: privateErr.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	privateErr = errors.Wrap(errors.New("internal_error"), localErrMsg)
	errResp := errors.Wrap(internal.ErrRespInternalUnexpected, localErrMsg)

	table = append(table, testCase{
		name:                 "Error getting methods internal",
		svcErrReturnTrigger:  2,
		expectedOutputStatus: http.StatusInternalServerError,
		expectedOutputResp:   nil,
		expectedOutputErrMsg: errResp.Error(),
		expectedLoggedErrMsg: privateErr.Error(),
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			var url string = "http://localhost:8001/api/httpmethods"

			req := httptest.NewRequest(http.MethodGet, url, nil)
			rwr := httptest.NewRecorder()

			// Execution
			ResourcesSVCMock_RetrieveAll_ReturnSwitch = test.svcErrReturnTrigger

			testHttpMethodsController.GETCollection(rwr, req)

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
				var resp []dto.HttpMethod
				err := json.NewDecoder(rwr.Result().Body).Decode(&resp)
				assert.Nil(t, err)
				assert.Equal(t, test.expectedOutputResp[0].ID, resp[0].ID)
				assert.Equal(t, test.expectedOutputResp[0].Name, resp[0].Name)
			}
		})

	}
}
