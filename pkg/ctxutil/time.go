package ctxutil

import (
	"context"
	"time"
)

type (
	timeKey      struct{}
	timeProvider interface {
		Now() time.Time
	}
)

func ContextWithTimeProvider(ctx context.Context, tp timeProvider) context.Context {
	return context.WithValue(ctx, timeKey{}, tp)
}

func Now(ctx context.Context) time.Time {
	var tp timeProvider = &timeProviderImpl{}
	if found, ok := ctx.Value(timeKey{}).(timeProvider); ok {
		tp = found
	}

	return tp.Now()
}

func getTimeProvider(ctx context.Context) timeProvider {
	if tp, ok := ctx.Value(timeKey{}).(timeProvider); ok {
		return tp
	}

	return nil
}

type timeProviderImpl struct{}

func (tp *timeProviderImpl) Now() time.Time {
	return time.Now()
}
