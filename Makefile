include .env
export

# ── Variables ─────────────────────────────────────────────────────────────────
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATIONS_DIR=./migrations
BINARY=hotel-booking

# ── Run ───────────────────────────────────────────────────────────────────────
.PHONY: run
run:
	go run main.go

.PHONY: build
build:
	go build -o $(BINARY) main.go

# ── Migrations ────────────────────────────────────────────────────────────────
.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

.PHONY: migrate-down-1
migrate-down-1:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

.PHONY: migrate-status
migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

.PHONY: migrate-force
migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(VERSION)

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)

# ── Database ──────────────────────────────────────────────────────────────────
.PHONY: db-create
db-create:
	createdb --username=$(DB_USER) --host=$(DB_HOST) $(DB_NAME)

.PHONY: db-drop
db-drop:
	dropdb --username=$(DB_USER) --host=$(DB_HOST) $(DB_NAME)

# ── Tools ─────────────────────────────────────────────────────────────────────
.PHONY: install-tools
install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test ./... -v

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run              - Run the server"
	@echo "  make build            - Build the binary"
	@echo "  make migrate-up       - Apply all migrations"
	@echo "  make migrate-down     - Rollback all migrations"
	@echo "  make migrate-down-1   - Rollback last migration"
	@echo "  make migrate-status   - Show current migration version"
	@echo "  make migrate-create NAME=migration_name  - Create new migration"
	@echo "  make migrate-force VERSION=N  - Force migration version"
	@echo "  make db-create        - Create database"
	@echo "  make db-drop          - Drop database"
	@echo "  make install-tools    - Install golang-migrate"
	@echo "  make tidy             - Tidy go modules"
