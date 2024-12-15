package config

// RateLimiter represents rate limiting configuration.
type RateLimiter struct {
	Enabled bool `default:"false"`

	// Client configuration, used for client rate limiting.
	Client Limiter `split_words:"true"`

	// Server configuration, used for server rate limiting (all clients).
	Server Limiter `split_words:"true"`
}

// Limiter defines rate limiting configuration.
type Limiter struct {
	// RequestsPerSec number of requests per second.
	RequestsPerSec float64 `default:"10000"`
	// BurstLimit burst limit.
	BurstLimit int `default:"1000"`
}
