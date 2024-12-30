# Variables
BINARY_NAME=pdu
BUILD_DIR=build
GO_FILES=$(shell find . -type f -name '*.go')

# Default target
.PHONY: all
all: build

# Build the project
.PHONY: build
build:
	@echo "Building project..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/pdu

# Run the project
.PHONY: run
run:
	@echo "Running project..."
	@go run ./cmd/pdu

# Run the samples
.PHONY: run-json
run-json:
	@echo "Running project..."
	@go run ./cmd/tools/json

# Test the project
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# Lint the code
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format the code
.PHONY: format
format:
	@echo "Formatting code..."
	@gofmt -s -w $(GO_FILES)

