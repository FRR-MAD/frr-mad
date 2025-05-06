BACKEND_SRC := src/backend
FRONTEND_SRC := src/frontend
TEMPCLIENT_SRC:= tempClient


PROTO_SRC := protobufSource/protocol.proto
PROTO_BACKEND_DEST := src/backend/pkg
PROTO_TEMPCLIENT_DEST := tempClient/pkg
PROTO_FRONTEND_DEST := src/frontend/pkg

.PHONY: run/backend run/backend/local run/backend/prod
run/backend:
	@cd $(BACKEND_SRC) && go mod tidy
	cd $(BACKEND_SRC) && go run -tags=dev cmd/frr-analytics/main.go

run/backend/local:
	@cd $(BACKEND_SRC) && go mod tidy
	cd $(BACKEND_SRC) && go run -tags=local cmd/frr-analytics/main.go

run/backend/prod:
	@cd $(BACKEND_SRC) && go mod tidy
	cd $(BACKEND_SRC) && go run cmd/frr-analytics/main.go

run/frontend:
	@cd $(FRONTEND_SRC) && go mod tidy
	cd $(FRONTEND_SRC) && go run cmd/tui/main.go


.PHONY: binaries
binaries:
	mkdir -p binaries

.PHONY: build/all build/frontend build/backend build/backend/prod build/local
build/all: build/frontend build/backend

.PHONY: build/frontend
build/frontend: binaries
	cd $(FRONTEND_SRC) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -tags=dev -o /tmp/frr-tui ./cmd/tui

build/backend: binaries
	@cd $(BACKEND_SRC) && go mod tidy
	cd $(BACKEND_SRC) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -tags=dev -o ../../binaries/analyzer_dev ./cmd/frr-analytics

build/local: binaries
	@cd $(BACKEND_SRC) && go mod tidy
	cd $(BACKEND_SRC) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -tags=local -o ../../binaries/analyzer_dev ./cmd/frr-analytics


build/backend/prod: binaries
	@cd $(BACKEND_SRC) && go mod tidy
	cd $(BACKEND_SRC) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o ../../binaries/analyzer ./cmd/frr-analytics


.PHONY: protobuf protobuf/mac protobuf/clean
protobuf: protobuf/clean
	protoc --proto_path=protobufSource --go_out=paths=source_relative:$(PROTO_BACKEND_DEST) protobufSource/protocol.proto
	protoc --proto_path=protobufSource --go_out=paths=source_relative:$(PROTO_FRONTEND_DEST) protobufSource/protocol.proto
#	protoc --proto_path=protobufSource --go_out=paths=source_relative:$(PROTO_TEMPCLIENT_DEST) protobufSource/protocol.proto

protobuf/mac: protobuf/clean
	protoc --proto_path=protobufSource --go_out=paths=source_relative:$(PROTO_BACKEND_DEST) protobufSource/protocol.proto
	protoc --proto_path=protobufSource --go_out=paths=source_relative:$(PROTO_FRONTEND_DEST) protobufSource/protocol.proto
#	protoc --proto_path=protobufSource --go_out=paths=source_relative:$(PROTO_TEMPCLIENT_DEST) protobufSource/protocol.proto

protobuf/clean:
	rm -f $(PROTO_BACKEND_DEST)/protocol.pb.go
	rm -f $(PROTO_FRONTEND_DEST)/protocol.pb.go
#	rm -f $(PROTO_TEMPCLIENT_DEST)/protocol.pb.go


#### Hot Module Reloading ####

.PHONY: hmr/docker hmr/run hmr/stop hmr/clean
hmr/docker:
	docker build -t frr-854-dev -f dockerfile/frr-dev.dockerfile .
	docker build -t frr-854 -f dockerfile/frr.dockerfile .

hmr/run:
	cd containerlab && chmod +x scripts/
	-cd containerlab && sh scripts/custom-bridges.sh
	cd containerlab && clab deploy --topo frr01-dev.clab.yml --reconfigure
	cd containerlab && sh scripts/pc-interfaces.sh
	cd containerlab && sh scripts/remove-default-route.sh

hmr/stop: 
	cd containerlab && clab destroy --topo frr01-dev.clab.yml --cleanup


hmr/restart: hmr/stop hmr/run

hmr/clean: hmr/stop
	docker container list -a -q | xargs -i{} docker container rm {}
	docker network prune -f


### Testing
.PHONY: test/backend
test/backend:
	cd src/backend && go test -v ./test/...