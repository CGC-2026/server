.PHONY: run build test clean dev

# Default target
all: build

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run the application with live reload
dev:
	docker-compose up dev

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/ tmp/
