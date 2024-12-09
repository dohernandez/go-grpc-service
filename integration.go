package go_grpc_service

import (
	"context"
	"github.com/bool64/dbdog"
	"github.com/bool64/sqluct"
	"github.com/cucumber/godog"
	"github.com/dohernandez/go-grpc-service/app"
	"github.com/dohernandez/go-grpc-service/test/feature"
	dbdogcleaner "github.com/dohernandez/go-grpc-service/test/feature/database"
	"github.com/nhatthm/clockdog"
	"testing"
	"time"
)

type FeaturesConfig struct {
	FeaturePath        string
	Locator            *app.Locator
	FeatureContextFunc func(t *testing.T, s *godog.ScenarioContext)
	Tables             map[string]any
}

func RunFeatures(t *testing.T, ctx context.Context, cfg *FeaturesConfig) {
	t.Helper()

	deps := cfg.Locator

	clock := clockdog.New()
	deps.ClockProvider = clock

	stop, errch := RunServicesTesting(t, ctx, deps)

	var baseRESTURL string

	select {
	case err := <-errch:
		if err != nil {
			t.Fatalf("failed to run service: %v", err)
		}
	case baseRESTURL = <-deps.GRPCRestService.AddrAssigned:
		break
	case <-time.After(timeout):
		t.Fatal("timeout waiting for service to start")
	}

	defer func() {
		// Stop the service.
		stop()
	}()

	local := feature.NewLocal(baseRESTURL)

	dbm := initDBManager(deps.Storage, cfg.Tables)
	dbmCleaner := initDBMCleaner(dbm)

	feature.RunFeatures(t, cfg.FeaturePath, func(_ *testing.T, s *godog.ScenarioContext) {
		local.RegisterSteps(s)

		dbm.RegisterSteps(s)
		dbmCleaner.RegisterSteps(s)

		clock.RegisterContext(s)

		cfg.FeatureContextFunc(t, s)
	})
}

func initDBManager(storage *sqluct.Storage, tables map[string]any) *dbdog.Manager {
	tableMapper := dbdog.NewTableMapper()

	dbm := dbdog.Manager{
		TableMapper: tableMapper,
	}

	dbm.Instances = map[string]dbdog.Instance{
		"postgres": {
			Storage: storage,
			Tables:  tables,
		},
	}

	return &dbm
}

func initDBMCleaner(dbm *dbdog.Manager) *dbdogcleaner.ManagerCleaner {
	return &dbdogcleaner.ManagerCleaner{
		Manager: dbm,
	}
}
