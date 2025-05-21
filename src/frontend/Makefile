# Frontend Makefile

include $(ROOT_DIR)/Makefile.inc

TUI_OUTPUT_NAME := frr-mad-tui
OUTPUT_PATH := $(BINARY_PATH)/$(TUI_OUTPUT_NAME)

.PHONY: all dev prod test clean

# Default target
all: dev

# Development build with dev tags
dev:
	@echo "Building frontend (development)..."
	mkdir -p $(BINARY_PATH)
	$(GO_FLAGS) go build -ldflags "$(LD_FLAGS)" -tags=dev -o $(OUTPUT_PATH) ./cmd/tui

# Production build
prod:
	@echo "Building frontend (production)..."
	mkdir -p $(BINARY_PATH)
	$(GO_FLAGS) go build -ldflags "$(LD_FLAGS)" -o $(OUTPUT_PATH) ./cmd/tui

# Run tests
test:
	@echo "Running frontend tests..."
#	go test -v ./test/...

# Clean frontend artifacts
clean:
	@echo "Cleaning frontend artifacts..."
	rm -f $(OUTPUT_PATH)