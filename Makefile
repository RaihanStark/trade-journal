.PHONY: help run build migrate-up migrate-down migrate-create sqlc-generate dev docker-up docker-down docker-logs

help:
	@echo "Available commands:"
	@echo "  make run            - Run the API server"
	@echo "  make build          - Build the API server"
	@echo "  make dev            - Run in development mode with hot reload"
	@echo "  make docker-up      - Start PostgreSQL in Docker"
	@echo "  make docker-down    - Stop PostgreSQL in Docker"
	@echo "  make docker-logs    - View PostgreSQL logs"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback database migrations"
	@echo "  make migrate-create - Create a new migration (usage: make migrate-create name=migration_name)"
	@echo "  make sqlc-generate  - Generate sqlc code"

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
