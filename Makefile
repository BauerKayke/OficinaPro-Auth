.PHONY: help build test test-coverage clean run-local deploy lint fmt vet install-deps

# Variables
BINARY_NAME=bootstrap
LAMBDA_ZIP=lambda.zip
GO_FILES=$(shell find . -name '*.go' -type f)
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod verify
	@echo "Dependencies installed successfully"

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "Vet completed"

lint: fmt vet ## Run linters
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run --timeout=5m
	@echo "Linting completed"

test: ## Run tests
	@echo "Running tests..."
	go test -v -race ./...
	@echo "Tests completed"

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"
	@echo "Opening coverage report..."
	@open $(COVERAGE_HTML) || xdg-open $(COVERAGE_HTML) || echo "Please open $(COVERAGE_HTML) manually"

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	go test -v -short ./...

build: clean ## Build Lambda binary
	@echo "Building Lambda binary..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_NAME) cmd/lambda/main.go
	@echo "Binary built: $(BINARY_NAME)"

build-mac: clean ## Build for macOS (development)
	@echo "Building for macOS..."
	go build -o $(BINARY_NAME)-mac cmd/lambda/main.go
	@echo "Binary built: $(BINARY_NAME)-mac"

package: build ## Package Lambda for deployment
	@echo "Packaging Lambda..."
	zip $(LAMBDA_ZIP) $(BINARY_NAME)
	@echo "Package created: $(LAMBDA_ZIP)"
	@ls -lh $(LAMBDA_ZIP)

run-local: build-mac ## Run locally (simulates Lambda)
	@echo "Running locally..."
	@echo "Make sure you have .env file configured"
	./$(BINARY_NAME)-mac

deploy: package ## Deploy to AWS Lambda
	@echo "Deploying to AWS Lambda..."
	@if [ -z "$(FUNCTION_NAME)" ]; then \
		echo "Error: FUNCTION_NAME environment variable is not set"; \
		echo "Usage: make deploy FUNCTION_NAME=your-function-name"; \
		exit 1; \
	fi
	aws lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--zip-file fileb://$(LAMBDA_ZIP) \
		--region $(AWS_REGION)
	@echo "Deployment completed"

deploy-staging: ## Deploy to staging environment
	$(MAKE) deploy FUNCTION_NAME=oficinapro-auth-staging

deploy-prod: ## Deploy to production environment
	$(MAKE) deploy FUNCTION_NAME=oficinapro-auth-prod

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-mac $(LAMBDA_ZIP) $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "Clean completed"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t oficinapro-auth:latest .
	@echo "Docker image built"

docker-run: docker-build ## Run in Docker
	@echo "Running in Docker..."
	docker run -p 9000:8080 --env-file .env oficinapro-auth:latest

docker-test: ## Run tests in Docker
	@echo "Running tests in Docker..."
	docker run --rm -v $(PWD):/app -w /app golang:1.21 make test

security-scan: ## Run security scan
	@echo "Running security scan..."
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securego/gosec/v2/cmd/gosec@latest)
	gosec ./...
	@echo "Security scan completed"

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

mod-tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	go mod tidy
	@echo "Modules tidied"

mod-update: ## Update Go dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies updated"

all: clean lint test build ## Run all checks and build

ci: lint test-coverage build ## CI pipeline

.DEFAULT_GOAL := help

