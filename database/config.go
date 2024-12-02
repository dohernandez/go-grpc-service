package database

import "time"

// Config represents the DB configuration fields and values.
type Config struct {
	DSN          string        `envconfig:"DATABASE_DSN" required:"true"`
	MaxLifetime  time.Duration `envconfig:"MAX_LIFETIME" default:"4h"`
	MaxIdleConns int           `envconfig:"MAX_IDLE_CONNECTIONS" default:"20"`
	MaxOpenConns int           `envconfig:"MAX_OPEN_CONNECTIONS" default:"20"`

	DriverName string
}
