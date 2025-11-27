package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds application configuration.
type Config struct {
	DBHost     string `env:"DB_HOST"              envDefault:"localhost"`
	DBPort     int    `env:"DB_PORT"              envDefault:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	DBSSLMode  string `env:"DB_SSLMODE"           envDefault:"disable"`
	JWTSecret  string `env:"JWT_SECRET,required"`
	Port       int    `env:"PORT"                 envDefault:"8080"`
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

// DatabaseURL returns the PostgreSQL connection string.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}
