.PHONY: help build install test lint clean fmt vet deps docker-build docker-run release

# Version information - source of truth is sdk/version.go
SDK_VERSION ?= $(shell grep 'const Version' sdk/version.go | sed 's/.*"\(.*\)".*/\1/')
VERSION ?= $(shell echo $(SDK_VERSION) | sed 's/^v//')
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build variables
BINARY_NAME = upwind
CMD_DIR = ./cmd/upwind
BUILD_DIR = ./build
LDFLAGS = -ldflags "-X main.version=$(SDK_VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE) -w -s"

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOMOD = $(GOCMD) mod
GOVET = $(GOCMD) vet
GOFMT = gofmt
GOLINT = golangci-lint

help: ## Display this help message
	@echo "Upwind Go SDK - Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

build: deps ## Build the CLI binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: deps ## Build binaries for all platforms
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	@echo "All binaries built in $(BUILD_DIR)/"

install: build ## Install the CLI binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) $(CMD_DIR)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "Tests complete. Coverage report: coverage.out"

test-coverage: test ## Run tests and display coverage report
	@echo "Displaying coverage report..."
	$(GOCMD) tool cover -html=coverage.out

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) -w -s .
	@echo "Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run ./...; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out
	@echo "Clean complete"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build \
		--build-arg VERSION=$(SDK_VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg DATE=$(DATE) \
		-t upwind-go-sdk:$(VERSION) .
	@echo "Docker image built: upwind-go-sdk:$(VERSION)"

docker-run: docker-build ## Run CLI in Docker container
	@echo "Running in Docker..."
	docker run --rm -it \
		-e UPWIND_CLIENT_ID \
		-e UPWIND_CLIENT_SECRET \
		-e UPWIND_ORGANIZATION_ID \
		-e UPWIND_REGION \
		upwind-go-sdk:$(VERSION) $(filter-out $@,$(MAKECMDGOALS))

release: clean build-all ## Create a release build
	@echo "Creating release $(VERSION)..."
	@mkdir -p $(BUILD_DIR)/release
	@cd $(BUILD_DIR) && \
		tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
		tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64 && \
		tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
		tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64 && \
		zip -q release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Release packages created in $(BUILD_DIR)/release/"
	@ls -lh $(BUILD_DIR)/release/

# Allow passing arguments to docker-run
%:
	@:

