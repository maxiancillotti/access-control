package middlewares

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/maxiancillotti/access-control/internal/controllers/requestutil"
	"github.com/maxiancillotti/logger"
)

// RESPONSE WRITER IMPLEMENTATION WITH statusCode and Body saving for log

// This ResponseWritter implementation will not support interfaces CloseNotifier, Flusher, Hijacker, or Pusher
type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.buf.Write(body)
	return w.ResponseWriter.Write(body)
}

func (w *LogResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *LogResponseWriter) BodyString() string {
	return w.buf.String()
}

// MIDDLEWARE

type LoggingMiddleware interface {
	LogInAndOut(next http.HandlerFunc) http.HandlerFunc
}

type loggingMiddleware struct {
	logger logger.Logger
}

func NewLoggingMiddleware(logger logger.Logger) LoggingMiddleware {
	return &loggingMiddleware{
		logger: logger,
	}
}

// TO-DO MAYBE: DEFINE INTERFACE AND PASS THROUGH CONSTRUCTOR
type Logger interface{}

func (m *loggingMiddleware) LogInAndOut(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		// START
		startTime := time.Now()

		// REQUEST ID
		ctx, requestID := requestutil.AssignRequestIDToCtx(req.Context())
		// req = req.Clone(ctx) //deep copy
		req = req.WithContext(ctx) // shallow copy

		// INIT LOG
		defer m.logger.Flush()
		m.logger.Info("incoming_request",
			requestID,
			m.logger.StringField("request_URI", req.RequestURI),
			m.logger.StringField("method", req.Method),
			m.logger.StringField("remote_address", req.RemoteAddr),
			m.logger.StringField("host_server", req.Host),
		)

		// REQUEST BODY READ FOR LOG
		debugEnabled := requestutil.IsDebugDataEnabled(req)
		if debugEnabled {
			reqBodyBytes, err := ioutil.ReadAll(req.Body)
			if err == nil || err == io.EOF {
				ctx = requestutil.AssignRequestBodyToCtx(ctx, reqBodyBytes)
				// req = req.Clone(ctx) //deep copy
				req = req.WithContext(ctx) // shallow copy
			}
		}

		// RESPONSE WRITER IMPLEMENTACION THAT SAVES STATUS AND RESPONSE BODY
		lrw := NewLogResponseWriter(rw)

		// NEXT HANDLER
		next.ServeHTTP(lrw, req) // Overwritten req

		// FINAL DEFAULT LOG
		m.logger.Info("response_sent",
			requestID,
			m.logger.Field("time_elapsed", time.Since(startTime).String()),
			m.logger.Field("status", lrw.statusCode),
		)

		// DEBUG DATA LOG FOR ERRORS
		if lrw.statusCode > 299 {
			if debugEnabled {

				reqBodyBytes := requestutil.GetRequestBodyFromCtx(req.Context())

				req.Header.Set("Authorization", "value removed for log") // shouldn't log this

				m.logger.Debug("processed_data",
					requestID,
					nil,
					lrw.statusCode,
					m.logger.Field("request_headers", req.Header),
					m.logger.StringField("request_body", string(reqBodyBytes)),
					m.logger.Field("response_headers", lrw.Header()),
					m.logger.StringField("response_body", lrw.BodyString()),
				)
			}
		}
	})
}
