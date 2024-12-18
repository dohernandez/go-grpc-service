GO ?= go
UNIT_TEST_COUNT ?= 2

# Override in app Makefile to exclude specific package from the test and coverage, example EXCLUDE_PACKAGES=/mocks|/pb|/cmd
EXCLUDE_PACKAGES ?= ""

TEST_PACKAGES := $(shell $(GO) list ./... | grep -vE '$(EXCLUDE_PACKAGES)')

## Run unit tests
test-unit:
	@echo "Running unit tests."
	@CGO_ENABLED=1 $(GO) test -short -coverprofile=unit.coverprofile -covermode=atomic -race $(TEST_PACKAGES)

## Run unit tests multiple times, use `UNIT_TEST_COUNT=10 make test-unit-multi` to control count
test-unit-multi:
	@echo "Running unit tests ${UNIT_TEST_COUNT} times."
	@CGO_ENABLED=1 $(GO) test -short -coverprofile=unit.coverprofile -count $(UNIT_TEST_COUNT) -covermode=atomic -race $(TEST_PACKAGES)

.PHONY: test-unit test-unit-multi
