GO ?= go

-include $(DEVSERVICEGO_PATH)/dev/makefiles/docker.mk

# Override in app Makefile to add custom ldflags, example BUILD_LDFLAGS="-s -w"
BUILD_LDFLAGS ?= ""
INTEGRATION_TEST_TARGET ?= -coverpkg ./internal/... integration_test.go
INTEGRATION_DOCKER_COMPOSE ?= ./docker-compose.yml
INTEGRATION_DOCKER_COMPOSE_PROFILE ?= integration-test

## Run integration tests
test-integration:
	@DOCKER_COMPOSE_PROFILE=$(INTEGRATION_DOCKER_COMPOSE_PROFILE) \
	DOCKER_COMPOSE_PATH=$(INTEGRATION_DOCKER_COMPOSE) \
	make dc-up
	@echo "Running integration tests."
	@CGO_ENABLED=1 $(GO) test -ldflags "$(shell bash $(DEVGO_SCRIPTS)/version-ldflags.sh && echo $(BUILD_LDFLAGS))" -race -cover -coverprofile ./integration.coverprofile $(INTEGRATION_TEST_TARGET)


.PHONY: test-integration start-deps stop-deps
