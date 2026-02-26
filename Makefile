BUILD_DIR = bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

.PHONY: all build install uninstall clean test lint

all: build

build:
	@echo "Building spretty..."
	go build $(LDFLAGS) -o $(BUILD_DIR)/spretty ./cmd/spretty

install:
	@echo "Installing spretty..."
	@bin_dir=$$(go env GOBIN); \
	if [ -z "$$bin_dir" ]; then \
		bin_dir=$$(go env GOPATH)/bin; \
	fi; \
	mkdir -p "$$bin_dir"; \
	echo "Installing to $$bin_dir"; \
	go build $(LDFLAGS) -o "$$bin_dir/spretty" ./cmd/spretty

uninstall:
	@echo "Uninstalling spretty..."
	@bin_dir=$$(go env GOBIN); \
	if [ -z "$$bin_dir" ]; then \
		bin_dir=$$(go env GOPATH)/bin; \
	fi; \
	rm -f "$$bin_dir/spretty"

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

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
	cd ./example && golangci-lint run ./...
