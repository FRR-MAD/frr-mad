# FRR-MAD: Top-level Makefile

# Set root directory for use in included files
FRR_VERSION := 8.5.4
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
test: test/clearcache test/frontend test/backend

test/clearcache:
	go clean -testcache

test/frontend:
	$(MAKE) -C src/frontend test

test/backend:
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
hmr/docker:
	cd $(DEV_DIR) && docker build -t frr-img-dev -f dockerfile/frr$(FRR_VERSION)-dev.dockerfile .
	cd $(DEV_DIR) && docker build -t frr-img -f dockerfile/frr$(FRR_VERSION).dockerfile .


hmr/run:
	mkdir -p $(CONTAINERLAB)/.shared/binary/
	cd $(CONTAINERLAB) && chmod +x scripts/
	-cd $(CONTAINERLAB) && sh scripts/custom-bridges.sh
	cd $(CONTAINERLAB) && clab deploy --topo frr01-dev.clab.yml --reconfigure
	cd $(CONTAINERLAB) && sh scripts/pc-interfaces.sh
	cd $(CONTAINERLAB) && sh scripts/remove-default-route.sh

hmr/stop: 
	cd $(CONTAINERLAB) && clab destroy --topo frr01-dev.clab.yml --cleanup

hmr/restart: hmr/stop hmr/run

hmr/clean: hmr/stop
	docker container list -a -q | xargs -i{} docker container rm {}
	docker network prune -f
	rm -rf $(CONTAINERLAB)/.shared

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
