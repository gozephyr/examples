# Makefile for running cbreak examples

.PHONY: all clean
.PHONY: basic-all basic-run-simple
.PHONY: advanced-all advanced-run-failure-detection
.PHONY: integration-all integration-run-http-client

# Default target
all: basic-all advanced-all integration-all

# Clean target
clean:
	@echo "Cleaning up..."

# Basic examples
basic-all: basic-run-simple

basic-run-simple:
	@echo "Running simple circuit breaker example..."
	cd basic/simple && go run main.go

# Advanced examples
advanced-all: advanced-run-failure-detection

advanced-run-failure-detection:
	@echo "Running custom failure detection example..."
	cd advanced/failure_detection && go run main.go

# Integration examples
integration-all: integration-run-http-client

integration-run-http-client:
	@echo "Running HTTP client integration example..."
	cd integration/http_client && go run main.go

# Help target
help:
	@echo "Available targets:"
	@echo "  all                    - Run all examples"
	@echo "  basic-all             - Run all basic examples"
	@echo "  advanced-all          - Run all advanced examples"
	@echo "  integration-all       - Run all integration examples"
	@echo "  clean                 - Clean up temporary files"
	@echo ""
	@echo "Basic examples:"
	@echo "  basic-run-simple      - Run simple circuit breaker example"
	@echo ""
	@echo "Advanced examples:"
	@echo "  advanced-run-failure-detection - Run custom failure detection example"
	@echo ""
	@echo "Integration examples:"
	@echo "  integration-run-http-client    - Run HTTP client integration example" 