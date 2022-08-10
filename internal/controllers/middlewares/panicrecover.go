package middlewares

import (
	"fmt"
	"net/http"

	"github.com/maxiancillotti/access-control/internal/controllers/responseutil"
	"github.com/pkg/errors"
)

type PanicRecoverMiddleware interface {
	PanicRecover(next http.HandlerFunc) http.HandlerFunc
}

func NewPanicRecoverMiddleware(writerLogger responseutil.WriterLogger) PanicRecoverMiddleware {
	return &panicRecoverMiddleware{writerLogger}
}

type panicRecoverMiddleware struct {
	writerLogger responseutil.WriterLogger
}

func (m *panicRecoverMiddleware) PanicRecover(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		defer func() {
			if pncErr := recover(); pncErr != nil {
				err := errors.New(fmt.Sprint(pncErr))
				m.writerLogger.WriteInternalErrorRespAndLogPanic(req.Context(), rw, err)
			}
		}()

		next.ServeHTTP(rw, req)
	})
}
