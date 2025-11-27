package ctxutiltest

import (
	"context"
	"time"

	"github.com/harusys/super-shiharai-kun/pkg/ctxutil"
)

func TestContext(p *TestContextProvider) context.Context {
	ctx := context.Background()
	ctx = ctxutil.ContextWithEnvProvider(ctx, p)
	ctx = ctxutil.ContextWithTimeProvider(ctx, p)
	ctx = ctxutil.ContextWithPasswordHashProvider(ctx, p)

	return ctx
}

type TestContextProvider struct {
	EnvVars      map[ctxutil.EnvKey]string
	CurrentTime  *time.Time
	PasswordHash *string
}
