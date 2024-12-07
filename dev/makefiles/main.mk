# MODULES is a list of dev modules (mk) to be included in the project.
ifndef MODULES
	MODULES := \
		DEVGO_PATH=github.com/bool64/dev \
        DEVGRPCGO_PATH=github.com/dohernandez/dev-grpc
else
	MODULES := \
		$(MODULES) \
		DEVGO_PATH=github.com/bool64/dev \
        DEVGRPCGO_PATH=github.com/dohernandez/dev-grpc
endif

# The head of Makefile determines location of dev-go to include standard targets.
GO ?= go
export GO111MODULE = on

PWD = $(shell pwd)

ifneq "$(GOFLAGS)" ""
  $(info GOFLAGS: ${GOFLAGS})
endif

# Use vendored dependencies if available.
ifneq ($(wildcard ./vendor),)
  modVendor := -mod=vendor
  ifeq (,$(findstring -mod,$(GOFLAGS)))
      export GOFLAGS := ${GOFLAGS} ${modVendor}
  endif
endif

# Set dev module paths or download them.
$(foreach module,$(MODULES), \
	$(eval key=$(word 1,$(subst =, ,$(module)))); \
	$(eval value=$(word 2,$(subst =, ,$(module)))); \
	\
	$(if $(wildcard ./vendor/$(value)), \
		$(eval export $(key)=./vendor/$(value)); \
	) \
	\
	$(if $(strip $($(key))), , \
    	$(eval export $(key)=$(shell GO111MODULE=on $(GO) list ${modVendor} -f '{{.Dir}}' -m $(value))); \
		$(if $(strip $($(key))), \
			$(info Module $(value) not found, downloading.); \
			$(eval export $(key)=$(shell export GO111MODULE=on && $(GO) get $(value) && $(GO) list -f '{{.Dir}}' -m $(value))); \
		) \
    ) \
)

# Detecting GOPATH and removing trailing "/" if any
GOPATH = $(realpath $(shell $(GO) env GOPATH))

ifneq "$(wildcard ./vendor )" ""
  modVendor = -mod=vendor
endif
export MODULE_NAME := $(shell test -f go.mod && GO111MODULE=on $(GO) list $(modVendor) -m)

# If DEVSERVICEGO_PATH is not exported, it is because we are in the root project therefore PWD.
DEVSERVICEGO_PATH ?= $(PWD)

DEVSERVICEGO_PATH := $(DEVSERVICEGO_PATH)/dev
DEVSERVICEGO_SCRIPTS ?= $(DEVSERVICEGO_PATH)/scripts

-include $(DEVGO_PATH)/makefiles/main.mk