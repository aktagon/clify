# Makefile for clify

.PHONY: build run clean test lint fmt deps help install

# Build the binary
build:
	go build -o bin/clify main.go

# Run the application in interactive mode
run:
	go run main.go

# Run with a specific query
run-query:
	go run main.go "$(QUERY)"

# Clean build artifacts
clean:
	rm -f bin/clify

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Install the binary to $GOPATH/bin
install: build
	go install .

# Show help
help:
	@echo "Available targets:"
	@echo "  build         Build the binary"
	@echo "  run           Run the application in interactive mode"
	@echo "  run-query     Run with a specific query (use QUERY=...)"
	@echo "  clean         Clean build artifacts"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  lint          Run linter"
	@echo "  fmt           Format code"
	@echo "  deps          Download and tidy dependencies"
	@echo "  install       Install binary to GOPATH/bin"
	@echo "  help          Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make run"
	@echo "  make run-query QUERY='find all .txt files'"
