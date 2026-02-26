.PHONY: all test lint

all: build

build:
	@echo "Building..."
	go build ./...

test:
	@echo "Running tests..."
	go test -race ./...

lint:
	@echo "Running golangci-lint..."
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "golangci-lint is not installed"; \
		exit 1; \
	}
	golangci-lint run ./...
