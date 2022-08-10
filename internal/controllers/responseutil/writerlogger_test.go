package responseutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxiancillotti/access-control/internal/controllers/internal"
	"github.com/maxiancillotti/access-control/internal/mock"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	pstrMock mock.PresenterMock
	lgrMock  = mock.LoggerMock

	testWriterLogger WriterLogger = NewWriterLogger(&pstrMock, lgrMock)
)

// TODO: test logging.

func TestWriteErrorRespAndLog(t *testing.T) {

	// Initialization
	type testCase struct {
		name           string
		expectedStatus int
		expectedErr    error
		errToLog       error
	}

	testTable := make([]testCase, 0)

	testTable = append(testTable, testCase{
		name:           "Status 500",
		expectedStatus: http.StatusInternalServerError,
		expectedErr:    errors.New("error msg response"),
		errToLog:       errors.New("error msg to log"),
	})

	testTable = append(testTable, testCase{
		name:           "Status 400",
		expectedStatus: http.StatusBadRequest,
		expectedErr:    errors.New("error msg response"),
		errToLog:       errors.New("error msg to log"),
	})

	for _, test := range testTable {

		t.Run(test.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
			rwr := httptest.NewRecorder()
			ctx := req.Context()

			// Execution
			testWriterLogger.WriteErrorRespAndLog(ctx, rwr, test.expectedStatus, test.expectedErr, test.errToLog)

			// Check
			resultStatus := rwr.Result().StatusCode
			assert.Equal(t, test.expectedStatus, resultStatus)

			var errResp mock.ResponseMockError
			err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
			assert.Nil(t, err)

			assert.Equal(t, test.expectedErr.Error(), errResp.Msg)
		})
	}
}

func TestWritePresentationHeaderErrorRespAndLog(t *testing.T) {

	// Initialization
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
	rwr := httptest.NewRecorder()

	ctx := req.Context()

	expectedStatus := http.StatusNotAcceptable
	errToWriteAndLog := errors.New("error msg")
	failingHeaderKey := "failingHeaderkey"
	failingHeaderValue := "failing header value"

	// Execution
	testWriterLogger.WritePresentationHeaderErrorRespAndLog(ctx, rwr, errToWriteAndLog, failingHeaderKey, failingHeaderValue)

	// Check
	resultStatus := rwr.Result().StatusCode
	assert.Equal(t, expectedStatus, resultStatus)

	var errResp mock.ResponseMockError
	err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
	assert.Nil(t, err)

	assert.Equal(t, errToWriteAndLog.Error(), errResp.Msg)
}

func TestWriteUnauthorizedErrorRespAndLog(t *testing.T) {

	// Initialization
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
	rwr := httptest.NewRecorder()

	ctx := req.Context()

	expectedStatus := http.StatusUnauthorized
	errToLog := errors.New("error msg to log")
	authHeaderValue := "Bearer pipipipipi"
	authHeaderTypeExpected := "Basic"
	expectedErrResp := errors.Errorf(internal.ErrMsgFmtAuthorization, authHeaderTypeExpected)

	// Execution
	testWriterLogger.WriteUnauthorizedErrorRespAndLog(ctx, rwr, errToLog, authHeaderValue, authHeaderTypeExpected)

	// Check
	resultStatus := rwr.Result().StatusCode
	assert.Equal(t, expectedStatus, resultStatus)

	var errResp mock.ResponseMockError
	err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
	assert.Nil(t, err)

	assert.Equal(t, expectedErrResp.Error(), errResp.Msg)

	headerAuthenticate := rwr.Header().Get("WWW-Authenticate")
	assert.Equal(t, "Basic", headerAuthenticate)
}

func TestWriteInternalErrorRespAndLogPanic(t *testing.T) {

	// Initialization
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
	rwr := httptest.NewRecorder()

	ctx := req.Context()

	expectedStatus := http.StatusInternalServerError
	errToLog := errors.New("error msg to log")

	// Execution
	testWriterLogger.WriteInternalErrorRespAndLogPanic(ctx, rwr, errToLog)

	// Check
	resultStatus := rwr.Result().StatusCode
	assert.Equal(t, expectedStatus, resultStatus)

	var errResp mock.ResponseMockError
	err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
	assert.Nil(t, err)

	assert.Equal(t, internal.ErrRespInternalUnexpected.Error(), errResp.Msg)
}
