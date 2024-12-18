package feature

import (
	"fmt"
	"github.com/bool64/httpmock"
	"github.com/cucumber/godog"
	"net/http"
)

// Health is step-driven HTTP client for application health check HTTP service.
type Health struct {
	*httpmock.Client
}

// NewHealth creates an instance of step-driven HTTP client.
func NewHealth(baseURL string) *Health {
	return &Health{
		Client: httpmock.NewClient(baseURL),
	}
}

func (h *Health) RegisterSteps(s *godog.ScenarioContext) {
	s.Step(`^I check server is up and running$`, h.iCheckServerIsUpAndRunning)
	s.Step(`^It should be up and running$`, h.itShouldBeUpAndRunning)

	s.Step(`^I check server is healthy$`, h.iCheckServerIsHealthy)
	s.Step(`^It should be healthy$`, h.itShouldBeHealthy)
}

func (h *Health) iCheckServerIsUpAndRunning() error {
	return h.iCheckSetup("/")
}

func (h *Health) iCheckSetup(method string) error {
	if err := h.CheckUnexpectedOtherResponses(); err != nil {
		return fmt.Errorf("unexpected other responses for previous request: %w", err)
	}

	h.Reset()
	h.WithMethod(http.MethodGet)
	h.WithURI(method)

	return nil
}

func (h *Health) itShouldBeUpAndRunning() error {
	return h.ExpectResponseStatus(200)
}

func (h *Health) iCheckServerIsHealthy() error {
	return h.iCheckSetup("/health")
}

func (h *Health) itShouldBeHealthy() error {
	return h.ExpectResponseStatus(200)
}
