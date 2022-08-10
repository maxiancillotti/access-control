package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testLoggingMDW = NewLoggingMiddleware(lgrMock)
)

func TestLogIntAndOut(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		// Initialization
		req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)
		rwr := httptest.NewRecorder()

		req.Header.Set("Authorization", "Basic user:password")

		// Execution
		testLoggingMDW.LogInAndOut(handlerMockforMDW).ServeHTTP(rwr, req)
	})

	t.Run("Success. Debug mode enabled", func(t *testing.T) {

		// Initialization

		reqBody := []byte(`{"key":"value"}`)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", bytes.NewBuffer(reqBody))
		rwr := httptest.NewRecorder()

		req.Header.Set("Authorization", "Basic user:password")
		req.Header.Set("X-Debug-Enabled", "true")

		// Execution
		testLoggingMDW.LogInAndOut(handlerMockforMDW).ServeHTTP(rwr, req)
	})

	t.Run("Status error. Debug mode enabled", func(t *testing.T) {

		// Initialization

		reqBody := []byte(`{"key":"value"}`)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", bytes.NewBuffer(reqBody))
		rwr := httptest.NewRecorder()

		req.Header.Set("Authorization", "Basic user:password")
		req.Header.Set("X-Debug-Enabled", "true")

		// Execution

		// This new handler will write status 500
		testLoggingMDW.LogInAndOut(handlerMockforMDWStatus500).ServeHTTP(rwr, req)
	})
}
