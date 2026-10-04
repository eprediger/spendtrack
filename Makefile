.PHONY: build run up down logs db lint lint-fix test bdd coverage ci clean dev

TARGETS := build run up down logs db lint lint-fix test bdd coverage ci clean dev

$(TARGETS):
	$(MAKE) -C backend $@

help: ## Show this help message
	$(MAKE) -C backend help
