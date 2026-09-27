# S1 — Technical Plan

**Status:** approved

Implementation plan for [`spec.md`](spec.md). This plan also covers migrating the existing Python scaffold to Go (per [`adr/0001`](../../../adr/0001-backend-language.md)).

## Technical Decisions

| Concern        | Decision                          | Rationale                                                        |
| -------------- | --------------------------------- | ---------------------------------------------------------------- |
| HTTP router    | stdlib `net/http` (Go 1.22+ patterns) | Method+path routing (`POST /api/v1/expenses`) covers S1 with zero deps; add `chi` only if middleware needs grow |
| Money          | `shopspring/decimal`              | Exact decimal arithmetic; float64 for money is a correctness bug  |
| IDs            | `google/uuid`                     | UUID v7 generated in the app: non-enumerable, k-sortable, identity known before persistence |
| Postgres       | `pgx/v5` (`pgxpool`)              | Standard, driver + pool; plain SQL keeps the adapter transparent  |
| Migrations     | `golang-migrate`                  | `migrations/` SQL files, versioned                                |
| BDD            | `godog`                           | Go counterpart of behave; `.feature` files unchanged              |
| Unit tests     | stdlib `testing` + `testify`      | Assertions readability                                            |
| Lint           | `golangci-lint` incl. `exhaustive`| `exhaustive` closes Go's missing sum-type check (per ADR-0001)    |
| Config         | env vars + `caarlos0/env`         | Minimal typed config                                              |

`sqlc` was considered for DB access (generates type-safe Go from SQL) — deferred: adds codegen tooling; revisit when queries grow.

## Target Layout

```text
backend/
  cmd/api/main.go                 # composition root
  internal/
    domain/
      money.go                    # Money value object (amount + currency, invariants)
      expense.go                  # Expense entity + NewExpense(...) (invariants here)
      errors.go                   # domain errors
    application/
      ports.go                    # ExpenseRepository (driven port)
      expense_service.go          # RegisterExpense, GetExpense (use cases)
    adapters/
      http/
        router.go                 # routes
        expense_handler.go        # decode → service → encode
        problem.go                # RFC 9457 problem+json writer
      postgres/
        expense_repository.go     # pgx adapter for ExpenseRepository
        migrations/               # golang-migrate SQL files
  features/                       # unchanged .feature files
    steps/
      expense_steps.go            # godog bindings
    expenses_test.go              # godog suite runner
  go.mod
  Dockerfile                      # multi-stage → static binary
  docker-compose.yml              # app + postgres
```

Python scaffold (`app/`, `tests/`, `features/environment.py`, `features/steps/*.py`, `pyproject.toml`, `poetry.lock`) is removed — history preserves it; Python returns later for analytics as its own directory.

## Layering (unchanged from existing convention)

* `domain/` — entities, value objects, invariants. No imports from other layers.
* `application/` — use cases + port interfaces. Depends only on domain.
* `adapters/http`, `adapters/postgres` — driving/driven adapters; depend on ports.
* `cmd/api` — wires adapters to ports.

## Validation Mapping

| `.feature` rule                     | Enforced by                                            |
| ----------------------------------- | ------------------------------------------------------ |
| required fields                     | domain constructor (`NewExpense`) + request decoding    |
| `amount > 0`                        | `Money`/`Expense` invariant                            |
| currency is valid ISO 4217          | currency list embedded in domain                       |
| `occurred_at` instant not future    | `Expense` invariant vs. `time.Now()` (instant comparison) |
| rejected → per-field errors         | `problem.go` renders 422 + `errors[]`                  |

## Migration Steps

1. Go module + layout + tooling (golangci-lint config with `exhaustive`).
2. `docker-compose`: add `postgres` service; Dockerfile → multi-stage Go build.
3. Makefile targets → `go build/test`, `godog`, `golangci-lint`, `migrate`; root Makefile unchanged (delegates).
4. CI: Go setup/test/lint/godog instead of poetry/pytest/behave; coverage stays via `go test -cover`.
5. Postgres migration `0001_expenses` → `expenses` table per spec.
6. Domain → service → repository → handler, TDD per scenario.
7. godog steps for `@s1` scenarios; remove `@todo` from `register_basic_expense.feature`.
8. Health endpoint re-implemented in Go (existing behavior preserved — `health.feature` keeps passing).

## Verification

`make ci` must run green: `golangci-lint`, `go test ./...`, `godog` — including the previously-skipped `@todo` S1 scenarios now executed for real.

## Known Costs Accepted

* Go migration rewrites scaffold + step bindings (one-time).
* Exhaustiveness/invariants by convention + `exhaustive` linter, not compiler.
