package presenters

import (
	"net/http"
	"strings"

	"github.com/maxiancillotti/access-control/internal/controllers/presenter"
	"github.com/maxiancillotti/access-control/internal/controllers/responseutil"
)

type jsonPresenterMDW struct {
	writerLogger responseutil.WriterLogger
}

func NewJSONPresenterMDW(wrlgr responseutil.WriterLogger) presenter.PresenterMiddleware {
	return &jsonPresenterMDW{
		writerLogger: wrlgr,
	}
}

// MIDDLEWARE HANDLER
func (m *jsonPresenterMDW) CheckPresentationHeaders(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		var acceptsJSON bool

		acceptHeader := strings.ToLower(req.Header.Get(acceptHeaderKey))
		acceptHeaderSplit := strings.Split(acceptHeader, ",")

		for _, v := range acceptHeaderSplit {
			v = strings.TrimSpace(v)

			if v == contentTypeJSON || v == "" || v == "*/*" {
				acceptsJSON = true
				break
			}
		}
		if !acceptsJSON {
			// status, respErr := http.StatusNotAcceptable, errRespCheckHeadersContentTypeNotSupported
			// p.ErrorResp(req.Context(), rw, status, respErr)
			// p.logger.Warn("client_error_response", controllers.GetRequestIDFromCtx(req.Context()), respErr, status, p.logger.StringField("request_header_accept", acceptHeader))
			// p.logger.Flush()
			m.writerLogger.WritePresentationHeaderErrorRespAndLog(req.Context(), rw, errRespCheckHeadersContentTypeNotSupported, acceptHeaderKey, acceptHeader)
			return
		}

		next.ServeHTTP(rw, req)
	})
}

// func (c *jsonPresenter) writePresentationHeaderErrorRespAndLog(ctx context.Context, rw http.ResponseWriter, err error, headerKey string, headerValue string) {

// 	status := http.StatusNotAcceptable
// 	rw.Header().Add("WWW-Authenticate", "Basic")
// 	c.ErrorResp(ctx, rw, status, err)

// 	var failingHeader = make(http.Header)
// 	failingHeader.Set(headerKey, headerValue)

// 	c.logger.Warn("client_error_response",
// 		requestutil.GetRequestIDFromCtx(ctx),
// 		err,
// 		status,
// 		//c.logger.StringField(fmt.Sprint("request_header_", headerKey), headerValue),
// 		/*
// 			c.logger.Field("request_headers",
// 				fmt.Sprintf(`{"%s":"%s"}`, headerKey, headerValue),
// 			),
// 		*/
// 		c.logger.Field("request_headers", failingHeader),
// 	)
// 	c.logger.Flush()
// }
