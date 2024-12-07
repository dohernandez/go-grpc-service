#GOLANGCI_LINT_VERSION := "v1.61.0" # Optional configuration to pinpoint golangci-lint version.

PWD := $(shell pwd)

MODULES := \
    DEVGO_PATH=github.com/bool64/dev \
    DEVGRPCGO_PATH=github.com/dohernandez/dev-grpc

-include $(PWD)/dev/makefiles/main.mk

# Add your include here with based path to the module.

-include $(DEVGO_PATH)/makefiles/lint.mk
-include $(DEVGO_PATH)/makefiles/test-unit.mk

-include $(PWD)/dev/makefiles/docker.mk
-include $(PWD)/dev/makefiles/test-integration.mk
-include $(PWD)/dev/makefiles/mockery.mk

# Add your custom targets here.

PHONY: