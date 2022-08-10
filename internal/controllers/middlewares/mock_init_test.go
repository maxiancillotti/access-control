package middlewares

import (
	"net/http"

	"github.com/maxiancillotti/access-control/internal/mock"
	"github.com/maxiancillotti/access-control/internal/mock/writerloggermock"
)

var (
	pstrMock mock.PresenterMock
	lgrMock  = mock.LoggerMock
	// TODO: MOCK LOGGER TOO SO THE INTERNAL ERRORS CAN BE SEEN AND TESTED
	// And obviously, so the test doesn't depend on the logger too, duh.
	wrlgrMock = writerloggermock.NewWriterLoggerMock(&pstrMock, lgrMock)
)

func handlerMockforMDW(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func handlerMockforMDWStatus500(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
}
