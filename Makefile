.PHONY: build up start stop down logs shell install lint lint-ci lint-fix \
		test test-ci coverage behave behave-ci test-all ci clean

TARGETS := build up start stop down logs shell install lint lint-ci lint-fix \
		test test-ci coverage behave behave-ci test-all ci clean

$(TARGETS):
	$(MAKE) -C backend $@

help: ## Show this help message
	$(MAKE) -C backend help
