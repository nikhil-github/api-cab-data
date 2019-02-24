SHELL                :=/bin/bash
REPO_NAME            := api-cab-data
GIT_HASH             := $$(git rev-parse --short HEAD)
REPO_VERSION         := 1.0
GOBUILD_ARGS         := -ldflags "-s -X main.Version=$(REPO_VERSION) -X main.gitCommit=$(GIT_HASH)"
BINARY               := $(REPO_NAME)


build-all: depend fmt test build

depend:
	dep ensure --vendor-only -v
	go get golang.org/x/tools/cmd/goimports

build:
	go build $(GOBUILD_ARGS) ./cmd/$(BINARY)

fmt:
	gofmt -w -s $$(find . -type f -name '*.go' -not -path "./vendor/*")
	goimports -w -local github.com/nikhil-github/ -d $$(find . -type f -name '*.go' -not -path "./vendor/*")

test:
	go test -v -race ./...

bench:
	go test -bench=. ./...

run: build-all
	./$(BINARY)

build-docker:
	docker-compose build

run-docker: build-docker
	docker-compose up -d && docker-compose logs -f

stop-docker:
	docker-compose stop

start-db:
	docker-compose up -d db

.PHONY: bench build build-all depend fmt run run-docker stop-docker test
