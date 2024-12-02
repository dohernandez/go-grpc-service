package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config represents config with variables needed for an app.
type Config struct {
	Environment string `envconfig:"ENVIRONMENT" default:"dev"`
}

// LoadConfig load config, filled from environment variables.
func LoadConfig(c any) error {
	return envconfig.Process("", c)
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
