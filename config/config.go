package config

import (
	"github.com/dohernandez/go-grpc-service/database"
	"github.com/dohernandez/go-grpc-service/logger"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config represents config with variables needed for an app.
type Config struct {
	ServiceName    string `envconfig:"SERVICE_NAME"`
	AppGRPCPort    int    `envconfig:"APP_GRPC_PORT" default:"8000"`
	AppRESTPort    int    `envconfig:"APP_REST_PORT" default:"8080"`
	AppMetricsPort int    `envconfig:"APP_METRICS_PORT" default:"8080"`
	Environment    string `envconfig:"ENVIRONMENT" default:"dev"`

	Log logger.Config
	DB  database.Config
}

// GetConfig returns service config, filled from environment variables.
func GetConfig() (*Config, error) {
	var c Config

	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// IsDev returns true for development environment.
func (c *Config) IsDev() bool {
	return len(c.Environment) > 2 && strings.ToLower(c.Environment[0:3]) == "dev"
}

// IsTest returns true for testing environment.
func (c *Config) IsTest() bool {
	return len(c.Environment) > 2 && strings.ToLower(c.Environment[0:4]) == "test"
}

// WithEnvFiles populates env vars from provided files.
//
// It returns an error if file does not exist.
func WithEnvFiles(files ...string) error {
	var found []string

	for _, f := range files {
		if fileExists(f) {
			found = append(found, f)
		}
	}

	if len(found) == 0 {
		return nil
	}

	return godotenv.Load(files...)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
