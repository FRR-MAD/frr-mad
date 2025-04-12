
.PHONY: run/backend run/backend/prod
run/backend:
	@cd backend && go mod tidy
	cd backend && go run -tags=dev cmd/frr-analytics/main.go

run/backend/prod:
	@cd backend && go mod tidy
	cd backend && go run cmd/frr-analytics/main.go

run/frontend:
	@cd frontend && go mod tidy
	cd frontend && go run cmd/tui/main.go


.PHONY: binaries
binaries:
	mkdir -p binaries

.PHONY: build/all build/frontend build/backend/dev build/backend/prod
build/all: build/frontend build/backend

#.PHONY: build/frontend
#build/frontend: binaries
#	cd frontend && GOOS=linux GOARCH=amd64 go build -o  ../binaries

build/backend/dev: binaries
	cd backend && go mod tidy && GOOS=linux GOARCH=amd64 go build -tags=dev -o ../binaries/analyzer ./cmd/frr-analytics

build/backend/prod: binaries
	cd backend && GOOS=linux GOARCH=amd64 go build -o ../binaries/analyzer ./cmd/frr-analytics


#.PHONY: translate/proto
#translate/proto:
#	protoc --go_out=. --go_opt=paths=source_relative \
#		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
#		./proto/*.proto