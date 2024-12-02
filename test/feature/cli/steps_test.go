package cli_test

import (
	"bytes"
	"github.com/dohernandez/go-grpc-service/test/feature/cli"
	"github.com/dohernandez/go-grpc-service/test/feature/cli/testdata"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

func TestNewContext(t *testing.T) {
	buf := bytes.NewBuffer(nil)

	app := cli.App{}
	app.Add("greet", testdata.NewApp)

	suite := godog.TestSuite{
		Name:                 "cliSteps",
		TestSuiteInitializer: nil,
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			cli.RegisterContext(s, &app)
		},
		Options: &godog.Options{
			Format:    "pretty",
			Output:    buf,
			Paths:     []string{"testdata/Test.feature"},
			Strict:    true,
			Randomize: time.Now().UTC().UnixNano(),
		},
	}

	status := suite.Run()
	if status != 0 {
		t.Fatal(buf.String())
	}
}
