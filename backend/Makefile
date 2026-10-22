.PHONY: help build run up down logs db test-db lint lint-fix test bdd coverage coverage-run ci clean

DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/spendtrack?sslmode=disable
TEST_DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/spendtrack_test?sslmode=disable

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go binary
	go build ./...

dev: db ## Rebuild + restart the API on save (host only, requires `air`; not used for deploys)
	DATABASE_URL="$(DATABASE_URL)" air

run: ## Build and run the full stack (app + postgres)
	docker compose up --build

up: ## Start the full stack in background
	docker compose up -d --wait

down: ## Stop the stack
	docker compose down

logs: ## Follow container logs
	docker compose logs -f

db: ## Start postgres only (for local tests)
	docker compose up -d --wait postgres

test-db: db ## Ensure the isolated test database exists
	@docker compose exec -T postgres psql -U postgres -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname = 'spendtrack_test'" | grep -q 1 || docker compose exec -T postgres createdb -U postgres spendtrack_test

lint: ## Run golangci-lint
	golangci-lint run ./...

lint-fix: ## Run golangci-lint with autofix
	golangci-lint run --fix ./...

test: ## Run unit tests (no database needed)
	go test ./internal/... ./cmd/...

bdd: test-db ## Run BDD suite against the isolated test database
	DATABASE_URL="$(TEST_DATABASE_URL)" go test ./features/...

coverage: test-db ## Run all tests with coverage report
	DATABASE_URL="$(TEST_DATABASE_URL)" $(MAKE) coverage-run

coverage-run: ## Run unit + BDD coverage (needs DATABASE_URL in env)
	rm -rf .coverdata && mkdir -p .coverdata/unit .coverdata/bdd
	go test -count=1 -cover -coverpkg=./internal/...,./cmd/... ./internal/... ./cmd/... -args -test.gocoverdir=$(abspath .coverdata/unit)
	GOCOVERDIR=$(abspath .coverdata/bdd) go test -count=1 ./features/...
	go tool covdata textfmt -i=.coverdata/unit,.coverdata/bdd -o coverage.out
	go tool cover -func=coverage.out | tail -1

ci: lint test-db ## Lint + full test suite (used by hooks/CI)
	DATABASE_URL="$(TEST_DATABASE_URL)" $(MAKE) coverage-run

clean: ## Stop the stack and delete volumes
	docker compose down -v
