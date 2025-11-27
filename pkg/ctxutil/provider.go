package ctxutil

import (
	"context"
	"time"
)

func WithProvider(ctx context.Context) context.Context {
	ctx = ContextWithEnvProvider(ctx, &envProviderImpl{})
	ctx = ContextWithTimeProvider(ctx, &timeProviderImpl{})
	ctx = ContextWithPasswordHashProvider(ctx, &passwordHashProviderImpl{})

	return ctx
}

// Clone creates a new context with the same providers from the source context
// This is useful when you need to create a background context with the same
// providers (e.g., for testing with fixed time/UUID)
func Clone(ctx context.Context) context.Context {
	newCtx := context.Background()

	// Copy all providers from source context to new context
	if envProvider := getEnvProvider(ctx); envProvider != nil {
		newCtx = ContextWithEnvProvider(newCtx, envProvider)
	}

	if timeProvider := getTimeProvider(ctx); timeProvider != nil {
		newCtx = ContextWithTimeProvider(newCtx, timeProvider)
	}

	if passwordHashProvider := getPasswordHashProvider(ctx); passwordHashProvider != nil {
		newCtx = ContextWithPasswordHashProvider(newCtx, passwordHashProvider)
	}

	return newCtx
}

// WithTimeout wraps context.WithTimeout and preserves providers from the source context
func WithTimeout(
	parent context.Context,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
	// Create timeout context
	ctx, cancel := context.WithTimeout(parent, timeout)

	// Providers are already in parent context, so just return the timeout context
	return ctx, cancel
}

// WithTimeoutClone creates a background context with timeout, copying providers from the source context
// This is useful for background tasks that need to be independent of the parent context's cancellation
func WithTimeoutClone(
	parent context.Context,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
	// Clone providers to background context
	ctx := Clone(parent)

	// Add timeout
	return context.WithTimeout(ctx, timeout)
}
