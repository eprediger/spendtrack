# Spendtrack backend

Go REST API. Current capability: health check.

## Stack

- Go 1.27 (standard library `net/http`)
- PostgreSQL 16 via Docker Compose
- `godog` (BDD) and Go tests
- `golangci-lint`

## Endpoints

- `GET /api/v1/health` — `200 {"status":"ok"}`

## Local development

```bash
make db       # start postgres
make dev      # run the API with hot-reload (requires `air`)
make test     # unit tests
make bdd      # BDD suite against postgres
make lint     # golangci-lint
make ci       # lint + full suite (what the pre-commit hook runs)
make run      # build the image and run the full stack
```

## Configuration

See `.env.example`: `PORT` (default 8000).
