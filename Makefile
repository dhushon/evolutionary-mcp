# Simple development helpers

# `make run` brings up supporting containers (Postgres/Redis) and runs the
# Go backend in the foreground.  The compose phase also ensures the database is
# configured with the credentials expected by config.yaml.

.PHONY: run db compose-up compose-down build

run: db
	@echo "Starting backend..."
	# wait for Postgres to accept connections before starting server
	until docker-compose exec -T postgres pg_isready -U user > /dev/null 2>&1; do \
		echo "waiting for postgres..."; sleep 1; \
	done
	cd backend && go run ./cmd/server

# start only the database (and optionally redis) services
# other services (frontend, ml-sidecar) are optional and may be started
# independently via docker-compose.
# the "redis" service is not strictly required; if pulling it fails you can
# comment it out or run it manually.
db:
	@echo "starting postgres (and redis) containers..."
	@if ! docker-compose up -d postgres; then \
		echo "\nERROR: Docker compose could not start the postgres container."; \
		echo "Please ensure your Docker daemon is running and can pull images."; \
		echo "You can also start a local Postgres manually or set DB connection via env vars."; \
		exit 1; \
	fi
	@docker-compose up -d redis || echo "redis service failed to start; continue anyway"

# full composition including build
compose-up:
	docker-compose up --build -d

compose-down:
	docker-compose down

build:
	go build ./...
