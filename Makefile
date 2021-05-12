.PHONY: all
all: lint build test

.PHONY: build
build:
	curl --silent https://www.unicode.org/Public/emoji/13.1/emoji-test.txt --output emoji-test.txt
# Remove comments and blank lines to save space (requires gnu-sed)
	sed '/^[ \t]*#/d' -i emoji-test.txt
	sed '/^$$/d' -i emoji-test.txt
	go build -v ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: lint
lint:
	golangci-lint run ./...
