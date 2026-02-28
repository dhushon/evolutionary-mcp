# Simple development helpers

# `make run` brings up supporting containers (Postgres/Redis) and runs the
# Go backend in the foreground.  The compose phase also ensures the database is
# configured with the credentials expected by config.yaml.

.PHONY: run db compose-up compose-down build clean setup-env deps

run: db
	@echo "Starting backend..."
	# wait for Postgres to accept connections before starting server
	until docker-compose exec -T postgres pg_isready -U user > /dev/null 2>&1; do \
		echo "waiting for postgres..." >&2; sleep 1; \
	done
	cd backend && go run ./cmd/server

# start only the database (and optionally redis) services
# other services (frontend, ml-sidecar) are optional and may be started
# independently via docker-compose.
db:
	@echo "starting postgres (and redis) containers..."
	@if ! docker-compose up -d postgres; then \
		echo "\nERROR: Docker compose could not start the postgres container."; \
		echo "Please ensure your Docker daemon is running and can pull images."; \
		echo "You can also start a local Postgres manually or set DB connection via env vars."; \
		exit 1; \
	fi; \
	# The "redis" service is not strictly required for `make run`; if it fails,
	# we log a warning but continue.
	@docker-compose up -d redis || echo "redis service failed to start; continue anyway"

# full composition including build
compose-up:
	docker-compose up --build -d

compose-down:
	docker-compose down

build:
	go build ./...

clean:
	@echo "Cleaning Go build artifacts..."
	@cd backend && go clean ./...
	@echo "Stopping and removing Docker containers..."
	@docker-compose down

setup-env:
	@echo "Building and running setup utility..."
	@cd backend && go run ./cmd/setup env
	@mv backend/.env .env 2>/dev/null || true
	@echo "✅ .env file created in project root."

deps:
	@echo "Resolving dependencies..."
	@cd backend && go mod tidy
