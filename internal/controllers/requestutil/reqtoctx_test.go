package requestutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssignRequestIDToCtx(t *testing.T) {

	ctx, reqID := AssignRequestIDToCtx(context.Background())

	assert.NotNil(t, ctx)
	assert.NotEmpty(t, reqID)

	ctxReqID, ok := ctx.Value(contextKeyRequestID).(string)

	assert.True(t, ok)
	assert.Equal(t, reqID, ctxReqID)
}

func TestGetRequestIDFromCtx(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		// Initialization
		expectedReqID := "requestIDxxx"

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, expectedReqID)

		// Execution
		reqID := GetRequestIDFromCtx(ctx)

		// Check
		assert.Equal(t, expectedReqID, reqID)
	})

	t.Run("Zero value return. Value type error.", func(t *testing.T) {

		// Initialization
		inputReqID := 123 // not string as expected by func

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestID, inputReqID)

		// Execution
		reqID := GetRequestIDFromCtx(ctx)

		// Check
		assert.Empty(t, reqID)
	})

	t.Run("Zero value return. Empty value error.", func(t *testing.T) {

		// Execution
		reqID := GetRequestIDFromCtx(context.Background())

		// Check
		assert.Empty(t, reqID)
	})

}

func TestAssignRequestBodyToCtx(t *testing.T) {

	expectedBody := []byte(`{"key":"value}`)

	ctx := AssignRequestBodyToCtx(context.Background(), expectedBody)

	assert.NotNil(t, ctx)

	ctxReqBody, ok := ctx.Value(contextKeyRequestBody).([]byte)

	assert.True(t, ok)
	assert.Equal(t, expectedBody, ctxReqBody)
}

func TestGetRequestBodyFromCtx(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		// Initialization
		expectedBody := []byte(`{"key":"value}`)

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestBody, expectedBody)

		// Execution
		ctxBody := GetRequestBodyFromCtx(ctx)

		// Check
		assert.Equal(t, expectedBody, ctxBody)
	})

	t.Run("Zero value return. Value type error.", func(t *testing.T) {

		// Initialization
		inputBody := 123

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyRequestBody, inputBody)

		// Execution
		ctxBody := GetRequestBodyFromCtx(ctx)

		// Check
		assert.Nil(t, ctxBody)
	})

	t.Run("Zero value return. Empty value error.", func(t *testing.T) {

		// Execution
		ctxBody := GetRequestBodyFromCtx(context.Background())

		// Check
		assert.Nil(t, ctxBody)
	})

}
