package ctxutiltest_test

import (
	"testing"
	"time"

	"github.com/harusys/super-shiharai-kun/pkg/ctxutil"
	"github.com/harusys/super-shiharai-kun/pkg/ctxutil/ctxutiltest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLookupEnv(t *testing.T) {
	t.Setenv("ENV", "LOCAL")

	// call EnvVar() without setting env
	got1, ok := ctxutil.LookupEnv(
		ctxutiltest.TestContext(&ctxutiltest.TestContextProvider{}),
		ctxutil.EnvKeyEnvName,
	)
	assert.True(t, ok)
	assert.Equal(t, "LOCAL", got1)

	// set env before construction
	p := &ctxutiltest.TestContextProvider{}
	p.SetEnvVar(ctxutil.EnvKeyEnvName, "PRODUCTION")

	ctx := ctxutiltest.TestContext(p)
	got2, ok := ctxutil.LookupEnv(ctx, ctxutil.EnvKeyEnvName)
	assert.True(t, ok)
	assert.Equal(t, "PRODUCTION", got2)

	// set env after construction
	p.SetEnvVar(ctxutil.EnvKeyEnvName, "STAGING")

	got3, ok := ctxutil.LookupEnv(ctx, ctxutil.EnvKeyEnvName)
	assert.True(t, ok)
	assert.Equal(t, "STAGING", got3)
}

func TestNow(t *testing.T) {
	t.Parallel()

	asiaTokyo, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	// call Now() without setting time
	got := ctxutil.Now(ctxutiltest.TestContext(&ctxutiltest.TestContextProvider{}))
	assert.NotEqual(t, 0, got.Unix())

	// set time befor construction
	p := &ctxutiltest.TestContextProvider{}
	p.SetAsiaTokyo(t, "2020-01-01 07:00:00")

	ctx := ctxutiltest.TestContext(p)
	want1 := time.Date(2020, 1, 1, 7, 0, 0, 0, asiaTokyo)
	got1 := ctxutil.Now(ctx)
	assert.Equal(t, want1, got1)

	// set time after construction
	p.SetAsiaTokyo(t, "2020-01-01 08:00:00")

	want2 := time.Date(2020, 1, 1, 8, 0, 0, 0, asiaTokyo)
	got2 := ctxutil.Now(ctx)
	assert.Equal(t, want2, got2)
}
