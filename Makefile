.PHONY: dev-frontend dev-backend install tidy migrate build up down

## Development
dev-frontend:
	cd frontend && pnpm dev

dev-backend:
	cd backend && go run cmd/api/main.go serve

## Setup
install:
	cd frontend && pnpm install

tidy:
	cd backend && go mod tidy

## Migrations — usage: make migrate name=add_posts_collection
migrate:
	@[ "$(name)" ] || (echo "Usage: make migrate name=<migration_name>" && exit 1)
	cd backend && go run cmd/api/main.go migrate create $(name)

## Docker
build:
	docker compose build

up:
	docker compose up --build

down:
	docker compose down
