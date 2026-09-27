# S1 — Register a Basic Expense: Specification

**Status:** approved

## Traceability

* Story (US + AC + BDD): [`backend/features/expenses/register_basic_expense.feature`](../../../backend/features/expenses/register_basic_expense.feature) (`@s1`)
* Agreements: [`product/epics/financial-management/features/expense-registration/decisions.md`](../../../product/epics/financial-management/features/expense-registration/decisions.md) — S1 Required Attributes and Validation Rules
* Stack: [`adr/0001-backend-language.md`](../../../adr/0001-backend-language.md) (Go), [`adr/0002-api-protocol.md`](../../../adr/0002-api-protocol.md) (REST + JSON), Postgres for persistence

## Scope

Register an expense as an economic event with occurrence date and time, amount, currency, and description, and review the registered expense. No payment mechanism, document, purchase detail, classification, or installments — those are later slices.

## API Contract

### `POST /api/v1/expenses`

Request body:

```json
{
  "occurred_at": "2026-09-20T14:30:00-03:00",
  "amount": "45.90",
  "currency": "USD",
  "description": "Groceries"
}
```

`amount` is transmitted as a string and interpreted as an exact decimal (no float). `currency` is an ISO 4217 alphabetic code (3 letters). `occurred_at` is an RFC 3339 datetime with explicit UTC offset — the client sends the moment as observed locally; the offset preserves the wall-clock time (needed for time-windowed promotion analysis).

* `201 Created` — returns the expense resource:

```json
{
  "id": "expense-uuid",
  "occurred_at": "2026-09-20T14:30:00-03:00",
  "amount": "45.90",
  "currency": "USD",
  "description": "Groceries"
}
```

* `422 Unprocessable Entity` — RFC 9457 `application/problem+json` body with a per-field error list when required information is missing or invalid:

```json
{
  "type": "https://spendtrack.dev/problems/validation-failed",
  "title": "Validation failed",
  "status": 422,
  "errors": [
    { "field": "amount", "code": "not_positive", "detail": "must be greater than 0" },
    { "field": "currency", "code": "invalid_iso4217", "detail": "unknown currency code" }
  ]
}
```

### `GET /api/v1/expenses/{id}`

* `200 OK` — returns the expense resource (same shape as POST response).
* `404 Not Found` — `application/problem+json` when the id does not exist:

```json
{
  "type": "https://spendtrack.dev/problems/not-found",
  "title": "Expense not found",
  "status": 404
}
```

## Validation Rules (contract level)

| Field            | Rule                                                        |
| ---------------- | ----------------------------------------------------------- |
| `occurred_at`    | required; RFC 3339 datetime with explicit offset; instant must not be in the future |
| `amount`         | required; decimal; must be > 0                              |
| `currency`       | required; valid ISO 4217 alphabetic code                    |
| `description`    | required; non-empty after trimming whitespace               |

Notes:

* "Not in the future" compares instants (the offset makes it exact — no date-boundary ambiguity).
* The original UTC offset is persisted alongside the instant so the local wall-clock time remains recoverable. The IANA zone name (`America/Argentina/Buenos_Aires`) is not stored — only needed for future recurring rules, revisit then.
* ISO 4217 validation: the code must be a known currency code (validated against the ISO 4217 list, not merely format).

## Domain Model (minimal, for this slice)

* `Expense` — entity holding `id`, `occurred_at` (instant + local offset), `amount`, `currency`, `description`; owns its invariants (positive amount, non-future instant, non-empty description).
* `Financial Context` — S1 uses a single implicit context; the expense record carries no context reference yet. Explicit context modeling arrives with the slice that requires multiple contexts or shared management.
* `Money` — value object (`amount` + `currency`) so both always travel together.
* `ExpenseRepository` — driven port; Postgres adapter (e.g., `sqlc`/`pgx` — decided in the technical plan).

## Persistence

* Postgres (per technical decision), schema managed with migrations (tooling decided in the technical plan).
* `expenses` table: `id` UUID PK, `occurred_at` TIMESTAMPTZ, `occurred_offset_minutes` SMALLINT, `amount` NUMERIC, `currency` CHAR(3), `description` TEXT, `created_at` TIMESTAMPTZ.

## Error Semantics

Validation failures return RFC 9457 problem details (`application/problem+json`) with `type`, `title`, `status`, and an `errors` array of `{field, code, detail}` entries — covering both missing and invalid fields, per the story's scenarios.

## Explicitly Out of This Slice

Payment mechanisms and settlements, supporting documents, purchase detail, classification structure, people and economic allocation, installments/payment commitments, multiple financial contexts, authentication/authorization.

## Open Points (resolved at plan time)

* HTTP error codes for domain rejection (proposed: 422 for validation).
* Timezone semantics resolved: instant + local offset stored. IANA zone name not stored — revisit only if recurring/future rules need it.
