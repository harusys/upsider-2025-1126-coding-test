package infrastructure

import "time"

// Security constants.
const (
	// AccessTokenExpiry is the expiry duration for access tokens.
	AccessTokenExpiry = 15 * time.Minute
	// RefreshTokenExpiry is the expiry duration for refresh tokens.
	RefreshTokenExpiry = 7 * 24 * time.Hour
	// BcryptCost is the cost parameter for bcrypt hashing.
	BcryptCost = 12
)

// Database connection pool constants.
const (
	// DBMaxConns is the maximum number of connections in the pool.
	DBMaxConns = 25
	// DBMinConns is the minimum number of connections in the pool.
	DBMinConns = 5
)
