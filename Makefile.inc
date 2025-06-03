# Makefile.inc - Common definitions for FRR-MAD

# Version information
DAEMON_VERSION := $(shell grep "daemon=" $(ROOT_DIR)/VERSION | cut -d'=' -f2)
TUI_VERSION := $(shell grep "tui=" $(ROOT_DIR)/VERSION | cut -d'=' -f2)
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date -u '+%Y-%m-%d\:%H:%M:%S')

# Common settings
GO_FLAGS := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
LD_FLAGS := -s \
	-X main.TUIVersion=$(TUI_VERSION) \
	-X main.DaemonVersion=$(DAEMON_VERSION) \
	-X main.GitCommit=$(COMMIT) \
	-X main.BuildDate=$(BUILD_DATE)

# Protobuf settings
PROTO_DIR := $(ROOT_DIR)/protobufSource
PROTO_SRC := $(ROOT_DIR)/protobufSource/protocol.proto

# Dev Environment settings
DEV_DIR := .devenvironment
CONTAINERLAB := $(DEV_DIR)/containerlab