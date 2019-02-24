SHELL:=/bin/bash
REPO_NAME := api-cab-data
VERSION_VAR := main.version
GIT_COMMIT_VAR := main.gitCommit
REPO_VERSION := $$(git describe --tags 2>/dev/null || echo "nil")
GIT_HASH := $$(git rev-parse --short HEAD)
GOBUILD_VERSION_ARGS := -ldflags "-s -X $(VERSION_VAR)=$(REPO_VERSION) -X $(GIT_COMMIT_VAR)=$(GIT_HASH)"
DOCKER_IMAGE := 175914186171.dkr.ecr.ap-southeast-2.amazonaws.com/$(REPO_NAME)
BINARY := $(REPO_NAME)
BASE_PKG := bitbucket.org/ffxblue/$(REPO_NAME)
MAIN_PKG := ${BASE_PKG}/cmd/$(REPO_NAME)
CLIENT_PKG := ${BASE_PKG}/cmd/$(REPO_NAME)-client
SPLIT_LOGS := tee >(grep --line-buffered -E '^{' | jq 1>&2) | grep -Ev '^{'

# Default target (since it's the first without '.' prefix)
build-all: depend fmt build

# Docker build will pull dependencies as a separate step
build-all-docker: generate check cover build

# CI will pull dependencies as a separate step, and should not apply formatting (that should be committed correctly)
build-all-ci: generate check cover build

depend:
	dep ensure --vendor-only -v
	go get golang.org/x/tools/cmd/goimports

build:
	go build $(GOBUILD_VERSION_ARGS) ./cmd/$(BINARY)

fmt:
	gofmt -w -s $$(find . -type f -name '*.go' -not -path "./vendor/*")

test:
	go test ./...

bench:
	go test -bench=. ./...

run: build
	./$(BINARY)

docker:
	docker build --tag "${DOCKER_IMAGE}" .

run-docker: docker
	docker-compose up -d

stop-docker:
	docker-compose stop

# None of the Make tasks generate files with the name of the task, so all must be declared as 'PHONY'
.PHONY: bench bench-race build build-all build-all-ci build-ci build-client build-docker build-race check check-go check-proto cover cover-report depend depend-ci fmt generate generate-proto run run-docker stop-docker test test-contract test-race watch
