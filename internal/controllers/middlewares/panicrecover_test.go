package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxiancillotti/access-control/internal/mock"
	"github.com/maxiancillotti/access-control/internal/mock/writerloggermock"
	"github.com/stretchr/testify/assert"
)

var (
	testPanicRecoverMDW = NewPanicRecoverMiddleware(wrlgrMock)
)

func handlerMockforMDWWithPanic(rw http.ResponseWriter, req *http.Request) {
	panic("this is a panic")
}

func TestPanicRecover(t *testing.T) {

	// Initialization
	req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken/authorize", nil)
	rwr := httptest.NewRecorder()

	// Execution
	testPanicRecoverMDW.PanicRecover(handlerMockforMDWWithPanic).ServeHTTP(rwr, req)

	// Check
	statusCode := rwr.Result().StatusCode
	assert.Equal(t, http.StatusInternalServerError, statusCode)

	var successResp mock.ResponseMockSuccess
	err := json.NewDecoder(rwr.Result().Body).Decode(&successResp)
	assert.Nil(t, err)

	assert.Equal(t, "ERROR_INTERNAL", successResp.Msg)

	errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
	assert.Contains(t, errLogged, "this is a panic")
}
