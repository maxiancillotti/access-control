package requestutil

import (
	"net/http"
	"strconv"
)

// IsDebugDataEnabled returns bool when the custom header "X-Debug-Enabled" is
// present and with a value of "true" in the request.
// The caller then will enable debug level logging
func IsDebugDataEnabled(req *http.Request) bool {

	debugEnabled, err := strconv.ParseBool(req.Header.Get("X-Debug-Enabled"))
	return err == nil && debugEnabled
}
