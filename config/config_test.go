package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/dohernandez/go-grpc-service/config"
	"github.com/stretchr/testify/require"
)

var cfg = &config.Config{
	ServiceName:    "test-server",
	AppGRPCPort:    0,
	AppRESTPort:    0,
	AppMetricsPort: 0,
	Environment:    "test",
}

func TestGetConfig_EnvSuccessfully(t *testing.T) {
	t.Parallel()

	// os.Setenv() can be replaced by `t.Setenv()` but can't be used in with t.Parallel()
	require.NoError(t, os.Setenv("SERVICE_NAME", "test-server"))
	require.NoError(t, os.Setenv("APP_GRPC_PORT", "0"))
	require.NoError(t, os.Setenv("APP_REST_PORT", "0"))
	require.NoError(t, os.Setenv("APP_METRICS_PORT", "0"))
	require.NoError(t, os.Setenv("ENVIRONMENT", "test"))

	got, err := config.GetConfig()
	require.NoError(t, err, "GetConfig() error = %v", err)

	if !reflect.DeepEqual(got, cfg) {
		t.Errorf("GetConfig() got = %v, want %v", got, cfg)
	}
}

func TestGetConfig_FileSuccessfully(t *testing.T) {
	t.Parallel()

	require.NoError(t, config.WithEnvFiles("./testdata/.env.template"))

	got, err := config.GetConfig()
	require.NoError(t, err, "GetConfig() error = %v", err)

	if !reflect.DeepEqual(got, cfg) {
		t.Errorf("GetConfig() got = %v, want %v", got, cfg)
	}
}
