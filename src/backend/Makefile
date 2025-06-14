# Backend Makefile

export ROOT_DIR := $(shell pwd)/../..
include $(ROOT_DIR)/Makefile.inc

ANALYZER_OUTPUT_NAME := frr-mad-analyzer
OUTPUT_PATH := $(BINARY_PATH)/$(ANALYZER_OUTPUT_NAME)
TEST_CASE := ""

.PHONY: all dev prod test test/analyzer test/aggregator test/exporter test/socket test/config clean

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
test: test/analyzer test/aggregator test/exporter test/socket test/config

# Individual test targets
test/analyzer:
	@echo "Running analyzer tests..."
ifeq ($(TEST_CASE),"")
	go test -v ./test/analyzer/... 
else
	go test -v ./test/analyzer -run $(TEST_CASE)
endif

test/aggregator:
	@echo "Running aggregator tests..."
	go test -v ./test/aggregator/...

test/exporter:
	@echo "Running exporter tests..."
	go test -v ./test/exporter/...

test/socket:
	@echo "Running communication tests..."
	go test -v ./test/socket/...

test/config:
	@echo "Running config tests..."
	go test -v ./test/configs/...

# Clean backend artifacts
clean:
	@echo "Cleaning backend artifacts..."
	rm -f $(OUTPUT_PATH)
