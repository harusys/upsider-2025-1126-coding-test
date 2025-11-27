package ctxutiltest

import (
	"context"
	"os"

	"github.com/harusys/super-shiharai-kun/pkg/ctxutil"
)

func (p *TestContextProvider) LookupEnv(_ context.Context, key ctxutil.EnvKey) (string, bool) {
	if v, ok := p.EnvVars[key]; ok {
		return v, true
	}

	return os.LookupEnv(string(key)) //nolint:forbidigo // for test purpose
}

func (p *TestContextProvider) SetEnvVar(k ctxutil.EnvKey, v string) {
	if p.EnvVars == nil {
		p.EnvVars = make(map[ctxutil.EnvKey]string)
	}

	p.EnvVars[k] = v
}
