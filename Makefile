EMOJI_TEST_FILE := emoji-test.txt
.PHONY: all
all: lint build test

.PHONY: build
build:
ifeq (,$(wildcard ${EMOJI_TEST_FILE}))
	curl --silent https://www.unicode.org/Public/emoji/13.1/emoji-test.txt --output ${EMOJI_TEST_FILE}
# Remove comments and blank lines to save space (requires gnu-sed)
	sed '/^[ \t]*#/d' -i ${EMOJI_TEST_FILE}
	sed '/^$$/d' -i ${EMOJI_TEST_FILE}
endif
	go build -v ./...

.PHONY: test
test:
	go test -v -race ./...

.PHONY: lint
lint:
	golangci-lint run ./...
