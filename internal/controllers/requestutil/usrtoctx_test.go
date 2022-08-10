package requestutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssignAuthenticatedUserIDToCxt(t *testing.T) {

	var expectedUserID uint = 123

	ctx := AssignAuthenticatedUserIDToCxt(context.Background(), expectedUserID)

	assert.NotNil(t, ctx)

	ctxUserID, ok := ctx.Value(contextKeyAuthBasicUserID).(uint)

	assert.True(t, ok)
	assert.Equal(t, expectedUserID, ctxUserID)
}

func TestGetAuthenticatedUserIDFromCtx(t *testing.T) {

	t.Run("Success", func(t *testing.T) {

		// Initialization
		var expectedUserID uint = 123

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyAuthBasicUserID, expectedUserID)

		// Execution
		ctxUserID := GetAuthenticatedUserIDFromCtx(ctx)

		// Check
		assert.Equal(t, expectedUserID, ctxUserID)
	})

	t.Run("Zero value return. Value type error.", func(t *testing.T) {

		// Initialization
		inputUserID := "123"

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKeyAuthBasicUserID, inputUserID)

		// Execution
		ctxUserID := GetAuthenticatedUserIDFromCtx(ctx)

		// Check
		assert.Zero(t, ctxUserID)
	})

	t.Run("Zero value return. Empty value error.", func(t *testing.T) {

		// Execution
		ctxUserID := GetAuthenticatedUserIDFromCtx(context.Background())

		// Check
		assert.Zero(t, ctxUserID)
	})

}
