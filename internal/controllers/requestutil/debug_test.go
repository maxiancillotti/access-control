package requestutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDebugDataEnabled(t *testing.T) {

	type testCase struct {
		name           string
		headerKey      string
		headerValue    string
		expectedResult bool
	}

	table := make([]testCase, 0)

	table = append(table, testCase{
		name:           "Success. True.",
		headerKey:      "X-Debug-Enabled",
		headerValue:    "true",
		expectedResult: true,
	})

	table = append(table, testCase{
		name:           "False. Header value false",
		headerKey:      "X-Debug-Enabled",
		headerValue:    "false",
		expectedResult: false,
	})

	table = append(table, testCase{
		name:           "False. Err empty header value",
		headerKey:      "X-Debug-Enabled",
		headerValue:    "",
		expectedResult: false,
	})

	table = append(table, testCase{
		name:           "False. Err empty header key",
		headerKey:      "",
		headerValue:    "",
		expectedResult: false,
	})

	for _, test := range table {

		t.Run(test.name, func(t *testing.T) {

			// Initialization
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8001/authtoken", nil)

			if test.headerKey != "" {
				req.Header.Set(test.headerKey, test.headerValue)
			}

			// Execution
			isEnabled := IsDebugDataEnabled(req)

			// Check
			assert.Equal(t, test.expectedResult, isEnabled)

		})
	}
}
