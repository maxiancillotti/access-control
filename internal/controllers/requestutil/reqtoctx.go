package requestutil

import (
	"context"

	"github.com/google/uuid"
)

// AssignRequestIDToCtx will assign a brand new request ID to a http request
func AssignRequestIDToCtx(ctx context.Context) (context.Context, string) {

	reqID := uuid.New()
	return context.WithValue(ctx, contextKeyRequestID, reqID.String()), reqID.String()
}

// GetRequestIDFromCtx will get reqID from a http request and return it as a string
func GetRequestIDFromCtx(ctx context.Context) string {

	reqIDCtx := ctx.Value(contextKeyRequestID)

	if reqID, ok := reqIDCtx.(string); ok {
		return reqID
	}
	return ""
}

// AssignRequestBodyToCtx
func AssignRequestBodyToCtx(ctx context.Context, body []byte) context.Context {

	return context.WithValue(ctx, contextKeyRequestBody, body)
}

// GetRequestBodyFromCtx
func GetRequestBodyFromCtx(ctx context.Context) []byte {

	reqBody := ctx.Value(contextKeyRequestBody)

	if reqBody != nil {
		if ret, ok := reqBody.([]byte); ok {
			return ret
		}
	}
	return nil
}
