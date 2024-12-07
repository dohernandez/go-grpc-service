package go_grpc_service

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/bool64/ctxd"
	"github.com/bool64/httptestbench"
	"github.com/bool64/sqluct"
	"github.com/dohernandez/go-grpc-service/app"
	"github.com/dohernandez/go-grpc-service/must"
	"github.com/nhatthm/clockdog"
	"github.com/valyala/fasthttp"
)

type BenchmarkCases struct {
	Name string
	Uri  string
	Data map[string]any
}

type BenchmarkConfig struct {
	Locator   *app.Locator
	TestCases []BenchmarkCases
}

// RunBenchmark for performance benchmarking.
func RunBenchmark(b *testing.B, ctx context.Context, cfg *BenchmarkConfig) {
	deps := cfg.Locator

	clock := clockdog.New()
	deps.ClockProvider = clock

	stop, errch := RunServicesTesting(b, ctx, deps)

	var baseRESTURL string

	select {
	case err := <-errch:
		if err != nil {
			b.Fatalf("failed to run service: %v", err)
		}
	case baseRESTURL = <-deps.GRPCRestService.AddrAssigned:
		break
	case <-time.After(timeout):
		b.Fatal("timeout waiting for service to start")
	}

	defer func() {
		// Stop the service.
		stop()
	}()

	baseRESTURL = strings.Replace(baseRESTURL, "[::]", "127.0.0.1", 1)

	for _, tt := range cfg.TestCases {
		requestURI := "http://" + baseRESTURL + tt.Uri

		tables := make([]string, 0, len(tt.Data))

		for table, _ := range tt.Data {
			tables = append(tables, table)
		}

		cleanDatabase(ctx, deps.Storage, tables)
		loadDatabase(ctx, deps.Storage, tt.Data)

		b.Run(tt.Name, func(b *testing.B) {
			httptestbench.RoundTrip(b, 50,
				func(i int, req *fasthttp.Request) {
					req.SetRequestURI(requestURI)
				},
				func(i int, resp *fasthttp.Response) bool {
					return resp.StatusCode() == http.StatusOK
				},
			)
		})
	}
}

func cleanDatabase(ctx context.Context, s *sqluct.Storage, tables []string) {
	// Deleting from table
	for _, table := range tables {
		_, err := s.Exec(ctx, s.DeleteStmt(table))
		must.NotFail(ctxd.WrapError(ctx, err, "failed cleaning table", "table", table))
	}
}

func loadDatabase(ctx context.Context, s *sqluct.Storage, data map[string]any) {
	for table, d := range data {
		_, err := s.Exec(ctx, s.InsertStmt(table, d))
		must.NotFail(ctxd.WrapError(ctx, err, "failed loading", "table", table, "data", d))
	}
}
