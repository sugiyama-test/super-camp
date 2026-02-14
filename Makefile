.PHONY: up down build logs migrate-up migrate-down health

# Docker
up:
	docker compose up -d

up-build:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f

logs-api:
	docker compose logs -f api

logs-frontend:
	docker compose logs -f frontend

# Database
MIGRATE=docker compose exec api migrate
DB_URL=postgres://supercamp:supercamp@db:5432/supercamp_dev?sslmode=disable

migrate-up:
	$(MIGRATE) -path /app/migrations -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path /app/migrations -database "$(DB_URL)" down 1

# Development
api-test:
	cd backend && go test ./...

frontend-lint:
	cd frontend && npm run lint

frontend-build:
	cd frontend && npm run build

# Quick checks
health:
	curl -s http://localhost:8081/api/health | python3 -m json.tool
