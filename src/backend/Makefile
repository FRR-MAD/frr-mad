# Backend Makefile

include $(ROOT_DIR)/Makefile.inc

ANALYZER_OUTPUT_NAME := frr-mad-analyzer
OUTPUT_PATH := $(BINARY_PATH)/$(ANALYZER_OUTPUT_NAME)

.PHONY: all dev prod test test-analyzer test-aggregator test-exporter test-comms clean

# Default target
all: dev

# Development build with dev tags
dev:
	@echo "Building backend (development)..."
	mkdir -p $(BINARY_PATH)
	$(GO_FLAGS) go build -ldflags "$(LD_FLAGS)" -tags=dev -o $(OUTPUT_PATH) ./cmd/frr-analyzer

# Production build
prod:
	@echo "Building backend (production)..."
	mkdir -p $(BINARY_PATH)
	$(GO_FLAGS) go build -ldflags "$(LD_FLAGS)" -o $(OUTPUT_PATH) ./cmd/frr-analyzer

# Run all backend tests
test: test-analyzer test-aggregator test-exporter test-comms

# Individual test targets
test-analyzer:
	@echo "Running analyzer tests..."
	go test -v ./test/analyzer/...

test-aggregator:
	@echo "Running aggregator tests..."
	go test -v ./test/aggregator/...

test-exporter:
	@echo "Running exporter tests..."
	go test -v ./test/exporter/...

test-comms:
	@echo "Running communication tests..."
	go test -v ./test/comms/...

# Clean backend artifacts
clean:
	@echo "Cleaning backend artifacts..."
	rm -f $(OUTPUT_PATH)