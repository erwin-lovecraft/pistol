PROJECT_NAME := pistol
TEAM_NAME := viebiz
ENV := dev

# Shorten cmd
DOCKER_BUILD_BIN := docker
COMPOSE_BIN := PROJECT_NAME=$(PROJECT_NAME) TEAM_NAME=$(TEAM_NAME) ENV=$(ENV) docker compose --file build/docker-compose.yaml --project-directory . -p $(PROJECT_NAME)
COMPOSE_TOOL_RUN := $(COMPOSE_BIN) run --rm --service-ports tool

# Run cmd
.PHONY: run dev debug
run:
	@$(COMPOSE_TOOL_RUN) sh -c "go run ./cmd/serverd"

# Setup cmd
.PHONY: build-dev-image build-prod-image pg
build-dev-image:
	@$(DOCKER_BUILD_BIN) build -f build/local.go.Dockerfile -t ${PROJECT_NAME}-local:latest .
	-docker images -q -f "dangling=true" | xargs docker rmi -f

build-server-image:
	@$(DOCKER_BUILD_BIN) build -f build/server.Dockerfile -t ${PROJECT_NAME}-prod:latest .

pg:
	@$(COMPOSE_BIN) up -d pg

# Helper cmd
.PHONY: buf-lint buf-gen
buf-lint:
	@$(DOCKER_BUILD_BIN) run --volume ".:/workspace" --workdir /workspace bufbuild/buf lint

buf-gen:
	@$(DOCKER_BUILD_BIN) run --volume ".:/workspace" --workdir /workspace bufbuild/buf generate

teardown:
	@$(COMPOSE_BIN) down
