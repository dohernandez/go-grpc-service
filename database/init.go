package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/bool64/ctxd"
	"github.com/bool64/sqluct"
	"github.com/cenkalti/backoff/v4"
	"github.com/jmoiron/sqlx"
	"github.com/opencensus-integrations/ocsql"
)

const pgxDriver = "pgx"

// initDBx initializes database.
func initDBx(cfg Config) (*sqlx.DB, error) {
	db, err := initDB(cfg)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("ping context: %w", err)
	}

	return sqlx.NewDb(db, cfg.DriverName), nil
}

// initDB initializes database.
func initDB(cfg Config) (*sql.DB, error) {
	driverName, err := ocsql.Register(cfg.DriverName,
		ocsql.WithQuery(true),
		ocsql.WithRowsClose(true),
		ocsql.WithRowsAffected(true),
		ocsql.WithAllowRoot(true),
	)
	if err != nil {
		return nil, err
	}

	ocsql.RegisterAllViews()

	db, err := sql.Open(driverName, cfg.DSN)
	if err != nil {
		return nil, err
	}

	ocsql.RecordStats(db, time.Second)

	db.SetConnMaxLifetime(cfg.MaxLifetime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

// ConnectPostgres initializes and connects to postgres database storage.
func ConnectPostgres(
	cfg Config,
	logger ctxd.Logger,
) (*sqluct.Storage, error) {
	var (
		dBx *sqlx.DB
		err error
	)

	cfg.DriverName = pgxDriver

	fn := func() error {
		dBx, err = initDBx(cfg)
		if err != nil {
			return fmt.Errorf("connect to postgres: %w", err)
		}

		return nil
	}

	if err = backoff.Retry(fn, backoff.WithMaxRetries(backoff.NewConstantBackOff(5*time.Second), 3)); err != nil {
		return nil, err
	}

	st := sqluct.NewStorage(dBx)

	st.Format = squirrel.Dollar
	st.OnError = func(ctx context.Context, err error) {
		logger.Error(ctx, "storage failure", "error", err)
	}

	return st, nil
}
