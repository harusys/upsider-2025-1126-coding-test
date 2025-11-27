package ctxutil

import (
	"context"
	"os"
)

type EnvKey string

const (
	EnvKeyDataDir EnvKey = "DATA_DIR"
	EnvKeyEnvName EnvKey = "ENV"
	EnvKeyPort    EnvKey = "PORT"
)

type (
	envKey      struct{}
	envProvider interface {
		LookupEnv(ctx context.Context, key EnvKey) (string, bool)
	}
)

func ContextWithEnvProvider(ctx context.Context, tp envProvider) context.Context {
	return context.WithValue(ctx, envKey{}, tp)
}

func LookupEnv(ctx context.Context, key EnvKey) (string, bool) {
	var ep envProvider = &envProviderImpl{}
	if found, ok := ctx.Value(envKey{}).(envProvider); ok {
		ep = found
	}

	return ep.LookupEnv(ctx, key)
}

func getEnvProvider(ctx context.Context) envProvider {
	if ep, ok := ctx.Value(envKey{}).(envProvider); ok {
		return ep
	}

	return nil
}

type envProviderImpl struct{}

func (tp *envProviderImpl) LookupEnv(_ context.Context, key EnvKey) (string, bool) {
	return os.LookupEnv(string(key))
}
