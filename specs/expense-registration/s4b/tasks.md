# S4b — Record Store Discounts: Implementation Plan

> **For agentic workers:** implement one acceptance scenario at a time: activate only that scenario, watch it fail, implement the minimum, pass it, then continue. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Record item discounts nested under `items[]` at registration, show them on `GET`, and require every itemized expense to reconcile exactly.

**Architecture:** `NewExpense` validates items (with their discounts) and reconciliation. The HTTP adapter parses decimal strings and merges parse errors into the problem response. The Postgres repository saves and loads discounts inside the existing expense aggregate; there is no port or service change.

**Tech Stack:** Go, `net/http`, `pgx/v5`, `shopspring/decimal`, PostgreSQL migrations, Godog, testify.

**Spec:** [`spec.md`](spec.md). Scenarios: [`record_store_discounts.feature`](../../../backend/features/expenses/record_store_discounts.feature).

## Global Constraints

* Record receipt-applied item discounts, not promotion rules, later refunds, expense-wide discounts, or merchant identity.
* Discounts are sent at registration inside the item they apply to; no item reference is needed.
* Item discounts are positive decimal amounts, in the expense currency, and keep their entered decimal scale.
* Multiple discounts can target one item; their sum cannot exceed that item's positive line total.
* Exact reconciliation is `sum(item line totals) - sum(item discounts) = charged total`, checked only when no other field error exists.
* Use fixed safe HTTP problem responses. Never expose SQL or runtime errors.
* Run only one database-backed suite at a time. Use `make -C backend bdd`; finish with `make -C backend ci`.

## File Structure

* `backend/internal/domain/purchase_item.go`, `expense.go`: `ItemDiscountInput`, `ItemDiscount`, per-discount validation, per-item cap, expense reconciliation.
* `backend/internal/adapters/http/expense_handler.go`: `discounts` request/response shapes, decimal parsing, parse-error merge.
* `backend/internal/adapters/postgres/expense_repository.go` and `migrations/0003_expense_item_discounts.{up,down}.sql`: ordered persistence inside the save transaction and the load query.
* `backend/features/expenses/record_store_discounts.feature`, `steps_test.go`: executable acceptance scenarios and bindings.
* Existing domain, HTTP, and repository test files: focused validation and round-trip coverage.

## Tasks

### Task 1: Reject an itemized expense whose items do not reconcile

**Files:**
- Modify: `backend/features/expenses/record_store_discounts.feature`, `backend/features/steps_test.go`
- Modify: `backend/internal/domain/expense.go`, `expense_test.go`
- Modify: `backend/internal/adapters/http/expense_handler_test.go` fixtures that use items

- [ ] Activate `Scenario: Reject an itemized expense whose items do not reconcile` and watch it fail.
- [ ] Add `NewExpense` reconciliation: with no field errors and at least one item, `sum(line_total)` must equal `amount` via `decimal.Equal`; mismatch reports `items`/`not_reconciled`.
- [ ] Fix S4a fixtures that no longer reconcile (single-item helpers send `amount` = the item line total when it parses; the two-item helper sends the items sum; rejection outlines keep `14000.00`).
- [ ] Focused `go test ./internal/...` plus `make -C backend bdd`; commit.

### Task 2: Register an expense with a store discount on one item

**Files:**
- Modify: `backend/internal/domain/purchase_item.go`, `expense.go`, `purchase_item_test.go`, `expense_test.go`
- Modify: `backend/internal/adapters/http/expense_handler.go`, `expense_handler_test.go`
- Modify: `backend/features/expenses/record_store_discounts.feature`, `backend/features/steps_test.go`

- [ ] Activate `Scenario: Register an expense with a store discount on one item` and watch it fail.
- [ ] Add `ItemDiscountInput`, `ItemDiscount`, `PurchaseItemInput.Discounts`, `PurchaseItem.Discounts`; `NewPurchaseItem` gains a `discounts` parameter (update all callers).
- [ ] Per-discount validation: `description` required (trimmed), `amount` required, `not_positive` when ≤ 0.
- [ ] Parse `items[].discounts[].amount` in the adapter; report `items[i].discounts[j].amount`/`invalid` through `withParseErrors`. Subtract discounts in the reconciliation.
- [ ] Response `items[].discounts` is always an array with scale-preserved `amount` and `description`.
- [ ] Focused tests plus BDD; commit.

### Task 3: Review the store discounts of an expense

**Files:**
- Create: `backend/internal/adapters/postgres/migrations/0003_expense_item_discounts.{up,down}.sql`
- Modify: `backend/internal/adapters/postgres/expense_repository.go`, `expense_repository_test.go`
- Modify: `backend/features/steps_test.go` truncate, `AGENTS.md`, `.github/copilot-instructions.md` (one sentence only)

- [ ] Activate `Scenario: Review the store discounts of an expense` and watch it fail.
- [ ] Add the migration per the spec schema.
- [ ] `Save` inserts discounts in the same transaction; `FindByID` loads them ordered by `item_position, position` into `PurchaseItemInput.Discounts` before `NewExpense`.
- [ ] Repository round-trip test: order and scale (`"1700.00"`).
- [ ] Update the BDD truncate to `TRUNCATE expense_item_discounts, expense_items, expenses` and the matching sentences in `AGENTS.md` and `.github/copilot-instructions.md`.
- [ ] Focused tests plus BDD; commit.

### Task 4: Reject item discounts that do not reconcile

- [ ] Activate `Scenario: Reject item discounts that do not reconcile`. If it already passes after Task 2, keep this task test-only; commit.

### Task 5: Register multiple discounts on one item

- [ ] Activate `Scenario: Register multiple discounts on one item`. If it already passes after Task 2, keep this task test-only; commit.

### Task 6: Reject discounts that exceed the item's line total

**Files:**
- Modify: `backend/features/expenses/record_store_discounts.feature`, `backend/features/steps_test.go`
- Modify: `backend/internal/domain/purchase_item.go`, `purchase_item_test.go`

- [ ] Activate `Scenario: Reject discounts that exceed the item's line total` and watch it fail.
- [ ] When all discount entries on an item are valid and their sum exceeds `max(line_total, 0)`, report `discounts`/`exceeds_line_total` (a sum equal to the line total is allowed; any discount on a zero or negative item is rejected).
- [ ] Focused tests plus BDD; commit.

### Task 7: Reject invalid item discount data

- [ ] Activate `Scenario Outline: Invalid item discount data is rejected`. Cover `amount` 0/−1 → `not_positive`, `abc` → `invalid` (adapter), missing → `required`, and missing `description` → `required`. If it already passes, keep this task test-only; commit.

### Task 8: Validate the S4b slice

- [ ] Run `make -C backend test`.
- [ ] Run `make -C backend bdd` by itself against the isolated test database.
- [ ] Run `make -C backend ci`.
- [ ] Confirm all active S4b scenarios pass, only the two PUT replacement scenarios stay `@todo`, and all S4a scenarios still pass.
