.PHONY: all help build run dev test clean \
        migrations-status migrations-version migrations-create \
        migrations-up migrations-down migrations-reset migrations-redo \
        migrations-up-to migrations-down-to \
		sqlc-gen sqlc-vet

# Default target
all: help

help:
	@echo "Available commands:"
	@echo "  build              Build the application"
	@echo "  run                Run the application"
	@echo "  dev                Run the application with live reload (docker-compose)"
	@echo "  test               Run tests"
	@echo "  clean              Clean build artifacts"
	@echo "  migrations-status  Show goose migration status"
	@echo "  migrations-create  Create a new migration (make migrations-create name=foo)"
	@echo "  migrations-up      Apply all up migrations"
	@echo "  migrations-down    Roll back one migration"
	@echo "  migrations-reset    Roll back ALL migrations to base (dangerous in prod!)"
	@echo "  migrations-redo       Re-run the latest migration (down then up)"
	@echo "  migrations-up-to      Migrate up to a specific version (use: make migrations-up-to version=<num>)"
	@echo "  migrations-down-to    Migrate down to a specific version (use: make migrations-down-to version=<num>)"
	@echo "  db-seed            Seed the database with test data (make db-seed)"
	@echo "  sqlc-gen           Generate sqlc code (make sqlc-gen)"
	@echo "  sqlc-vet           Validate sqlc without generating (make sqlc-vet)"
	@echo ""
	@echo "Tip: Run 'make migrations-status' to see the version numbers."
	
# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run the application with live reload
dev:
	docker-compose up -d pgweb
	docker-compose up -d db
	docker-compose up dev; docker-compose down

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/ tmp/


migrations-status:
	docker compose exec dev goose status

migrations-create:
	if [ -z "$$name" ]; then \
		name="migration"; \
	fi; \
	docker compose exec dev goose -s create $$name sql

migrations-up:
	@echo Running goose migrations up
	docker compose exec dev goose up

migrations-down:
	@echo Running goose migrations down
	docker compose exec dev goose down

migrations-reset:
	@echo "Resetting ALL migrations"
	docker compose exec dev goose reset

migrations-redo:
	@echo "Re-running latest migration (down then up)"
	docker compose exec dev goose redo

# Migrate to a specific version (pass version=<number>)
migrations-up-to:
	@echo "Migrating UP to version $(version)"
	docker compose exec dev goose up-to $(version)

migrations-down-to:
	@echo "Migrating DOWN to version $(version)"
	docker compose exec dev goose down-to $(version)

## Seed the database with test data
db-seed:
	@echo "Seeding base types"
	docker compose exec dev sh -c 'psql $$DATABASE_URL -f db/seeds/seed.sql'
	@echo "Seeding development data..."
	docker compose exec dev sh -c 'psql -d "$$DATABASE_URL" -f db/seeds/seed.dev.sql'

## Generate sqlc code
sqlc-gen: 
	docker compose exec dev sqlc generate

## Validate code without generating
sqlc-vet:
	docker compose exec dev sqlc vet