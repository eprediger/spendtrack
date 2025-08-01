# Spendtrack

Keep track of income, expense and savings for a healthy budget

## Quick Start

### Git Hooks

**This repository sets up the following git hooks:**
- 
- `post-checkout`: to automate the process of executing the `setup-hooks.sh` script when the repository gets cloned
- `commit-msg`: to validate locally that the commits messages follows the [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/)

### **Build and start:**

 ```bash
 make build
 make up
 ```

### **Access the API:**
   - Health endpoint: `curl http://localhost:8000/api/v1/health`
   - API docs: http://localhost:8000/docs

## Available Commands

### Using Make (recommended):
```bash
make help           # Show all available commands
make build          # Build Docker containers
make up             # Start the application
make stop           # Stop the application
make down           # Stop the application and delete all artifacts
make logs           # Show application logs
make shell          # Enter container shell

# Development & Testing
make lint           # Run ruff linter
make lint-fix       # Fix linting issues automatically
make test           # Run unit tests with coverage
make coverage       # Generate detailed coverage report
make behave         # Run BDD tests
make test-all       # Run all tests and linting

make clean          # Clean up containers and volumes
```

## Development Workflow && Testing Strategy

This project adopt a Behavior-Driven Development methodology for it's development 
and a Hexagonal Architecture as the software design pattern leveraging it's 
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
3. **Infrastructure**: Create adapters and web routes
4. **Tests**: Add unit tests and BDD scenarios
5. **Documentation**: Update this README

## Code Quality

### Linting with Ruff

- Fast Python linter and formatter
- Replaces flake8, isort, and black
- Configuration in `pyproject.toml`
- Run with: `make lint` or `make lint-fix`

---

## Nice-to-have features

- [ ] Add authentication and authorization
- [ ] Implement proper error handling
- [ ] Add API rate limiting

## Nice-to-have non-functional requirements

- [ ] Generate CHANGELOG.md automatically
- [ ] Define tree structured with a Hexagonal Architecture approach
- [ ] Linting rules (that executes previous a commit)
- [ ] Testing
  - [ ] unit
  - [ ] integration
  - [ ] e2e
  - [ ] performance
  - [ ] mutation
- [ ] Coverage rules
- [ ] Vulnerability checks for dependencies
- [ ] Versioning (SemVer)
- [ ] Environment configuration && deployment (Set up CI/CD pipeline)
- [ ] Monitoring && observability
  - [ ] Structured logging 
  - [ ] Monitoring
  - [ ] Healthcheck
- [ ] Security
  - [ ] Configure SAST
  - [ ] Container scanning
  - [ ] Define a secrets management policy
- [ ] Technical documentation
