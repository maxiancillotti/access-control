package mock

import (
	"context"
	"encoding/json"
	"net/http"
)

type PresenterMock struct{}

type ResponseMockSuccess struct {
	Msg string `json:"response"`
}

type ResponseMockError struct {
	Msg string `json:"response"`
}

// MIDDLEWARE HANDLER
func (p *PresenterMock) CheckPresentationHeaders(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		// No hago nada

		next.ServeHTTP(rw, req)
	})
}

func (p *PresenterMock) SuccessResp(ctx context.Context, rw http.ResponseWriter, status int, sucessMsg string) {

	responseBody := ResponseMockSuccess{sucessMsg}
	p.PresentResponse(ctx, rw, status, responseBody)
}

func (p *PresenterMock) ErrorResp(ctx context.Context, rw http.ResponseWriter, status int, errorResp error) {

	responseBody := ResponseMockError{errorResp.Error()}
	p.PresentResponse(ctx, rw, status, responseBody)
}

func (p *PresenterMock) PresentResponse(ctx context.Context, rw http.ResponseWriter, status int, responseBody interface{}) {

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(responseBody)
}
