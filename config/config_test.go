package config_test

import (
	"os"
	"testing"

	"github.com/dohernandez/go-grpc-service/config"
	"github.com/stretchr/testify/require"
)

func TestGetConfig_EnvSuccessfully(t *testing.T) {
	t.Parallel()

	// os.Setenv() can be replaced by `t.Setenv()` but can't be used in with t.Parallel()
	require.NoError(t, os.Setenv("ENVIRONMENT", "test"))

	var got config.Config

	err := config.LoadConfig(&got)
	require.NoError(t, err)

	require.Equal(t, "test", got.Environment)
}

func TestGetConfig_FileSuccessfully(t *testing.T) {
	t.Parallel()

	require.NoError(t, config.WithEnvFiles("./testdata/.env.template"))

	var got config.Config

	err := config.LoadConfig(&got)
	require.NoError(t, err)

	require.Equal(t, "test", got.Environment)
}
