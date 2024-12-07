
# Override in app Makefile to control docker file path.
DOCKERFILE_PATH ?= Dockerfile

# Override in app Makefile to control docker build context.
DOCKERBUILD_CONTEXT ?= .

# Override in app Makefile to control docker image tag.
DOCKER_IMAGE_TAG ?= latest

# Override in app Makefile to control docker image github token in case the docker required it into build.
DOCKER_GITHUB_TOKEN ?= ""

# Override in app Makefile to control docker docker-compose.yml path.
DOCKER_COMPOSE_PATH ?= ./docker-compose.yml

# Override in app Makefile to control docker docker-compose.yml project name.
DOCKER_COMPOSE_PROJECT_NAME ?= $(shell basename $$PWD)

# Override in app Makefile to control docker docker-compose.yml profile name.
DOCKER_COMPOSE_PROFILE ?= ""

# Override in app Makefile to control docker docker-compose.yml build using secret instead of args.
DOCKER_SECRET ?= false


DOCKER_COMPOSE_OPTIONS=$(foreach file,$(DOCKER_COMPOSE_PATH),-f $(file)) \
    -p $(DOCKER_COMPOSE_PROJECT_NAME) \
    $(if $(PROFILE),--profile $(PROFILE),$(if $(DOCKER_COMPOSE_PROFILE),--profile $(DOCKER_COMPOSE_PROFILE)))

## Build docker image
build-image:
	@DOCKER_IMAGE_TAG=$(DOCKER_IMAGE_TAG) \
	DOCKERFILE_PATH=$(DOCKERFILE_PATH) \
	DOCKERBUILD_CONTEXT=$(DOCKERBUILD_CONTEXT) \
	DOCKER_GITHUB_TOKEN=$(DOCKER_GITHUB_TOKEN) \
	DOCKER_SECRET=$(DOCKER_SECRET) \
	bash $(DEVSERVICEGO_SCRIPTS)/docker-build.sh

## Run docker-compose up from file DOCKER_COMPOSE_PATH with project name DOCKER_COMPOSE_PROJECT_NAME and profile DOCKER_COMPOSE_PROFILE.
## Usage: "make dc-up PROFILE=<profile>, if PROFILE is not provide, start only default services"
dc-up:
	@for file in $(DOCKER_COMPOSE_PATH); do \
	  if [ ! -f $$file ]; then \
	    echo "Error: File $$file not found!"; \
	    exit 1; \
	  fi; \
	done; \
	if command -v docker-compose >/dev/null 2>&1; then \
	  docker-compose $(DOCKER_COMPOSE_OPTIONS) up -d --remove-orphans; \
	else \
	  docker compose $(DOCKER_COMPOSE_OPTIONS) up -d --remove-orphans; \
	fi


## Run docker-compose down from file DOCKER_COMPOSE_PATH with project name DOCKER_COMPOSE_PROJECT_NAME
dc-down:
	@for file in $(DOCKER_COMPOSE_PATH); do \
	  if [ ! -f $$file ]; then \
		echo "Error: File $$file not found!"; \
		exit 1; \
	  fi; \
	done; \
	if command -v docker-compose >/dev/null 2>&1; then \
	  docker-compose $(DOCKER_COMPOSE_OPTIONS) down -v; \
	else \
	  docker compose $(DOCKER_COMPOSE_OPTIONS) down -v; \
	fi

## Run docker-compose logs from file DOCKER_COMPOSE_PATH with project name DOCKER_COMPOSE_PROJECT_NAME. Usage: "make generate APP=<docker-composer-service-name>"
dc-logs:
	@for file in $(DOCKER_COMPOSE_PATH); do \
	  if [ ! -f $$file ]; then \
		echo "Error: File $$file not found!"; \
		exit 1; \
	  fi; \
	done; \
	if command -v docker-compose >/dev/null 2>&1; then \
	  docker-compose $(DOCKER_COMPOSE_OPTIONS) logs $(APP); \
	else \
	  docker compose $(DOCKER_COMPOSE_OPTIONS) logs $(APP); \
	fi


.PHONY: build-image dc-up dc-down dc-logs