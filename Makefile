.PHONY: help run build migrate-up migrate-down migrate-create sqlc-generate dev docker-up docker-down docker-logs test test-verbose test-coverage test-race clean seed

help:
	@echo "Available commands:"
	@echo "  make run            - Run the API server"
	@echo "  make build          - Build the API server"
	@echo "  make dev            - Run in development mode with hot reload"
	@echo "  make test           - Run all tests"
	@echo "  make test-verbose   - Run all tests with verbose output"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-race      - Run tests with race detector"
	@echo "  make clean          - Clean build artifacts and test cache"
	@echo "  make docker-up      - Start PostgreSQL in Docker"
	@echo "  make docker-down    - Stop PostgreSQL in Docker"
	@echo "  make docker-logs    - View PostgreSQL logs"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback database migrations"
	@echo "  make migrate-create - Create a new migration (usage: make migrate-create name=migration_name)"
	@echo "  make sqlc-generate  - Generate sqlc code"
	@echo "  make seed           - Seed the database"

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

migrate-up:
	dbmate up

migrate-down:
	dbmate down

migrate-create:
	dbmate new $(name)

sqlc-generate:
	sqlc generate

dev:
	@which air > /dev/null || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	air

docker-up:
	docker compose -f docker-compose-dev.yml up -d

docker-down:
	docker compose -f docker-compose-dev.yml down

docker-logs:
	docker compose -f docker-compose-dev.yml logs -f postgres

test:
	@echo "Running tests..."
	go test ./...

test-verbose:
	@echo "Running tests (verbose)..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...
	@echo ""
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	@echo "Running tests with race detector..."
	go test -race ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -testcache
	@echo "Clean complete"

seed:
	go run cmd/seed/main.go