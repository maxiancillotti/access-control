package writerloggermock

import (
	"context"
	"errors"
	"net/http"

	"github.com/maxiancillotti/access-control/internal/controllers/presenter"
	"github.com/maxiancillotti/access-control/internal/controllers/responseutil"

	"github.com/maxiancillotti/logger"
)

var (
	TestErrLogHeaderKey    = "TEST-ErrLog"
	TestFailingHKHeaderKey = "TEST-failingHeaderKey"
	TestFailingHVHeaderKey = "TEST-failingHeaderValue"
)

type writerLoggerMock struct {
	presenter presenter.Presenter
	logger    logger.Logger
}

func NewWriterLoggerMock(pstrMock presenter.Presenter, lgrMock logger.Logger) responseutil.WriterLogger {
	return &writerLoggerMock{pstrMock, lgrMock}
}

func (wl *writerLoggerMock) WriteErrorRespAndLog(ctx context.Context, rw http.ResponseWriter, status int, respErr error, logErr error) {

	rw = wl.logWriteMock(rw, logErr, "", "")

	wl.presenter.ErrorResp(ctx, rw, status, respErr)
}

func (wl *writerLoggerMock) WritePresentationHeaderErrorRespAndLog(ctx context.Context, rw http.ResponseWriter, err error, headerKey string, headerValue string) {

	rw = wl.logWriteMock(rw, err, headerKey, headerValue)

	status := http.StatusNotAcceptable
	wl.presenter.ErrorResp(ctx, rw, status, err)
}

func (wl *writerLoggerMock) WriteUnauthorizedErrorRespAndLog(ctx context.Context, rw http.ResponseWriter, errLog error, authHeader string, expectedHeaderType string) {

	rw = wl.logWriteMock(rw, errLog, "Authorization", authHeader)

	status := http.StatusUnauthorized
	rw.Header().Add("WWW-Authenticate", expectedHeaderType)
	wl.presenter.ErrorResp(ctx, rw, status, errors.New("ERROR_UNAUTHORIZED"))
}

func (wl *writerLoggerMock) WriteInternalErrorRespAndLogPanic(ctx context.Context, rw http.ResponseWriter, errLog error) {

	rw = wl.logWriteMock(rw, errLog, "", "")

	status := http.StatusInternalServerError
	wl.presenter.ErrorResp(ctx, rw, status, errors.New("ERROR_INTERNAL"))
}

func (wl *writerLoggerMock) logWriteMock(rw http.ResponseWriter, err error, failingHeaderKey string, failingHeaderValue string) http.ResponseWriter {

	rw.Header().Add(TestErrLogHeaderKey, err.Error())

	if failingHeaderKey != "" {
		rw.Header().Add(TestFailingHKHeaderKey, failingHeaderKey)
		rw.Header().Add(TestFailingHVHeaderKey, failingHeaderValue)
	}

	return rw
}
