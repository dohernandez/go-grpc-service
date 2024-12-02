package app

import (
	"context"
	"fmt"
	"github.com/bool64/ctxd"
	"github.com/bool64/sqluct"
	"github.com/bool64/zapctxd"
	"github.com/dohernandez/go-grpc-service/config"
	"github.com/dohernandez/go-grpc-service/logger"
	"github.com/dohernandez/servers"
	grpcLogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jmoiron/sqlx"
	clockSrv "github.com/nhatthm/go-clock/service"
)

const pgxDriver = "pgx"

// Option sets up service locator.
type Option func(l *Locator)

func WithPostgresDBx() Option {
	return func(l *Locator) {
		l.opts.postgresEnabled = true
	}
}

// WithGRPC sets up gRPC and with server options.
func WithGRPC(opts ...servers.Option) Option {
	return func(l *Locator) {
		l.opts.grpcEnabled = true

		if len(opts) == 0 {
			return
		}

		l.opts.grpcOpts = append(l.opts.grpcOpts, opts...)
	}
}

// WithGRPCRest sets up gRPC REST and with server options.
func WithGRPCRest(opts ...servers.Option) Option {
	return func(l *Locator) {
		l.opts.grpcRestEnabled = true

		if len(opts) == 0 {
			return
		}

		l.opts.grpcRestOpts = append(l.opts.grpcRestOpts, opts...)
	}
}

// WithMetrics sets up metrics and with server options.
func WithMetrics(opts ...servers.Option) Option {
	return func(l *Locator) {
		l.opts.metricsEnabled = true

		if len(opts) == 0 {
			return
		}

		l.opts.metricsOpts = append(l.opts.metricsOpts, opts...)
	}
}

type locatorOptions struct {
	postgresEnabled bool

	grpcEnabled bool
	grpcOpts    []servers.Option

	grpcRestEnabled bool
	grpcRestOpts    []servers.Option

	metricsEnabled bool
	metricsOpts    []servers.Option
}

// Locator defines application resources.
type Locator struct {
	Config *config.Config

	opts locatorOptions

	DBx     *sqlx.DB
	Storage *sqluct.Storage

	logger *zapctxd.Logger
	ctxd.LoggerProvider

	clockSrv.ClockProvider

	GRPCService     *servers.GRPC
	GRPCRestService *servers.GRPCRest
	MetricsService  *servers.Metrics
}

// NewServiceLocator creates application locator.
func NewServiceLocator(cfg *config.Config, opts ...Option) (*Locator, error) {
	l := Locator{
		Config: cfg,
	}

	for _, o := range opts {
		o(&l)
	}

	var err error

	l.logger = logger.InitLogger(l.Config.Log, l.Config.IsDev())
	l.LoggerProvider = l.logger

	// init db
	if l.opts.postgresEnabled {
		l.DBx, err = sqlx.Connect(pgxDriver, cfg.DB.DSN)
		if err != nil {
			return nil, fmt.Errorf("connect to postgres: %w", err)
		}

		l.Storage = sqluct.NewStorage(l.DBx)
	}

	l.opts.grpcOpts = append(
		l.opts.grpcOpts,
		servers.WithChainUnaryInterceptor(
			// recovering from panic
			grpcRecovery.UnaryServerInterceptor(),
			grpcLogging.UnaryServerInterceptor(grpcInterceptorLogger(l.logger)),
		),
	)

	if l.opts.metricsEnabled {
		l.opts.grpcOpts = append(
			l.opts.grpcOpts,
			// metrics
			servers.WithChainUnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		)

		grpc_prometheus.EnableHandlingTimeHistogram()
	}

	return &l, nil
}

func (l *Locator) InitGRPCService(opts ...servers.Option) {
	if !l.opts.grpcEnabled {
		return
	}

	grpcOpts := append(
		[]servers.Option{},
		l.opts.grpcOpts...,
	)

	grpcOpts = append(
		grpcOpts,
		opts...,
	)

	l.GRPCService = servers.NewGRPC(
		servers.Config{
			Name: "grpc " + l.Config.ServiceName,
		},
		grpcOpts...,
	)
}

func (l *Locator) InitGRPCRestService(opts ...servers.Option) error {
	var err error

	if !l.opts.grpcRestEnabled {
		return nil
	}

	grpcRestOpts := append(
		[]servers.Option{},
		l.opts.grpcRestOpts...,
	)

	grpcRestOpts = append(
		grpcRestOpts,
		opts...,
	)

	l.GRPCRestService, err = servers.NewGRPCRest(
		servers.Config{
			Name: "grpc rest " + l.Config.ServiceName,
		},
		grpcRestOpts...,
	)
	if err != nil {
		return fmt.Errorf("creating grpc rest service: %w", err)
	}

	return nil
}

func (l *Locator) InitMetricsService(opts ...servers.Option) {
	// Check if metrics service is enabled.
	if !l.opts.metricsEnabled {
		return
	}

	metricsOpts := append(
		[]servers.Option{},
		l.opts.metricsOpts...,
	)

	metricsOpts = append(
		metricsOpts,
		opts...,
	)

	l.MetricsService = servers.NewMetrics(
		servers.Config{
			Name: "metrics " + l.Config.ServiceName,
		},
		metricsOpts...,
	)
}

// grpcInterceptorLogger adapts zapctxd logger to interceptor logger.
func grpcInterceptorLogger(l *zapctxd.Logger) grpcLogging.Logger {
	return grpcLogging.LoggerFunc(func(ctx context.Context, lvl grpcLogging.Level, msg string, _ ...any) {
		switch lvl {
		case grpcLogging.LevelDebug:
			l.Debug(ctx, msg)
		case grpcLogging.LevelInfo:
			l.Info(ctx, msg)
		case grpcLogging.LevelWarn:
			l.Warn(ctx, msg)
		case grpcLogging.LevelError:
			l.Error(ctx, msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func (l *Locator) Close() error {
	if l.DBx != nil {
		return l.DBx.Close()
	}

	return nil
}
