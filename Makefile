# gql2j Makefile

# Variables
BINARY_NAME := gql2j
MAIN_PATH := ./cmd/gql2j
BUILD_DIR := ./build
VERSION := 1.0.0
GO := go
GOFLAGS := -v
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

# Java integration test variables
JAVA_TEST_DIR := ./build/java-test
JAVA_LIBS_DIR := ./build/libs
LOMBOK_VERSION := 1.18.34
JAKARTA_VALIDATION_VERSION := 3.0.2
# Lombok supports up to Java 23 (as of 1.18.34). Java 24+ may not work.
LOMBOK_MAX_JAVA_VERSION := 23
LOMBOK_JAR := $(JAVA_LIBS_DIR)/lombok-$(LOMBOK_VERSION).jar
JAKARTA_VALIDATION_JAR := $(JAVA_LIBS_DIR)/jakarta.validation-api-$(JAKARTA_VALIDATION_VERSION).jar

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
	rm -rf $(JAVA_TEST_DIR)

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

# ============================================================================
# Java Integration Tests
# ============================================================================
# These targets generate Java code and verify it compiles correctly.
# This provides end-to-end validation that gql2j produces valid Java code.

# Check if Java is available
.PHONY: check-java
check-java:
	@which javac > /dev/null || (echo "Error: javac not found. Please install JDK." && exit 1)
	@echo "Java compiler found: $$(javac -version 2>&1)"

# Get Java major version number
JAVA_VERSION := $(shell javac -version 2>&1 | sed -n 's/javac \([0-9]*\).*/\1/p')

# Check if Lombok is compatible with current Java version
.PHONY: check-lombok-compat
check-lombok-compat:
	@if [ "$(JAVA_VERSION)" -gt "$(LOMBOK_MAX_JAVA_VERSION)" ]; then \
		echo "WARNING: Lombok $(LOMBOK_VERSION) may not support Java $(JAVA_VERSION)."; \
		echo "         Lombok currently supports up to Java $(LOMBOK_MAX_JAVA_VERSION)."; \
		echo "         See: https://projectlombok.org/changelog"; \
		echo "         Skipping Lombok test."; \
		exit 1; \
	fi

# Download Lombok jar for annotation processing
$(LOMBOK_JAR):
	@mkdir -p $(JAVA_LIBS_DIR)
	@echo "Downloading Lombok $(LOMBOK_VERSION)..."
	@curl -fsSL -o $(LOMBOK_JAR) \
		"https://repo1.maven.org/maven2/org/projectlombok/lombok/$(LOMBOK_VERSION)/lombok-$(LOMBOK_VERSION).jar"
	@echo "Downloaded: $(LOMBOK_JAR)"

# Download Jakarta Validation API for validation annotations
$(JAKARTA_VALIDATION_JAR):
	@mkdir -p $(JAVA_LIBS_DIR)
	@echo "Downloading Jakarta Validation API $(JAKARTA_VALIDATION_VERSION)..."
	@curl -fsSL -o $(JAKARTA_VALIDATION_JAR) \
		"https://repo1.maven.org/maven2/jakarta/validation/jakarta.validation-api/$(JAKARTA_VALIDATION_VERSION)/jakarta.validation-api-$(JAKARTA_VALIDATION_VERSION).jar"
	@echo "Downloaded: $(JAKARTA_VALIDATION_JAR)"

# Download all Java dependencies
.PHONY: java-deps
java-deps: $(LOMBOK_JAR) $(JAKARTA_VALIDATION_JAR)
	@echo "All Java dependencies downloaded."

