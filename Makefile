# FRR-MAD: Top-level Makefile

# Set root directory for use in included files
export ROOT_DIR := $(shell pwd)
export BINARY_PATH := $(ROOT_DIR)/build

# Include common definitions
include Makefile.inc

.PHONY: all dev clean test protobuf help frontend backend

# Default target
all: protobuf frontend-prod backend-prod

# Development build
dev: protobuf frontend backend

# Frontend targets
frontend:
	$(MAKE) -C src/frontend dev

frontend-prod:
	$(MAKE) -C src/frontend prod

# Backend targets
backend:
	$(MAKE) -C src/backend dev

backend-prod:
	$(MAKE) -C src/backend prod

# Testing targets
test: test-frontend test-backend

test-frontend:
	$(MAKE) -C src/frontend test

test-backend:
	$(MAKE) -C src/backend test

# Protobuf generation
protobuf:
	@echo "Generating protobuf files..."
	mkdir -p src/frontend/pkg src/backend/pkg
	protoc --proto_path=$(PROTO_DIR) --go_out=paths=source_relative:src/backend/pkg $(PROTO_SRC)
	protoc --proto_path=$(PROTO_DIR) --go_out=paths=source_relative:src/frontend/pkg $(PROTO_SRC)

# Go module management
go-tidy:
	cd src/backend && go mod tidy
	cd src/frontend && go mod tidy

go-sync: go-tidy
	cd src && go work sync
	cd src && go work vendor

# Development environment
hmr-docker:
	docker build -t frr-854-dev -f dockerfile/frr-dev.dockerfile .
	docker build -t frr-854 -f dockerfile/frr.dockerfile .

# More HMR targets here...

# Clean everything
clean:
	$(MAKE) -C src/frontend clean
	$(MAKE) -C src/backend clean
	rm -rf $(BINARY_PATH)
	go clean -testcache

# Help information
help:
	@echo "FRR-MAD Makefile Help"
	@echo ""
	@echo "Main targets:"
	@echo "  all              - Build production frontend and backend"
	@echo "  dev              - Build development frontend and backend"
	@echo "  clean            - Clean all build artifacts"
	@echo "  test             - Run all tests"
	@echo ""
	@echo "Component targets:"
	@echo "  frontend         - Build frontend (development)"
	@echo "  frontend-prod    - Build frontend (production)"
	@echo "  backend          - Build backend (development)"
	@echo "  backend-prod     - Build backend (production)"
	@echo ""
	@echo "For more options, see component-specific Makefiles in src/frontend and src/backend"

