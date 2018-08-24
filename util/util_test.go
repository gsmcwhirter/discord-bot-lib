package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequestContext(t *testing.T) {
	ctx := NewRequestContext()
	rid := ctx.Value(ContextKey("request_id"))
	if assert.NotNil(t, rid) {
		rid2, ok := rid.(string)
		assert.True(t, ok)
		assert.NotEqual(t, rid2, "")
	}

	assert.True(t, true)
}
