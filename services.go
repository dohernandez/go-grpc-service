package go_grpc_service

import (
	"context"
	"sync"
	"time"

	"github.com/dohernandez/go-grpc-service/app"
	"github.com/dohernandez/goservicing"
)

const timeout = time.Second * 5

func RunServices(ctx context.Context, l *app.Locator) error {
	srvs := make([]goservicing.Service, 0, 3)

	if l.GRPCService != nil {
		srvs = append(srvs, l.GRPCService)
	}

	if l.GRPCRestService != nil {
		srvs = append(srvs, l.GRPCRestService)
	}

	if l.MetricsService != nil {
		srvs = append(srvs, l.MetricsService)
	}

	if l.HealthService != nil {
		srvs = append(srvs, l.HealthService)
	}

	if len(srvs) == 0 {
		panic("no services to start")
	}

	services := goservicing.WithGracefulShutDown(
		func(ctx context.Context) {
			_ = l.Close() //nolint:errcheck
		},
	)

	return services.Start(
		ctx,
		time.Second*5,
		func(ctx context.Context, msg string) {
			l.CtxdLogger().Important(ctx, msg)
		},
		srvs...,
	)
}

type ServicesT interface {
	Helper()
	Fatalf(format string, args ...any)
	Log(args ...any)
}

func RunServicesTesting(t ServicesT, ctx context.Context, l *app.Locator) (func(), chan error) {
	t.Helper()

	var (
		errch = make(chan error, 1)
		done  bool
		sm    sync.Mutex
	)

	// Check if the service is done, to avoid failing during context cancellation when the feature finished due to it is
	// the signal use to stop the services.
	isDone := func() bool {
		sm.Lock()
		defer sm.Unlock()

		return done
	}

	ctx, cancel := context.WithCancel(ctx)

	stopFunc := func() {
		sm.Lock()
		done = true
		sm.Unlock()

		cancel()

		// Wait for the service to stop.
		select {
		case err := <-errch:
			if err != nil && !isDone() {
				t.Fatalf("failed to stop services: %v", err)
			}
		case <-time.After(timeout):
			t.Log("timeout waiting for services to stop")
		}

		close(errch)
	}

	// Run the service.
	go func() {
		errch <- RunServices(ctx, l)
	}()

	return stopFunc, errch
}
