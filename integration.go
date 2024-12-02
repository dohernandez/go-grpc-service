package go_grpc_service

import (
	"context"
	"testing"
	"time"

	"github.com/bool64/ctxd"
	"github.com/bool64/dbdog"
	"github.com/bool64/httpdog"
	"github.com/bool64/sqluct"
	"github.com/cucumber/godog"
	"github.com/dohernandez/go-grpc-service/app"
	"github.com/dohernandez/go-grpc-service/must"
	"github.com/dohernandez/go-grpc-service/test/feature"
	dbdogcleaner "github.com/dohernandez/go-grpc-service/test/feature/database"
	"github.com/dohernandez/goservicing"
	"github.com/nhatthm/clockdog"
)

type FeaturesConfig struct {
	FeaturePath        string
	Locator            *app.Locator
	FeatureContextFunc func(t *testing.T, s *godog.ScenarioContext)
	Tables             map[string]any
}

func RunFeatures(t *testing.T, ctx context.Context, cfg FeaturesConfig) {
	t.Helper()

	deps := cfg.Locator

	clock := clockdog.New()

	deps.ClockProvider = clock

	dbm := initDBManager(deps.Storage, cfg.Tables)
	dbmCleaner := initDBMCleaner(dbm)

	services := goservicing.WithGracefulShutDown(
		func(ctx context.Context) {
			deps.Close() //nolint:errcheck
		},
	)

	go func() {
		err := services.Start(
			ctx,
			time.Second*5,
			func(ctx context.Context, msg string) {
				deps.CtxdLogger().Important(ctx, msg)
			},
			deps.GRPCService,
			deps.GRPCRestService,
		)
		must.NotFail(ctxd.WrapError(ctx, err, "failed to start the services"))
	}()

	baseRESTURL := <-deps.GRPCRestService.AddrAssigned
	local := httpdog.NewLocal(baseRESTURL)

	feature.RunFeatures(t, cfg.FeaturePath, func(_ *testing.T, s *godog.ScenarioContext) {
		local.RegisterSteps(s)

		dbm.RegisterSteps(s)
		dbmCleaner.RegisterSteps(s)

		clock.RegisterContext(s)

		cfg.FeatureContextFunc(t, s)
	})

	must.NotFail(services.Close())
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
