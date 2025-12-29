# Makefile for shop-retail-employee-service
#
# To use this Makefile, you should have a .env file in the root of the project,
# or have the required environment variables (DB_USER, DB_PASSWORD, DB_NAME) exported.
# You can copy .env.example to .env to get started.
#
# The 'godotenv' package in the go application will load the .env file for the 'run' target.
# For 'docker-up' and 'migrate', the variables need to be available in your shell environment.
# A simple way to do this is to run 'export $(cat .env | xargs)' on Linux/macOS.
# On Windows, you might need to set them manually or use a script.

# Go variables
MODULE_PATH=github.com/zuyatna/shop-retail-employee-service

.PHONY: help run docker-up docker-down migrate

all: help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run          - Run the Go application (loads .env file)."
	@echo "  docker-up    - Start PostgreSQL container (requires DB_* env vars)."
	@echo "  docker-down  - Stop and remove PostgreSQL container."
	@echo "  migrate      - Apply database migrations (requires DB_* env vars)."
	@echo ""
	@echo "Note: For 'docker-up' and 'migrate', ensure DB_USER, DB_PASSWORD, and DB_NAME are set in your environment."


run:
	@echo "Running the Go application..."
	go run $(MODULE_PATH)/cmd/api


docker-up:
	@echo "Starting PostgreSQL container..."
	docker run --name shop-retail \
	  -e POSTGRES_USER=$(DB_USER) \
	  -e POSTGRES_PASSWORD=$(DB_PASSWORD) \
	  -e POSTGRES_DB=$(DB_NAME) \
	  -p 5432:5432 \
	  -d postgres:15

docker-down:
	@echo "Stopping and removing PostgreSQL container..."
	docker stop shop-retail || true
	docker rm shop-retail || true

migrate:
	@echo "Applying database migrations..."
	docker cp migrations/001_init.sql shop-retail:/001_init.sql
	docker cp migrations/002_seed_supervisor.sql shop-retail:/002_seed_supervisor.sql
	docker exec -it shop-retail psql -U $(DB_USER) -d $(DB_NAME) -f 001_init.sql
	docker exec -it shop-retail psql -U $(DB_USER) -d $(DB_NAME) -f 002_seed_supervisor.sql
