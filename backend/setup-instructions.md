# Setup Instructions

Go backend + PostgreSQL via Docker Compose.

## Requirements

- Go 1.27+
- Docker + Docker Compose
- `golangci-lint` (for `make lint` / `make ci`)

## Run the stack

```bash
make run
```

Builds the API image and starts it with Postgres. The API listens on `http://localhost:8000` and applies migrations on startup.

Endpoints: `POST /api/v1/expenses`, `GET /api/v1/expenses/{id}`, `GET /api/v1/health`.

## Local development

```bash
make db       # start Postgres only
make test     # unit tests (no DB needed)
make bdd      # BDD suite against local Postgres
make lint     # golangci-lint
make ci       # lint + full suite (what pre-commit runs)
```

## Configuration

See `.env.example`: `PORT` (default 8000), `DATABASE_URL` (required by the app; compose sets it automatically).

## Git hooks

Run `backend/setup-hooks.sh` from the repository root — it links `.hooks/` into `.git/hooks/`. `pre-commit` runs `make ci`.
