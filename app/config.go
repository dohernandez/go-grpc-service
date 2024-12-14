package app

import (
	"github.com/dohernandez/go-grpc-service/config"
	"github.com/dohernandez/go-grpc-service/database"
	"github.com/dohernandez/go-grpc-service/logger"
)

// Config represents config with variables needed for an app.
type Config struct {
	config.Config

	Environment    string `envconfig:"ENVIRONMENT" default:"dev"`
	ServiceName    string `envconfig:"SERVICE_NAME"`
	AppGRPCPort    int    `envconfig:"APP_GRPC_PORT" default:"8000"`
	AppRESTPort    int    `envconfig:"APP_REST_PORT" default:"8080"`
	AppMetricsPort int    `envconfig:"APP_METRICS_PORT" default:"8080"`
	AppHealthPort  int    `envconfig:"APP_HEALTH_PORT" default:"8081"`

	Database database.Config
	Logger   logger.Config
}
