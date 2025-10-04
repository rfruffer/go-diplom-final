.PHONY: build build-server build-client clean test test-cover lint vet fmt generate proto docker-build docker-run help

# Build info
VERSION ?= 1.0.0
BUILD_DATE ?= $(shell date -u +"%Y-%m-%d %H:%M:%S UTC")

# Go build flags
LDFLAGS = -ldflags "-X 'github.com/rfruffer/go-diplom-final/pkg/version.Version=$(VERSION)' \
                   -X 'github.com/rfruffer/go-diplom-final/pkg/version.BuildDate=$(BUILD_DATE)'"

# Build directories
BUILD_DIR = build
SERVER_BINARY = $(BUILD_DIR)/server
CLIENT_BINARY = $(BUILD_DIR)/client

# Default target
all: build

## Build all binaries
build: build-server build-client

## Build server binary
build-server:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(SERVER_BINARY) ./cmd/server

## Build client binary
build-client:
	@echo "Building client..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(CLIENT_BINARY) ./cmd/client

## Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/server-linux-amd64 ./cmd/server
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/client-linux-amd64 ./cmd/client
	# Windows
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/server-windows-amd64.exe ./cmd/server
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/client-windows-amd64.exe ./cmd/client
	# macOS Intel
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/server-darwin-amd64 ./cmd/server
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/client-darwin-amd64 ./cmd/client
	# macOS Apple Silicon
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/server-darwin-arm64 ./cmd/server
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/client-darwin-arm64 ./cmd/client

## Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

## Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## Run unit tests
test-unit:
	@echo "Running unit tests..."
	@go test -v ./tests/unit/...

## Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v ./tests/integration/...

## Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

## Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

## Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

## Generate code
generate:
	@echo "Generating code..."
	@go generate ./...

## Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/api/proto/*.proto

## Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

## Install development tools
dev-deps:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

## Build Docker images
docker-build:
	@echo "Building Docker images..."
	@docker build -t gophkeeper-server -f docker/Dockerfile.server .
	@docker build -t gophkeeper-client -f docker/Dockerfile.client .

## Run with Docker Compose
docker-run:
	@echo "Running with Docker Compose..."
	@docker-compose up -d

## Show help
help:
	@echo "Available commands:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //' | awk -F ':' '{printf "  %-20s %s\n", $$1, $$2}'