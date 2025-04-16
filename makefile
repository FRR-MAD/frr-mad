
PROTO_SRC := protobufSource/protocol.proto
BACKEND_DEST := backend/pkg
TEMP_CLIENT_DEST := temporary-client/pkg
FRONTEND_DEST := frontend/pkg

.PHONY: run/backend run/backend/local run/backend/prod
run/backend:
	@cd backend && go mod tidy
	cd backend && go run -tags=dev cmd/frr-analytics/main.go

run/backend/local:
	@cd backend && go mod tidy
	cd backend && go run -tags=local cmd/frr-analytics/main.go

run/backend/prod:
	@cd backend && go mod tidy
	cd backend && go run cmd/frr-analytics/main.go

run/frontend:
	@cd frontend && go mod tidy
	cd frontend && go run cmd/tui/main.go


.PHONY: binaries
binaries:
	mkdir -p binaries

.PHONY: build/all build/frontend build/backend build/backend/prod
build/all: build/frontend build/backend

#.PHONY: build/frontend
#build/frontend: binaries
#	cd frontend && GOOS=linux GOARCH=amd64 go build -o  ../binaries

build/backend: binaries
	cd backend && go mod tidy && GOOS=linux GOARCH=amd64 go build -ldflags='-extldflags=-static' -tags=dev -o ../binaries/analyzer ./cmd/frr-analytics

build/backend/prod: binaries
	cd backend && go mod tidy && GOOS=linux GOARCH=amd64 go build -o ../binaries/analyzer ./cmd/frr-analytics


.PHONY: protobuf protobuf/clean
protobuf: protobuf/clean
	@chmod +x proto-binary/bin/protoc
	./proto-binary/bin/protoc --proto_path=protobufSource --go_out=paths=source_relative:backend/pkg protobufSource/protocol.proto
	./proto-binary/bin/protoc --proto_path=protobufSource --go_out=paths=source_relative:frontend/pkg protobufSource/protocol.proto
	./proto-binary/bin/protoc --proto_path=protobufSource --go_out=paths=source_relative:tempClient/pkg protobufSource/protocol.proto

protobuf/clean:
	rm -f $(BACKEND_DEST)/protocol.pb.go
	rm -f $(TEMP_CLIENT_DEST)/protocol.pb.go
	rm -f $(FRONTEND_DEST)/protocol.pb.go

