package ctxutiltest

import (
	"context"

	"github.com/harusys/super-shiharai-kun/pkg/ctxutil"
)

func (p *TestContextProvider) HashPassword(password string, cost int) (string, error) {
	if p.PasswordHash != nil {
		return *p.PasswordHash, nil
	}

	return ctxutil.HashPassword(context.TODO(), password, cost)
}
