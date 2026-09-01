.PHONY: dev infra api web setup clean migrate

# All services (infra + api + web)
dev: infra
	@echo ""
	@echo "  Thankly Dev Environment"
	@echo "  ─────────────────────────"
	@echo "  Frontend:  http://localhost:4321"
	@echo "  API:       http://localhost:8080"
	@echo "  Postgres:  localhost:5432"
	@echo "  Redis:     localhost:6379"
	@echo ""
	@make -j2 api web

# Infrastructure only (Docker)
infra:
	docker compose up -d
	@echo "Waiting for services..."
	@sleep 2
	@docker compose ps

# Go backend (native)
api:
	cd backend && go run ./cmd/server

# Astro frontend (native)
web:
	cd frontend && bun run dev

# Setup project for first time
setup:
	cp -n .env.example .env 2>/dev/null || true
	make infra
	@sleep 3
	@echo "Running migrations..."
	@docker compose exec -T postgres psql -U thankly -d thankly -f /docker-entrypoint-initdb.d/001_initial_schema.sql
	cd frontend && bun install
	@echo ""
	@echo "Setup complete!"
	@echo "Run 'make dev' to start all services."

# Stop everything
clean:
	docker compose down

# Run migrations
migrate:
	docker compose exec -T postgres psql -U thankly -d thankly -f /docker-entrypoint-initdb.d/001_initial_schema.sql

# Backend tests
test:
	cd backend && go test ./...

# Backend build
build:
	cd backend && go build -o bin/thankly-api ./cmd/server
