# Makefile for Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
BINARY_NAME=cal-blog-service

# .Phony is directive that tells Make the following commands are not actual files
.PHONY: all build start test clean lint fmt help

all: lint test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go

# Start the local server
start:
	air

# Run tests
test:
	$(GOTEST) -v -cover ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	gofmt -s -w .

# Update dependencies
tidy:
	$(GOMOD) tidy

# Show help
help:
	@echo "Usage:"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make lint     - Run linter"
	@echo "  make fmt      - Format code"
	@echo "  make clean    - Clean build files"
	@echo "  make tidy     - Update dependencies"
	@echo "  make all      - Run lint, test, and build"