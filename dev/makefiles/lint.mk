GO ?= go

## Check with golangci-lint
lint:
	@GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION) bash $(DEVSERVICEGO_SCRIPTS)/lint.sh

## Apply goimports and gofmt
fix-lint:
	@bash $(DEVSERVICEGO_SCRIPTS)/fix.sh

.PHONY: lint fix-lint
