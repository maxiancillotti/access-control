package presenters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxiancillotti/access-control/internal/controllers/presenter"
	"github.com/maxiancillotti/access-control/internal/mock"
	"github.com/maxiancillotti/access-control/internal/mock/writerloggermock"
	"github.com/stretchr/testify/assert"
)

// TODO: MOCK LOGGER TOO SO THE INTERNAL ERRORS CAN BE SEEN AND TESTED
// And obviously, so the test doesn't depend on the logger too, duh.

var (
	pstrMock mock.PresenterMock
	lgrMock  = mock.LoggerMock

	wrlgrMock = writerloggermock.NewWriterLoggerMock(&pstrMock, lgrMock)

	testPresenterMDW presenter.PresenterMiddleware = NewJSONPresenterMDW(wrlgrMock)
)

func handlerMockforMDW(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

//////////////////////////////////////////////

func TestCheckPresentationHeaders(t *testing.T) {

	// Initialization

	type testCheckPresentationHeadersCase struct {
		name              string
		acceptHeaderValue string
		expectedStatus    int
		expectedErr       error
		expectedErrLogged error
	}

	testTable := make([]testCheckPresentationHeadersCase, 0)

	testTable = append(testTable, testCheckPresentationHeadersCase{
		name:              "Success. JSON and other mime types.",
		acceptHeaderValue: "application/xml, text/html, application/json",
		expectedStatus:    http.StatusOK,
		expectedErr:       nil,
	})

	testTable = append(testTable, testCheckPresentationHeadersCase{
		name:              "Success. Only JSON.",
		acceptHeaderValue: "application/json",
		expectedStatus:    http.StatusOK,
		expectedErr:       nil,
	})

	testTable = append(testTable, testCheckPresentationHeadersCase{
		name:              "Success. All types",
		acceptHeaderValue: "*/*",
		expectedStatus:    http.StatusOK,
		expectedErr:       nil,
	})

	testTable = append(testTable, testCheckPresentationHeadersCase{
		name:              "Success. Empty type.",
		acceptHeaderValue: "",
		expectedStatus:    http.StatusOK,
		expectedErr:       nil,
	})

	testTable = append(testTable, testCheckPresentationHeadersCase{
		name:              "Error. Doesn't accept JSON.",
		acceptHeaderValue: "application/xml",
		expectedStatus:    http.StatusNotAcceptable,
		expectedErr:       errRespCheckHeadersContentTypeNotSupported,
		expectedErrLogged: errRespCheckHeadersContentTypeNotSupported,
	})

	for _, test := range testTable {

		t.Run(test.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
			rwr := httptest.NewRecorder()

			req.Header.Set("Accept", test.acceptHeaderValue)

			// Execution
			testPresenterMDW.CheckPresentationHeaders(handlerMockforMDW).ServeHTTP(rwr, req)

			// Check
			statusCode := rwr.Result().StatusCode
			assert.Equal(t, test.expectedStatus, statusCode)

			if test.expectedErr != nil {
				var errResp mock.ResponseMockError
				err := json.NewDecoder(rwr.Result().Body).Decode(&errResp)
				assert.Nil(t, err)

				assert.Equal(t, test.expectedErr.Error(), errResp.Msg)

				failingHeaderKeyLogged := rwr.Result().Header.Get(writerloggermock.TestFailingHKHeaderKey)
				assert.Contains(t, failingHeaderKeyLogged, acceptHeaderKey)

				failingHeaderValueLogged := rwr.Result().Header.Get(writerloggermock.TestFailingHVHeaderKey)
				assert.Contains(t, failingHeaderValueLogged, test.acceptHeaderValue)

				errLogged := rwr.Result().Header.Get(writerloggermock.TestErrLogHeaderKey)
				assert.Contains(t, errLogged, test.expectedErrLogged.Error())
			}
		})

	}
}
