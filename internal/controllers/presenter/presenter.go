package presenter

import (
	"context"
	"net/http"
)

type Presenter interface {
	PresentResponse(ctx context.Context, resp http.ResponseWriter, status int, responseBody interface{})
	SuccessResp(ctx context.Context, rw http.ResponseWriter, status int, sucessMsg string)
	ErrorResp(ctx context.Context, resp http.ResponseWriter, status int, errorResp error)
}

// Middlewares to run before the call to the main handlers
type PresenterMiddleware interface {
	CheckPresentationHeaders(next http.HandlerFunc) http.HandlerFunc
}
