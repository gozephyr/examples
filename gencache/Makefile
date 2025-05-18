# Makefile for running gencache examples

.PHONY: all clean
.PHONY: basic-all basic-run-simple basic-run-capacity basic-run-error
.PHONY: advanced-all advanced-run-pooling advanced-run-batch advanced-run-file-store advanced-run-metrics advanced-run-policy

# Default target
all: basic-all advanced-all

# Clean target
clean:
	@echo "Cleaning up..."
	rm -rf /tmp/gencache

# Basic examples
basic-all: basic-run-simple basic-run-capacity basic-run-error

basic-run-simple:
	@echo "Running simple operations example..."
	cd basic/simple_operations && go run main.go

basic-run-capacity:
	@echo "Running capacity limits example..."
	cd basic/capacity_limits && go run main.go

basic-run-error:
	@echo "Running error handling example..."
	cd basic/error_handling && go run main.go

# Advanced examples
advanced-all: advanced-run-pooling advanced-run-batch advanced-run-file-store advanced-run-metrics advanced-run-policy

advanced-run-pooling:
	@echo "Running object pooling example..."
	cd advanced/pooling && go run main.go

advanced-run-batch:
	@echo "Running batch operations example..."
	cd advanced/batch && go run main.go

advanced-run-file-store:
	@echo "Running file store example..."
	cd advanced/file_store && go run main.go

advanced-run-metrics:
	@echo "Running metrics example..."
	cd advanced/metrics && go run main.go

advanced-run-policy:
	@echo "Running eviction policies example..."
	cd advanced/policy && go run main.go

# Help target
help:
	@echo "Available targets:"
	@echo "  all              - Run all examples"
	@echo "  basic-all        - Run all basic examples"
	@echo "  advanced-all     - Run all advanced examples"
	@echo "  clean            - Clean up temporary files"
	@echo ""
	@echo "Basic examples:"
	@echo "  basic-run-simple    - Run simple operations example"
	@echo "  basic-run-capacity  - Run capacity limits example"
	@echo "  basic-run-error     - Run error handling example"
	@echo ""
	@echo "Advanced examples:"
	@echo "  advanced-run-pooling      - Run object pooling example"
	@echo "  advanced-run-batch        - Run batch operations example"
	@echo "  advanced-run-file-store   - Run file store example"
	@echo "  advanced-run-metrics      - Run metrics example"
	@echo "  advanced-run-policy       - Run eviction policies example" 