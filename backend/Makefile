.PHONY: help build run up down logs db lint lint-fix test bdd coverage ci clean

DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/spendtrack?sslmode=disable

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

lint: ## Run golangci-lint
	golangci-lint run ./...

lint-fix: ## Run golangci-lint with autofix
	golangci-lint run --fix ./...

test: ## Run unit tests (no database needed)
	go test ./internal/... ./cmd/...

bdd: db ## Run BDD suite against real Postgres
	DATABASE_URL="$(DATABASE_URL)" go test ./features/...

coverage: db ## Run all tests with coverage report
	DATABASE_URL="$(DATABASE_URL)" go test -coverprofile=coverage.out -coverpkg=./internal/...,./cmd/... ./...
	go tool cover -func=coverage.out | tail -1

ci: lint db ## Lint + full test suite (used by hooks/CI)
	DATABASE_URL="$(DATABASE_URL)" go test -coverprofile=coverage.out -coverpkg=./internal/...,./cmd/... ./...
	go tool cover -func=coverage.out | tail -1

clean: ## Stop the stack and delete volumes
	docker compose down -v
