# Spendtrack

[![Coverage Status](https://coveralls.io/repos/github/eprediger/spendtrack/badge.svg?branch=main)](https://coveralls.io/github/eprediger/spendtrack?branch=main)

Keep track of income, expense and savings for a healthy budget.

Go backend, hexagonal architecture, PostgreSQL persistence, REST + JSON API.
BDD requirements live in `features/` (Gherkin, executed by godog).

## Quick Start

### Prerequisites

- Go 1.27+
- Docker, Docker Compose
- `golangci-lint` (for `make lint` / `make ci`)
- Make (optional, for convenience commands)

### Git Hooks

**This repository sets up the following git hooks:**

- `post-checkout`: automates the process of executing the `setup-hooks.sh` script when the repository gets cloned
- `pre-commit`: verifies that the linter, unit and e2e tests run successfully
- `commit-msg`: validates locally that the commits messages follows the [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/)

### **Build and start:**

```bash
make run
```

### **Access the API:**

- `POST /api/v1/expenses` — register an expense
- `GET /api/v1/expenses/{id}` — review an expense
- `curl http://localhost:8000/api/v1/health` — health check

## Available Commands

### Using Make (recommended):
```bash
make help           # Show all available commands
```

## Development Workflow && Testing Strategy

This project adopts a Behavior-Driven Development methodology for its development
and a Hexagonal Architecture as the software design pattern leveraging its
testability, flexibility, maintainability and framework agnosticism

### 1. Unit Tests (go test)
- Located next to the code (`internal/**/…_test.go`)
- Test individual components in isolation
- Run with: `make test`

### 2. BDD Tests (godog)
- Located in `features/`
- Test user scenarios end-to-end against real Postgres
- Written in Gherkin syntax
- Run with: `make bdd`
- Use the isolated `spendtrack_test` database; `spendtrack` remains for manual app use
- Direct `go test ./features/...` also defaults to `spendtrack_test` (create it first with `make test-db`)

### 3. Coverage Reporting
- Run with: `make coverage`
- The total merges unit tests and the BDD-driven API binary

## Adding New Features

1. **Domain First**: Define entities and invariants in `internal/domain/`
2. **Application**: Implement use cases and ports in `internal/application/`
3. **Infrastructure**:
   1. **Driving side**: adapters under `internal/adapters/http/` depend on the port implemented by the Application
   2. **Driven side**: adapters under `internal/adapters/postgres/` implement ports the Application depends on
4. **Tests**: Add unit tests and BDD scenarios
5. **Documentation**: Update this README

## Code Quality

### Linting with golangci-lint

- Go linter aggregator; configuration in `.golangci.yml`
- Includes `exhaustive` (switch exhaustiveness over enums)
- Run with: `make lint` or `make lint-fix`

---

## Nice-to-have features

- [ ] Add authentication and authorization

## Nice-to-have non-functional requirements

- [ ] Generate CHANGELOG.md automatically
- [x] Define tree structured with a Hexagonal Architecture approach
- [x] Linting rules (that executes previous a commit)
- [ ] Testing setup
  - [x] unit
  - [ ] integration
  - [x] e2e
  - [ ] performance
  - [ ] stress
  - [ ] mutation
- [ ] Coverage rules
- [ ] Environment configuration && deployment (Set up CI/CD pipeline)
- [x] Implement proper error handling [(RFC 9457: Problem Details for HTTP APIs)](https://www.rfc-editor.org/rfc/rfc9457.html#name-the-problem-details-json-ob)
- [ ] Add API rate limiting
- [ ] Graceful shutdown improvements
  - [x] Shutdown timeout (8s) below Docker's default 10s stop grace period
    - Why 8s: the app is not deployed and there is no data about request durations. Docker Compose is the only platform. 8s is the longest drain that still leaves 2s for `pool.Close()` and process exit before Docker sends SIGKILL.
    - Missing information to tune it:
      - The deploy platform and its grace period (for example, Kubernetes `terminationGracePeriodSeconds`, 30s by default)
      - Request durations per endpoint (p99 and maximum) from production metrics or performance tests
      - The time that a load balancer needs to stop sending traffic to an instance that is shutting down
  - [x] Release the signal handler after the first signal, so a second Ctrl-C stops the API at once
  - [x] Fail the BDD suite when the API does not exit cleanly, so BDD coverage cannot drop silently
  - [x] Log when shutdown starts and when it ends
  - [ ] Cancel in-flight request contexts when the shutdown timeout expires
    - Deferred: each request is one short database round trip, so no request is likely to run for 8s. `srv.Shutdown` does not cancel request contexts, so a slow query continues until the process exits.
    - Do this when an endpoint can run longer than the shutdown timeout (for example, reports or imports). Give the server a `BaseContext` and cancel it when the timeout expires.
  - [ ] Read the shutdown timeout from configuration when the app has more than one deploy target
- [ ] Versioning (SemVer)
- [ ] Vulnerability checks for dependencies
- [ ] Monitoring && observability
  - [ ] Structured logging
  - [ ] Monitoring
  - [ ] Healthcheck
- [ ] Security
  - [ ] Configure SAST
  - [ ] Container scanning
  - [ ] Define a secrets management policy
- [ ] Technical documentation
