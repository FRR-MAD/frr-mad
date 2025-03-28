.PHONY: lint test build run clean

lint:
	golangci-lint run
#	staticcheck ./...
	errcheck ./...
#	govulncheck ./...

test:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# main folder still reeds renaming
build:
	go build -o bin/yourtui cmd/yourtui/main.go

# main folder still reeds renaming
run:
	go run cmd/yourtui/main.go

clean:
	rm -rf bin/ coverage.out coverage.html