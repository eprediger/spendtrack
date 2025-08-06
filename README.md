# Spendtrack
git
[![Coverage Status](https://coveralls.io/repos/github/eprediger/spendtrack/badge.svg?branch=main)](https://coveralls.io/github/eprediger/spendtrack?branch=main)

Keep track of income, expense and savings for a healthy budget

## Quick Start

## Prerequisites

- Docker, Docker Compose or Python 3.13+
- Make (optional, for convenience commands)

### Git Hooks

**This repository sets up the following git hooks:**

- `post-checkout`: automates the process of executing the `setup-hooks.sh` script when the repository gets cloned
- `pre-commit`: verifies that the linter, unit and e2e tests run successfully
- `commit-msg`: validates locally that the commits messages follows the [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/)

### **Build and start:**

 ```bash
 make build
 make up
 ```

### **Access the API:**
- API docs: http://localhost:8000/docs (OpenAPI Specification)
- Health endpoint: `curl http://localhost:8000/api/v1/health` (URL Path versioning)

## Available Commands

### Using Make (recommended):
```bash
make help           # Show all available commands
```

## Development Workflow && Testing Strategy

This project adopts a Behavior-Driven Development methodology for its development 
and a Hexagonal Architecture as the software design pattern leveraging its 
testability, flexibility, maintainability and framework agnosticism 

### 1. Unit Tests (pytest)
- Located in `tests/`
- Test individual components in isolation
- Mock external dependencies
- Run with: `make test`

### 2. BDD Tests (behave)
- Located in `features/`
- Test user scenarios end-to-end
- Written in Gherkin syntax
- Run with: `make behave`

### 3. Coverage Reporting
- Tracks code coverage across all tests
- Generates HTML reports in `htmlcov/`
- Run with: `make coverage`

1. **Start development environment:**
   ```bash
   make up
   ```

2. **Make code changes** (hot reload enabled)

3. **Run tests:**
   ```bash
   make test-all
   ```

4. **Check coverage:**
   ```bash
   make coverage
   ```

5. **Fix linting issues:**
   ```bash
   make lint-fix
   ```

## Adding New Features

1. **Domain First**: Define entities and ports in `app/application/domain/`
2. **Application**: Implement services in `app/application/`
3. **Infrastructure**: 
   1. **Driving side**: Primary actors initiate the interaction. depends on the port, which is implemented by the Application
   2. **Driven side**: Secondary actors are kicked into behavior by the Application, and this depends on the defined Ports, which this side implements
4. **Tests**: Add unit tests and BDD scenarios
5. **Documentation**: Update this README

## Code Quality

### Linting with Ruff

- Fast Python linter and formatter
- Replaces flake8, isort, and black
- Configuration in `pyproject.toml`
- Run with: `make lint` or `make lint-fix`

**Code Style Rules**:
- Line length: 88 characters
- Import sorting and formatting
- Python 3.13+ syntax enforcement
- Comprehensive error checking

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
- [ ] Implement proper error handling [(RFC 9457: Problem Details for HTTP APIs)](https://www.rfc-editor.org/rfc/rfc9457.html#name-the-problem-details-json-ob)
- [ ] Add API rate limiting
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
