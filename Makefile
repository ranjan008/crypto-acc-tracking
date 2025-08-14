# Crypto Account Tracking Makefile

.PHONY: build run test clean install help

# Build variables
BINARY_NAME=crypto-tracker
BUILD_DIR=./bin
MAIN_PATH=./main.go

# Sample addresses for testing
SAMPLE_SMALL=0xa39b189482f984388a34460636fea9eb181ad1a6
SAMPLE_MEDIUM=0xd620AADaBaA20d2af700853C4504028cba7C3333
SAMPLE_LARGE=0xfb50526f49894b78541b776f5aaefe43e3bd8590

# Default target
all: build

## Build the application
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## Build for multiple platforms
build-all:
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	
	@echo "✅ Multi-platform build complete"

## Run the application with sample address
run:
	@echo "🚀 Running with sample address..."
	@go run $(MAIN_PATH) -a $(SAMPLE_SMALL) -o sample_output.csv

## Run with medium sample address
run-medium:
	@echo "🚀 Running with medium sample address..."
	@go run $(MAIN_PATH) -a $(SAMPLE_MEDIUM) -o medium_output.csv

## Run with large sample address (requires API key)
run-large:
	@if [ -z "$(API_KEY)" ]; then \
		echo "❌ Error: API_KEY environment variable is required for large addresses"; \
		echo "Usage: make run-large API_KEY=your_api_key"; \
		exit 1; \
	fi
	@echo "🚀 Running with large sample address..."
	@go run $(MAIN_PATH) -a $(SAMPLE_LARGE) -k $(API_KEY) -o large_output.csv

## Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod tidy
	@go mod download
	@echo "✅ Dependencies installed"

## Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

## Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

## Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f *.csv
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

## Install the binary to $GOPATH/bin
install: build
	@echo "📦 Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "✅ Installed to $(GOPATH)/bin/$(BINARY_NAME)"

## Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

## Run linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run
	@echo "✅ Linting complete"

## Show help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  run          - Run with small sample address"
	@echo "  run-medium   - Run with medium sample address"
	@echo "  run-large    - Run with large sample address (requires API_KEY)"
	@echo "  deps         - Install dependencies"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make run"
	@echo "  make run-large API_KEY=your_etherscan_api_key"
