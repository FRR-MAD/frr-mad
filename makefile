
.PHONY: run/backend
run/backend:
	@cd backend && go mod tidy
	cd backend && go run cmd/frr-analytics/main.go

.PHONY: run/frontend
run/frontend:
	@cd frontend && go mod tidy
	cd frontend && go run cmd/tui/main.go


.PHONY: binaries
binaries:
	mkdir -p binaries


.PHONY: build/all
build/all: build/frontend build/backend

#.PHONY: build/frontend
#build/frontend: binaries
#	cd frontend && GOOS=linux GOARCH=amd64 go build -o  ../binaries

.PHONY: build/backend
build/backend: binaries
	cd frontend && GOOS=linux GOARCH=amd64 go build -o ../binaries
