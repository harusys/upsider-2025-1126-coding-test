package ctxutil

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type (
	passwordHashKey      struct{}
	passwordHashProvider interface {
		HashPassword(password string, cost int) (string, error)
	}
)

func ContextWithPasswordHashProvider(
	ctx context.Context,
	php passwordHashProvider,
) context.Context {
	return context.WithValue(ctx, passwordHashKey{}, php)
}

func HashPassword(ctx context.Context, password string, cost int) (string, error) {
	var php passwordHashProvider = &passwordHashProviderImpl{}
	if found, ok := ctx.Value(passwordHashKey{}).(passwordHashProvider); ok {
		php = found
	}

	return php.HashPassword(password, cost)
}

func getPasswordHashProvider(ctx context.Context) passwordHashProvider {
	if php, ok := ctx.Value(passwordHashKey{}).(passwordHashProvider); ok {
		return php
	}

	return nil
}

type passwordHashProviderImpl struct{}

func (php *passwordHashProviderImpl) HashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
