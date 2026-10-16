# Copilot instructions

Read [`AGENTS.md`](../AGENTS.md) for the project’s product, domain, documentation, and SDD rules. It is the source for requirements workflow and domain integrity; do not restate user stories or acceptance criteria outside executable feature files.

## Build, test, and lint

Commands run from the repository root:

| Purpose | Command |
| --- | --- |
| Build | `make -C backend build` |
| Unit tests (no database) | `make -C backend test` |
| One unit test | `cd backend && go test ./internal/domain -run '^TestNewMoneyValid$'` |
| BDD tests | `make -C backend bdd` |
| Lint | `make -C backend lint` |
| CI checks (lint and full tests) | `make -C backend ci` |

BDD and CI tests need Docker and PostgreSQL. BDD scenarios use the same database and truncate `expense_items` and `expenses` before each scenario; do not run test suites concurrently against that database. The Godog test wrapper fixes its feature path and excludes `@todo`, so it has no single-scenario selector.

## Architecture

The Go module is in [`backend/`](../backend/). The API uses REST and JSON, with a hexagonal architecture:

- [`backend/cmd/api/`](../backend/cmd/api/main.go) is the composition root: it reads configuration, connects to PostgreSQL, and wires the service to the HTTP router.
- [`backend/internal/domain/`](../backend/internal/domain/expense.go) owns business types and invariants. The application layer in [`backend/internal/application/`](../backend/internal/application/expense_service.go) implements use cases and defines ports such as `ExpenseRepository`.
- [`backend/internal/adapters/http/`](../backend/internal/adapters/http/expense_handler.go) handles HTTP input/output and invokes the application. [`backend/internal/adapters/postgres/`](../backend/internal/adapters/postgres/expense_repository.go) implements persistence; its SQL migrations are embedded and applied on connection.
- [`backend/features/`](../backend/features/) contains executable Gherkin requirements. The Godog suite starts the real API and uses PostgreSQL, so it covers the HTTP-to-database path.

Keep business rules in the domain, use cases and ports in the application, and transport or storage concerns in adapters. Map errors at the HTTP boundary to fixed, safe problem responses; do not expose internal error details.

## Repository-specific conventions

- `.feature` files carry the user story, acceptance criteria, and executable scenarios. Keep them under `backend/features/<domain>/`; tag unfinished scenarios `@todo` and use slice tags such as `@s1`. See [`SDD.md`](../SDD.md) for artifact roles and traceability.
- Implement approved behavior one acceptance scenario at a time. Check the feature’s decisions and slice specification before changing behavior; unresolved decisions need approval. For Expense Registration, see [`decisions.md`](../product/epics/financial-management/features/expense-registration/decisions.md) and [`specs/`](../specs/).
- Keep domain terms, code, tests, and documentation in English. When adding a core glossary term, add its Spanish wording as provisional in [`product/glossary.md`](../product/glossary.md); follow [`adr/0003-domain-language.md`](../adr/0003-domain-language.md).
- Represent money with `shopspring/decimal`, not floating point. Preserve the entered decimal scale in API and database representations where the existing expense code does so.
- Add schema changes as numbered up/down SQL migrations under `backend/internal/adapters/postgres/migrations/`; the PostgreSQL adapter embeds and applies them at startup.
- New Go linters or formatters should follow [`backend/.golangci.yml`](../backend/.golangci.yml): `golangci-lint` runs `gofmt`, `goimports`, and exhaustive switch checks.
