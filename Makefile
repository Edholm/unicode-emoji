.PHONY: all
all: lint build test

.PHONY: build
build:
	curl -q https://www.unicode.org/Public/emoji/13.1/emoji-test.txt > emoji-test.txt
# Remove comments and blank lines to save space
	sed -i -e '/^[ \t]*#/d' emoji-test.txt
	sed -i -e '/^$$/d' emoji-test.txt
	go build ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: lint
lint:
	golangci-lint run ./...
