include .env

BINARY_NAME=weatherboyapi

.PHONY: run
run: build
	@set -a && source .env && set +a && ./tmp/$(BINARY_NAME)

.PHONY: build
build:
	go build -o tmp/$(BINARY_NAME) ./cmd/weatherboyapi

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	go clean
	rm -f tmp/$(BINARY_NAME)
