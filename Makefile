.PHONY: build run test clean docker-build docker-up docker-down migrate-up migrate-down help

# Build the application
build:
	@echo "Building application..."
	go build -o bin/main ./cmd/main.go

# Run the application
run:
	@echo "Running application..."
	go run ./cmd/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f main

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t tz1-app .

# Start Docker Compose
docker-up:
	@echo "Starting Docker Compose..."
	docker-compose up -d

# Stop Docker Compose
docker-down:
	@echo "Stopping Docker Compose..."
	docker-compose down

# Run migrations up
migrate-up:
	@echo "Running migrations up..."
	@if command -v goose >/dev/null 2>&1; then \
		goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=device_db sslmode=disable" up; \
	else \
		echo "goose not installed. Install with: go install github.com/pressly/goose/v3/cmd/goose@latest"; \
	fi

# Run migrations down
migrate-down:
	@echo "Running migrations down..."
	@if command -v goose >/dev/null 2>&1; then \
		goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=device_db sslmode=disable" down; \
	else \
		echo "goose not installed. Install with: go install github.com/pressly/goose/v3/cmd/goose@latest"; \
	fi

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Install goose for migrations
install-goose:
	@echo "Installing goose..."
	go install github.com/pressly/goose/v3/cmd/goose@latest

# Create new migration
migrate-create:
	@echo "Creating new migration..."
	@if command -v goose >/dev/null 2>&1; then \
		goose -dir migrations create $(name) sql; \
	else \
		echo "goose not installed. Run: make install-goose"; \
	fi

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-up      - Start Docker Compose"
	@echo "  docker-down    - Stop Docker Compose"
	@echo "  migrate-up     - Run migrations up"
	@echo "  migrate-down   - Run migrations down"
	@echo "  migrate-create - Create new migration (make migrate-create name=migration_name)"
	@echo "  install-goose  - Install goose migration tool"
	@echo "  deps           - Install Go dependencies"
	@echo "  help           - Show this help message"
