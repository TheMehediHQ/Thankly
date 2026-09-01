.PHONY: dev dev-backend dev-frontend docker-up docker-down migrate test

# Development
dev: docker-up
	@echo "Services starting..."
	@echo "API:      http://localhost:8080"
	@echo "Frontend: http://localhost:4321"
	@echo "Postgres: localhost:5432"
	@echo "Redis:    localhost:6379"

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

# Backend
dev-backend:
	cd backend && go run ./cmd/server

build-backend:
	cd backend && go build -o bin/thankly-api ./cmd/server

test-backend:
	cd backend && go test ./...

test-backend-verbose:
	cd backend && go test -v ./...

# Frontend
dev-frontend:
	cd frontend && npm run dev

build-frontend:
	cd frontend && npm run build

# Database
migrate:
	psql $(DATABASE_URL) -f backend/migrations/001_initial_schema.sql

migrate-docker:
	docker compose exec postgres psql -U thankly -d thankly -f /docker-entrypoint-initdb.d/001_initial_schema.sql

# Setup
setup:
	cp .env.example .env
	docker compose up -d postgres redis
	@echo "Waiting for services..."
	@sleep 3
	docker compose exec postgres psql -U thankly -d thankly -f /docker-entrypoint-initdb.d/001_initial_schema.sql
	cd frontend && npm install
	@echo "Setup complete!"
