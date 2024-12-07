#MOCKERY_VERSION="2.36.0" # Optional configuration to pinpoint mockery version.

GO ?= go

## Check/install mockery tool
mockery-cli:
	@MOCKERY_VERSION=$(MOCKERY_VERSION) bash $(DEVSERVICEGO_SCRIPTS)/mockery-cli.sh

.PHONY: mockery-cli
