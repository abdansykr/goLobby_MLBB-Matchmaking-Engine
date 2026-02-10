.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down

# Variables
APP_NAME=matchmaking
DOCKER_IMAGE=antigravity/matchmaking
GO_FILES=$(shell find . -name '*.go' -type f)

# Help command
help:
	@echo "Antigravity Matchmaking - Available Commands:"
	@echo ""
	@echo "  make build         - Build the Go binary"
	@echo "  make run           - Run the application locally"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo ""
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start all services with Docker Compose"
	@echo "  make docker-down   - Stop all Docker services"
	@echo "  make docker-logs   - View Docker logs"
	@echo ""
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Rollback last migration"
	@echo ""

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) ./cmd/server
	@echo "Build complete!"

# Run the application
run:
	@echo "Starting $(APP_NAME)..."
	go run ./cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Test coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(APP_NAME)
	rm -f coverage.out coverage.html
	go clean
	@echo "Clean complete!"

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):latest .
	@echo "Docker image built!"

# Start Docker Compose services
docker-up:
	@echo "Starting services..."
	docker-compose up -d
	@echo "Services started! Check status with: docker-compose ps"

# Stop Docker Compose services
docker-down:
	@echo "Stopping services..."
	docker-compose down
	@echo "Services stopped!"

# View Docker logs
docker-logs:
	docker-compose logs -f app

# Run database migrations up
migrate-up:
	@echo "Running migrations..."
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/antigravity?sslmode=disable" up
	@echo "Migrations complete!"

# Rollback last migration
migrate-down:
	@echo "Rolling back migration..."
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/antigravity?sslmode=disable" down 1
	@echo "Rollback complete!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed!"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted!"

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run
	@echo "Linting complete!"

# Full deployment workflow
deploy: clean build docker-build docker-up
	@echo "Deployment complete!"
