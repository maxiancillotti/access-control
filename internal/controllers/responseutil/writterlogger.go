package responseutil

import (
	"context"
	"net/http"

	"github.com/maxiancillotti/access-control/internal/controllers/internal"
	"github.com/maxiancillotti/access-control/internal/controllers/presenter"
	"github.com/maxiancillotti/access-control/internal/controllers/requestutil"

	"github.com/maxiancillotti/logger"

	"github.com/pkg/errors"
)

type WriterLogger interface {
	WriteErrorRespAndLog(ctx context.Context, rw http.ResponseWriter, status int, respErr error, logErr error)
	WritePresentationHeaderErrorRespAndLog(ctx context.Context, rw http.ResponseWriter, err error, headerKey string, headerValue string)
	WriteUnauthorizedErrorRespAndLog(ctx context.Context, rw http.ResponseWriter, errLog error, authHeader string, expectedHeaderType string)
	WriteInternalErrorRespAndLogPanic(ctx context.Context, rw http.ResponseWriter, errLog error)
}

type writerLogger struct {
	presenter presenter.Presenter
	logger    logger.Logger
}

func NewWriterLogger(presenter presenter.Presenter, logger logger.Logger) WriterLogger {
	return &writerLogger{
		presenter: presenter,
		logger:    logger,
	}
}

// GENERAL

func (wl *writerLogger) WriteErrorRespAndLog(ctx context.Context, rw http.ResponseWriter, status int, respErr error, logErr error) {

	wl.presenter.ErrorResp(ctx, rw, status, respErr)

	if status > 499 {
		wl.logger.Error(logMsgKeyInternalErr, requestutil.GetRequestIDFromCtx(ctx), logErr, status)
	} else {
		wl.logger.Warn(logMsgKeyClientErr, requestutil.GetRequestIDFromCtx(ctx), logErr, status)
	}
	wl.logger.Flush()
}

// SPECIFIC

func (wl *writerLogger) WritePresentationHeaderErrorRespAndLog(ctx context.Context, rw http.ResponseWriter, err error, headerKey string, headerValue string) {

	status := http.StatusNotAcceptable
	wl.presenter.ErrorResp(ctx, rw, status, err)

	var failingHeader = make(http.Header)
	failingHeader.Set(headerKey, headerValue)

	wl.logger.Warn(logMsgKeyClientErr, requestutil.GetRequestIDFromCtx(ctx), err, status, wl.logger.Field(logKeyRequestHeaders, failingHeader))
	wl.logger.Flush()
}

func (wl *writerLogger) WriteUnauthorizedErrorRespAndLog(ctx context.Context, rw http.ResponseWriter, errLog error, authHeaderVal string, expectedHeaderType string) {

	status := http.StatusUnauthorized
	rw.Header().Add("WWW-Authenticate", expectedHeaderType)
	wl.presenter.ErrorResp(ctx, rw, status, errors.Errorf(internal.ErrMsgFmtAuthorization, expectedHeaderType))

	var failingHeader = make(http.Header)
	failingHeader.Set("Authorization", authHeaderVal)

	wl.logger.Warn(logMsgKeyClientErr, requestutil.GetRequestIDFromCtx(ctx), errLog, status, wl.logger.Field(logKeyRequestHeaders, failingHeader))
	wl.logger.Flush()
}

func (wl *writerLogger) WriteInternalErrorRespAndLogPanic(ctx context.Context, rw http.ResponseWriter, errLog error) {

	status := http.StatusInternalServerError
	wl.presenter.ErrorResp(ctx, rw, status, internal.ErrRespInternalUnexpected)

	wl.logger.Error(logMsgKeyPANIC, requestutil.GetRequestIDFromCtx(ctx), errLog, status)
	wl.logger.Flush()
}
