package ctxutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/harusys/super-shiharai-kun/pkg/ctxutil"
	"github.com/harusys/super-shiharai-kun/pkg/ctxutil/ctxutiltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupEnv(t *testing.T) {
	ctx := context.Background()

	t.Setenv("ENV", "local")

	got, ok := ctxutil.LookupEnv(ctx, ctxutil.EnvKeyEnvName)
	require.True(t, ok)
	assert.Equal(t, "local", got)
}

func TestNow(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	got := ctxutil.Now(ctx)
	assert.NotEqual(t, 0, got.Unix())
}

func TestClone(t *testing.T) {
	t.Parallel()

	// Create a context with test providers
	fixedTime := time.Date(2000, 1, 1, 9, 0, 0, 0, time.UTC)
	provider := &ctxutiltest.TestContextProvider{
		CurrentTime: &fixedTime,
	}
	ctx := ctxutiltest.TestContext(provider)

	// Clone the context
	clonedCtx := ctxutil.Clone(ctx)

	// Verify that time provider is preserved
	originalTime := ctxutil.Now(ctx)
	clonedTime := ctxutil.Now(clonedCtx)
	assert.Equal(t, originalTime, clonedTime, "Time provider should be cloned correctly")
}

func TestClone_WithoutProviders(t *testing.T) {
	t.Parallel()

	// Create a context without custom providers
	ctx := context.Background()

	// Clone should not panic and should return a valid context
	clonedCtx := ctxutil.Clone(ctx)

	// Should be able to use the cloned context
	now := ctxutil.Now(clonedCtx)

	// Verify default providers work
	assert.NotEqual(t, 0, now.Unix(), "Default time provider should work")
}
