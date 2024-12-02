package go_grpc_service

import (
	"context"
	"time"

	"github.com/bool64/ctxd"
	"github.com/dohernandez/go-grpc-service/app"
	"github.com/dohernandez/go-grpc-service/must"
	"github.com/dohernandez/goservicing"
)

func Run(ctx context.Context, l *app.Locator) {
	services := goservicing.WithGracefulShutDown(
		func(ctx context.Context) {
			l.Close() //nolint:errcheck
		},
	)

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

	err := services.Start(
		ctx,
		time.Second*5,
		func(ctx context.Context, msg string) {
			l.CtxdLogger().Important(ctx, msg)
		},
		srvs...,
	)
	must.NotFail(ctxd.WrapError(ctx, err, "failed to start the services"))
}
