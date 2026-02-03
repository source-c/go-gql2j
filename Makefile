# gql2j Makefile

# Variables
BINARY_NAME := gql2j
MAIN_PATH := ./cmd/gql2j
BUILD_DIR := ./build
VERSION := 1.0.0
GO := go
GOFLAGS := -v
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

# Default target
.DEFAULT_GOAL := build

# Build the binary
.PHONY: build
build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

# Build with version info and optimizations
.PHONY: build-release
build-release:
	$(GO) build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

# Build for multiple platforms
.PHONY: build-all
build-all: clean-build
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Install to GOPATH/bin
.PHONY: install
install:
	$(GO) install $(LDFLAGS) $(MAIN_PATH)

# Run tests
.PHONY: test
test:
	$(GO) test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	$(GO) mod tidy

# Verify dependencies
.PHONY: verify
verify:
	$(GO) mod verify

# Download dependencies
.PHONY: deps
deps:
	$(GO) mod download

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -rf $(BUILD_DIR)

# Clean build directory only
.PHONY: clean-build
clean-build:
	rm -rf $(BUILD_DIR)

# Run the tool with example schema
.PHONY: run
run: build
	./$(BINARY_NAME) -schema example.graphql -output ./output -package com.example.model -verbose

# Run with Lombok enabled
.PHONY: run-lombok
run-lombok: build
	./$(BINARY_NAME) -schema example.graphql -output ./output -package com.example.model -lombok -validation -verbose

# Generate example output
.PHONY: example
example: build
	@mkdir -p ./output
	./$(BINARY_NAME) -schema example.graphql -output ./output -package com.example.model
	@echo "\nGenerated files:"
	@ls -la ./output/

# Clean generated output
.PHONY: clean-output
clean-output:
	rm -rf ./output

# Show help
.PHONY: help
help:
	@echo "gql2j Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build targets:"
	@echo "  build          Build the binary"
	@echo "  build-release  Build optimized binary with version info"
	@echo "  build-all      Build for all platforms (darwin, linux, windows)"
	@echo "  install        Install to GOPATH/bin"
	@echo ""
	@echo "Test targets:"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  lint           Run golangci-lint"
	@echo ""
	@echo "Development targets:"
	@echo "  fmt            Format code"
	@echo "  tidy           Tidy go.mod"
	@echo "  deps           Download dependencies"
	@echo "  verify         Verify dependencies"
	@echo ""
	@echo "Run targets:"
	@echo "  run            Build and run with example schema"
	@echo "  run-lombok     Build and run with Lombok/validation enabled"
	@echo "  example        Generate example output files"
	@echo ""
	@echo "Clean targets:"
	@echo "  clean          Remove all build artifacts"
	@echo "  clean-output   Remove generated output directory"
	@echo ""
	@echo "Other:"
	@echo "  help           Show this help message"
