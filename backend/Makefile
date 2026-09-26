.PHONY: help build up start stop down logs shell install lint lint-ci lint-fix \
		test test-ci coverage behave behave-ci test-all ci clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Docker containers
	docker-compose build

up: ## Run the application containers
	docker-compose up -d

start: ## Start the application
	docker-compose start

stop: ## Stop the application
	docker-compose stop

down: ## Stop the application and delete all artifacts
	docker-compose down

logs: ## Show application logs
	docker-compose logs -f

shell: ## Enter the application container
	docker-compose exec app bash

install: ## Install dependencies
	docker-compose exec app poetry install

lint: ## Run linter (interactive)
	docker-compose exec app poetry run ruff check app/
	docker-compose exec app poetry run ruff format --check app/

lint-ci: ## Run linter (non-interactive, for CI/hooks)
	docker-compose exec -T app poetry run ruff check app/
	docker-compose exec -T app poetry run ruff format --check app/

lint-fix: ## Fix linting issues
	docker-compose exec app poetry run ruff check --fix app/
	docker-compose exec app poetry run ruff format app/

test: ## Run unit tests with coverage (interactive)
	docker-compose exec app poetry run -m pytest
	docker-compose exec app poetry run coverage report

test-ci: ## Run unit tests with coverage (non-interactive, for CI/hooks)
	docker-compose exec -T app poetry run coverage run -m pytest
	docker-compose exec -T app poetry run coverage report

coverage: ## Generate coverage report
	docker-compose exec app poetry run coverage run -m pytest
	docker-compose exec app poetry run coverage run -a -m behave
	docker-compose exec app poetry run coverage report --show-missing
	docker-compose exec app poetry run coverage html

behave: ## Run BDD tests
	docker-compose exec app poetry run behave --tags=~@todo

behave-ci: ## Run BDD tests (non-interactive, for CI/hooks)
	docker-compose exec -T app poetry run behave --tags=~@todo

test-all: lint test behave ## Run all tests and linting (interactive)

ci: lint-ci test-ci behave-ci ## Run all tests and linting (non-interactive, for CI/hooks)

clean: ## Clean up containers and volumes
	docker-compose down -v
	docker system prune -f
