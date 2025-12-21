.PHONY: help build test test-coverage clean run-local deploy lint fmt vet install-deps terraform-init terraform-plan terraform-apply terraform-destroy

# Variables
BINARY_NAME=bootstrap
LAMBDA_ZIP=lambda.zip
GO_FILES=$(shell find . -name '*.go' -type f)
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html
TERRAFORM_DIR=terraform/lambda
AWS_REGION?=us-east-1

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
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY_NAME) ./cmd/lambda
	@echo "Binary built: $(BINARY_NAME)"

build-mac: clean ## Build for macOS (development)
	@echo "Building for macOS..."
	go build -o $(BINARY_NAME)-mac ./cmd/lambda
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

ci: install-deps fmt vet test-coverage build ## CI pipeline (Complete validation)

ci-fast: fmt vet test build ## CI pipeline (Fast - no coverage)

# ==============================================================================
# Terraform Targets
# ==============================================================================

terraform-init: ## Initialize Terraform
	@echo "Initializing Terraform..."
	cd $(TERRAFORM_DIR) && terraform init

terraform-validate: ## Validate Terraform configuration
	@echo "Validating Terraform..."
	cd $(TERRAFORM_DIR) && terraform validate

terraform-fmt: ## Format Terraform files
	@echo "Formatting Terraform files..."
	terraform fmt -recursive terraform/

terraform-plan: build ## Run Terraform plan (requires Lambda binary)
	@echo "Planning Terraform changes..."
	@if [ ! -f $(BINARY_NAME) ]; then \
		echo "Error: Lambda binary not found. Run 'make build' first"; \
		exit 1; \
	fi
	cd $(TERRAFORM_DIR) && terraform plan -out=tfplan

terraform-apply: ## Apply Terraform changes
	@echo "Applying Terraform changes..."
	@if [ ! -f $(TERRAFORM_DIR)/tfplan ]; then \
		echo "Error: tfplan not found. Run 'make terraform-plan' first"; \
		exit 1; \
	fi
	cd $(TERRAFORM_DIR) && terraform apply tfplan
	@rm -f $(TERRAFORM_DIR)/tfplan
	@echo "✅ Deployment completed!"
	@echo ""
	@echo "API Gateway URL:"
	@cd $(TERRAFORM_DIR) && terraform output -raw api_gateway_url
	@echo ""

terraform-apply-auto: build ## Build and apply Terraform (auto-approve - use with caution!)
	@echo "Building and deploying Lambda..."
	cd $(TERRAFORM_DIR) && terraform apply -auto-approve
	@echo "✅ Deployment completed!"

terraform-output: ## Show Terraform outputs
	@cd $(TERRAFORM_DIR) && terraform output

terraform-destroy: ## Destroy Terraform infrastructure
	@echo "⚠️  WARNING: This will destroy all infrastructure!"
	@echo "Press Ctrl+C to cancel, or Enter to continue..."
	@read confirm
	cd $(TERRAFORM_DIR) && terraform destroy

terraform-clean: ## Clean Terraform files
	@echo "Cleaning Terraform files..."
	@rm -rf $(TERRAFORM_DIR)/.terraform $(TERRAFORM_DIR)/.terraform.lock.hcl
	@rm -f $(TERRAFORM_DIR)/tfplan $(TERRAFORM_DIR)/terraform.tfstate*
	@rm -f $(TERRAFORM_DIR)/lambda-deployment.zip
	@echo "Terraform files cleaned"

# Complete deployment workflow
deploy-infra: build terraform-plan terraform-apply ## Complete deployment workflow (build + plan + apply)

# Quick redeploy (updates Lambda code only via Terraform)
redeploy: build terraform-apply-auto ## Quick redeploy (build + auto-apply)

# Test deployed Lambda
test-lambda: ## Test deployed Lambda endpoint
	@echo "Testing Lambda endpoint..."
	@API_URL=$$(cd $(TERRAFORM_DIR) && terraform output -raw api_gateway_url 2>/dev/null); \
	if [ -z "$$API_URL" ]; then \
		echo "Error: API Gateway URL not found. Deploy first with 'make deploy-infra'"; \
		exit 1; \
	fi; \
	echo "Testing health endpoint: $$API_URL/health"; \
	curl -s $$API_URL/health | jq; \
	echo ""; \
	echo "Testing auth endpoint: $$API_URL/auth"; \
	curl -s -X POST $$API_URL/auth \
		-H "Content-Type: application/json" \
		-d '{"email":"admin@oficinapro.com","senha":"senha123"}' | jq

# View Lambda logs
logs-lambda: ## View Lambda CloudWatch logs
	@FUNCTION_NAME=$$(cd $(TERRAFORM_DIR) && terraform output -raw lambda_function_name 2>/dev/null); \
	if [ -z "$$FUNCTION_NAME" ]; then \
		echo "Error: Lambda function name not found."; \
		exit 1; \
	fi; \
	aws logs tail /aws/lambda/$$FUNCTION_NAME --follow --region $(AWS_REGION)

.DEFAULT_GOAL := help