# Basic Java compilation test (no external dependencies)
# Generates plain Java classes and compiles them with javac
.PHONY: test-java
test-java: build check-java
	@echo "=== Java Integration Test: Basic ==="
	@rm -rf $(JAVA_TEST_DIR)/basic
	@mkdir -p $(JAVA_TEST_DIR)/basic
	@echo "Generating Java code..."
	./$(BINARY_NAME) \
		-schema ./testdata/schemas/comprehensive.graphql \
		-output $(JAVA_TEST_DIR)/basic \
		-package com.example.generated \
		-verbose
	@echo ""
	@echo "Compiling generated Java code..."
	@find $(JAVA_TEST_DIR)/basic -name "*.java" -print | head -5
	@echo "..."
	javac -d $(JAVA_TEST_DIR)/basic/classes \
		$(JAVA_TEST_DIR)/basic/*.java
	@echo ""
	@echo "✓ Basic Java compilation test PASSED"
	@echo "  Generated: $$(find $(JAVA_TEST_DIR)/basic -name '*.java' | wc -l | tr -d ' ') Java files"
	@echo "  Compiled:  $$(find $(JAVA_TEST_DIR)/basic/classes -name '*.class' | wc -l | tr -d ' ') class files"

# Java compilation test with Lombok
# Requires lombok.jar for annotation processing
# Note: Lombok may not support the latest Java versions. Check https://projectlombok.org/changelog
.PHONY: test-java-lombok
test-java-lombok: build check-java $(LOMBOK_JAR)
	@echo "=== Java Integration Test: Lombok ==="
	@if [ "$(JAVA_VERSION)" -gt "$(LOMBOK_MAX_JAVA_VERSION)" ]; then \
		echo "⚠ SKIPPED: Lombok $(LOMBOK_VERSION) does not support Java $(JAVA_VERSION)"; \
		echo "  Lombok currently supports up to Java $(LOMBOK_MAX_JAVA_VERSION)."; \
		echo "  See: https://projectlombok.org/changelog"; \
	else \
		rm -rf $(JAVA_TEST_DIR)/lombok && \
		mkdir -p $(JAVA_TEST_DIR)/lombok && \
		echo "Generating Java code with Lombok..." && \
		./$(BINARY_NAME) \
			-schema ./testdata/schemas/comprehensive.graphql \
			-output $(JAVA_TEST_DIR)/lombok \
			-package com.example.generated \
			-lombok \
			-verbose && \
		echo "" && \
		echo "Compiling with Lombok annotation processor..." && \
		javac -d $(JAVA_TEST_DIR)/lombok/classes \
			-cp $(LOMBOK_JAR) \
			-processorpath $(LOMBOK_JAR) \
			$(JAVA_TEST_DIR)/lombok/*.java && \
		echo "" && \
		echo "✓ Lombok Java compilation test PASSED" && \
		echo "  Generated: $$(find $(JAVA_TEST_DIR)/lombok -name '*.java' | wc -l | tr -d ' ') Java files" && \
		echo "  Compiled:  $$(find $(JAVA_TEST_DIR)/lombok/classes -name '*.class' | wc -l | tr -d ' ') class files"; \
	fi

# Java compilation test with validation annotations
# Requires jakarta.validation-api.jar
.PHONY: test-java-validation
test-java-validation: build check-java $(JAKARTA_VALIDATION_JAR)
	@echo "=== Java Integration Test: Validation ==="
	@rm -rf $(JAVA_TEST_DIR)/validation
	@mkdir -p $(JAVA_TEST_DIR)/validation
	@echo "Generating Java code with validation annotations..."
	./$(BINARY_NAME) \
		-schema ./testdata/schemas/comprehensive.graphql \
		-output $(JAVA_TEST_DIR)/validation \
		-package com.example.generated \
		-validation \
		-verbose
	@echo ""
	@echo "Compiling with Jakarta Validation API..."
	javac -d $(JAVA_TEST_DIR)/validation/classes \
		-cp $(JAKARTA_VALIDATION_JAR) \
		$(JAVA_TEST_DIR)/validation/*.java
	@echo ""
	@echo "✓ Validation Java compilation test PASSED"
	@echo "  Generated: $$(find $(JAVA_TEST_DIR)/validation -name '*.java' | wc -l | tr -d ' ') Java files"
	@echo "  Compiled:  $$(find $(JAVA_TEST_DIR)/validation/classes -name '*.class' | wc -l | tr -d ' ') class files"

# Java compilation test with all features
.PHONY: test-java-full
test-java-full: build check-java $(LOMBOK_JAR) $(JAKARTA_VALIDATION_JAR)
	@echo "=== Java Integration Test: Full (Lombok + Validation) ==="
	@if [ "$(JAVA_VERSION)" -gt "$(LOMBOK_MAX_JAVA_VERSION)" ]; then \
		echo "⚠ SKIPPED: Lombok $(LOMBOK_VERSION) does not support Java $(JAVA_VERSION)"; \
		echo "  Lombok currently supports up to Java $(LOMBOK_MAX_JAVA_VERSION)."; \
		echo "  See: https://projectlombok.org/changelog"; \
	else \
		rm -rf $(JAVA_TEST_DIR)/full && \
		mkdir -p $(JAVA_TEST_DIR)/full && \
		echo "Generating Java code with all features..." && \
		./$(BINARY_NAME) \
			-schema ./testdata/schemas/comprehensive.graphql \
			-output $(JAVA_TEST_DIR)/full \
			-package com.example.generated \
			-lombok \
			-validation \
			-verbose && \
		echo "" && \
		echo "Compiling with Lombok and Jakarta Validation..." && \
		javac -d $(JAVA_TEST_DIR)/full/classes \
			-cp "$(LOMBOK_JAR):$(JAKARTA_VALIDATION_JAR)" \
			-processorpath $(LOMBOK_JAR) \
			$(JAVA_TEST_DIR)/full/*.java && \
		echo "" && \
		echo "✓ Full Java compilation test PASSED" && \
		echo "  Generated: $$(find $(JAVA_TEST_DIR)/full -name '*.java' | wc -l | tr -d ' ') Java files" && \
		echo "  Compiled:  $$(find $(JAVA_TEST_DIR)/full/classes -name '*.class' | wc -l | tr -d ' ') class files"; \
	fi

# Run all Java integration tests
.PHONY: test-java-all
test-java-all: test-java test-java-lombok test-java-validation test-java-full
	@echo ""
	@echo "============================================"
	@echo "✓ All Java integration tests completed"
	@if [ "$(JAVA_VERSION)" -gt "$(LOMBOK_MAX_JAVA_VERSION)" ]; then \
		echo "  Note: Lombok tests were skipped (Java $(JAVA_VERSION) > max supported $(LOMBOK_MAX_JAVA_VERSION))"; \
	fi
	@echo "============================================"

# Quick test - just basic (no downloads needed)
.PHONY: test-java-quick
test-java-quick: test-java
	@echo "Quick Java test completed (use 'make test-java-all' for full tests)"

# Clean Java test artifacts
.PHONY: clean-java-test
clean-java-test:
	rm -rf $(JAVA_TEST_DIR)

# Clean downloaded Java libraries
.PHONY: clean-java-libs
clean-java-libs:
	rm -rf $(JAVA_LIBS_DIR)

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
	@echo "  test           Run Go unit tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  lint           Run golangci-lint"
	@echo ""
	@echo "Java Integration Tests (compile generated code with javac):"
	@echo "  test-java            Basic test - plain Java classes (no dependencies)"
	@echo "  test-java-lombok     Test with Lombok annotations (requires Java ≤$(LOMBOK_MAX_JAVA_VERSION))"
	@echo "  test-java-validation Test with JSR-380 validation (downloads jakarta-validation-api.jar)"
	@echo "  test-java-full       Test with all features (Lombok + Validation, Java ≤$(LOMBOK_MAX_JAVA_VERSION))"
	@echo "  test-java-all        Run all Java integration tests"
	@echo "  test-java-quick      Alias for test-java (quick, no downloads)"
	@echo "  java-deps            Download Java dependencies only"
	@echo ""
	@echo "  Note: Lombok tests auto-skip if Java version > $(LOMBOK_MAX_JAVA_VERSION)"
	@echo ""
	@echo "Development targets:"
	@echo "  fmt            Format code"
	@echo "  tidy           Tidy go.mod"
	@echo "  deps           Download Go dependencies"
	@echo "  verify         Verify dependencies"
	@echo ""
	@echo "Run targets:"
	@echo "  run            Build and run with example schema"
	@echo "  run-lombok     Build and run with Lombok/validation enabled"
	@echo "  example        Generate example output files"
	@echo ""
	@echo "Clean targets:"
	@echo "  clean            Remove all build artifacts"
	@echo "  clean-output     Remove generated output directory"
	@echo "  clean-java-test  Remove Java test output"
	@echo "  clean-java-libs  Remove downloaded Java libraries"
	@echo ""
	@echo "Other:"
	@echo "  help           Show this help message"
