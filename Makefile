.PHONY: all
all: lint build test

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: lint
lint:
	golangci-lint run ./...
