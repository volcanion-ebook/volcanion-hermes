# Makefile for Volcanion Hermes

# Variables
APP_NAME=volcanion-hermes
SERVER_PATH=cmd/server/main.go
BUILD_DIR=bin
DOCKER_IMAGE=$(APP_NAME):latest

# Default target
.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
.PHONY: dev
dev: ## Run the application in development mode
	go run $(SERVER_PATH)

.PHONY: build
build: ## Build the application
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SERVER_PATH)

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)

# Dependencies
.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod verify

.PHONY: deps-update
deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

# Testing
.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Code quality
.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

# Docker
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Database
.PHONY: mongo-start
mongo-start: ## Start MongoDB container
	docker run --name mongodb -d -p 27017:27017 mongo:latest

.PHONY: mongo-stop
mongo-stop: ## Stop MongoDB container
	docker stop mongodb && docker rm mongodb

.PHONY: minio-start
minio-start: ## Start MinIO container
	docker run --name minio -d \
		-p 9000:9000 \
		-p 9001:9001 \
		-e "MINIO_ROOT_USER=minioadmin" \
		-e "MINIO_ROOT_PASSWORD=minioadmin" \
		minio/minio server /data --console-address ":9001"

.PHONY: minio-stop
minio-stop: ## Stop MinIO container
	docker stop minio && docker rm minio

# Development environment
.PHONY: dev-up
dev-up: mongo-start minio-start ## Start development environment
	@echo "Development environment started"
	@echo "MongoDB: http://localhost:27017"
	@echo "MinIO Console: http://localhost:9001"

.PHONY: dev-down
dev-down: mongo-stop minio-stop ## Stop development environment
	@echo "Development environment stopped"

# Production
.PHONY: prod-build
prod-build: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(APP_NAME) $(SERVER_PATH)

.PHONY: install
install: build ## Install the binary
	cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
