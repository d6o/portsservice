package appcontext

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestChecker_CheckContext_NotDone(t *testing.T) {
	checker := NewChecker()

	err := checker.CheckContext(context.Background())
	assert.NoError(t, err)
}

func TestChecker_CheckContext_Done(t *testing.T) {
	checker := NewChecker()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	// Wait for the context to be done
	time.Sleep(2 * time.Millisecond)

	err := checker.CheckContext(ctx)
	assert.Error(t, err)
}