#BINARY_PATH := ./build
#BACKEND_SRC := src/backend
#FRONTEND_SRC := src/frontend
#TEMPCLIENT_SRC:= tempClient
#
#PROTO_SRC := protobufSource/protocol.proto
#PROTO_BACKEND_DEST := src/backend/pkg
#PROTO_TEMPCLIENT_DEST := tempClient/pkg
#PROTO_FRONTEND_DEST := src/frontend/pkg
#
#TUI_OUTPUT_NAME := frr-mad-tui
#ANALYZER_OUTPUT_NAME := frr-mad-analyzer
#
#DAEMON_VERSION := $(shell grep "daemon=" VERSION | cut -d'=' -f2)
#TUI_VERSION := $(shell grep "tui=" VERSION | cut -d'=' -f2)
#COMMIT := $(shell git rev-parse --short HEAD)
#BUILD_DATE := $(shell date -u '+%Y-%m-%d\:%H:%M:%S')
#
#LD_FLAGS = -s \
#	-X main.TUIVersion=$(TUI_VERSION) \
#	-X main.DaemonVersion=$(DAEMON_VERSION) \
#	-X main.GitCommit=$(COMMIT) \
#	-X main.BuildDate=$(BUILD_DATE)
#
#.PHONY: all
#
#all: go/tidy go/sync build/backend/prod  build/frontend/prod
#
#dev: go/tidy go/sync build/backend build/frontend
#
#.PHONY: binaries build/all build/frontend build/backend build/backend/prod
#binaries:
#	mkdir -p $(BINARY_PATH)
#
#build/all: build/frontend build/backend
#
#build/frontend: binaries
#	cd $(FRONTEND_SRC) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LD_FLAGS)" -tags=dev -o ../../$(BINARY_PATH)/$(TUI_OUTPUT_NAME) ./cmd/tui
#
#build/backend: binaries
#	cd $(BACKEND_SRC) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LD_FLAGS)" -tags=dev -o ../../$(BINARY_PATH)/$(ANALYZER_OUTPUT_NAME) ./cmd/frr-analyzer
#
#build/frontend/prod: binaries
#	cd $(FRONTEND_SRC) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LD_FLAGS)" -o ../../$(BINARY_PATH)/$(TUI_OUTPUT_NAME) ./cmd/tui
#
#build/backend/prod: binaries
#	cd $(BACKEND_SRC) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LD_FLAGS)" -o ../../$(BINARY_PATH)/$(ANALYZER_OUTPUT_NAME) ./cmd/frr-analyzer
#
#.PHONY: protobuf protobuf/mac protobuf/clean
#protobuf: protobuf/clean
#	protoc --proto_path=protobufSource --go_out=paths=source_relative:$(PROTO_BACKEND_DEST) protobufSource/protocol.proto
#	protoc --proto_path=protobufSource --go_out=paths=source_relative:$(PROTO_FRONTEND_DEST) protobufSource/protocol.proto
#
#protobuf/mac: protobuf/clean
#	protoc --proto_path=protobufSource --go_out=paths=source_relative:$(PROTO_BACKEND_DEST) protobufSource/protocol.proto
#	protoc --proto_path=protobufSource --go_out=paths=source_relative:$(PROTO_FRONTEND_DEST) protobufSource/protocol.proto
#
#protobuf/clean:
#	rm -f $(PROTO_BACKEND_DEST)/protocol.pb.go
#	rm -f $(PROTO_FRONTEND_DEST)/protocol.pb.go
##	rm -f $(PROTO_TEMPCLIENT_DEST)/protocol.pb.go
#
#
##### Hot Module Reloading ####
#
#.PHONY: hmr/docker hmr/run hmr/stop hmr/clean hmr/restart
#hmr/docker:
#	docker build -t frr-854-dev -f dockerfile/frr-dev.dockerfile .
#	docker build -t frr-854 -f dockerfile/frr.dockerfile .
#
#hmr/run:
#	cd containerlab && chmod +x scripts/
#	-cd containerlab && sh scripts/custom-bridges.sh
#	cd containerlab && clab deploy --topo frr01-dev.clab.yml --reconfigure
#	cd containerlab && sh scripts/pc-interfaces.sh
#	cd containerlab && sh scripts/remove-default-route.sh
#
#hmr/stop: 
#	cd containerlab && clab destroy --topo frr01-dev.clab.yml --cleanup
#
#
#hmr/restart: hmr/stop hmr/run
#
#hmr/clean: hmr/stop
#	docker container list -a -q | xargs -i{} docker container rm {}
#	docker network prune -f
#
#.PHONY: go/mod go/sync
#go/tidy:
#	cd src/backend && go mod tidy
#	cd src/frontend && go mod tidy
#
#go/sync: go/tidy
#	cd src && go work sync 
#	cd src && go work vendor
#
#
#### Testing
#.PHONY: test/all test/backend test/analyzer test/aggregator test/exporter test/clean
#test/all: test/backend test/clean
#
#test/frontend: test/clean
#	cd src/frontend && go test -v ./test/...
#
#test/backend: test/clean
#	cd src/backend && go test -v ./test/...
#
#test/analyzer: test/clean
#	cd src/backend && go test -v ./test/analyzer/...
#
#test/aggregator: test/clean
#	cd src/backend && go test -v ./test/aggregator/...
#
#test/exporter: test/clean
#	cd src/backend && go test -v ./test/exporter/...
#
#test/comms: test/clean
#	cd src/backend && go test -v ./test/comms/...
#
#test/clean:
#	go clean -testcache