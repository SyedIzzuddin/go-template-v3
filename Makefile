# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integration Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Database operations
migrate-up:
	@echo "Running migrations..."
	@goose -dir db/migrations postgres "$$DATABASE_URL" up

migrate-down:
	@echo "Rolling back migrations..."
	@goose -dir db/migrations postgres "$$DATABASE_URL" down

migrate-status:
	@echo "Migration status..."
	@goose -dir db/migrations postgres "$$DATABASE_URL" status

migrate-create:
	@echo "Creating new migration: $(name)"
	@goose -dir db/migrations create $(name) sql

# Generate code with sqlc
sqlc-generate:
	@echo "Generating sqlc code..."
	@sqlc generate

# Generate Swagger documentation
swagger-gen:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go

# Install tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: all build run test clean watch docker-run docker-down itest migrate-up migrate-down migrate-status migrate-create sqlc-generate swagger-gen install-tools
