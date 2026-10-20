# S4b — Record Store Discounts: Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `superpowers:subagent-driven-development` or `superpowers:executing-plans`, as selected after review. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Record and review the complete set of item-specific discounts already applied on a receipt, while rejecting any set that does not reconcile exactly with the expense charged total.

**Architecture:** Keep discount values inside the `Expense` aggregate and each targeted `PurchaseItem`. The application loads the expense, asks the domain to validate a complete replacement set, then asks the Postgres repository to replace the persisted set in one transaction. The HTTP adapter parses decimal strings and maps domain errors to existing problem responses.

**Tech Stack:** Go, `net/http`, `pgx/v5`, `shopspring/decimal`, PostgreSQL migrations, Godog, testify.

**Spec:** [`spec.md`](spec.md). Scenarios: [`record_store_discounts.feature`](../../../backend/features/expenses/record_store_discounts.feature).

## Global Constraints

* Record receipt-applied item discounts, not promotion rules, future refunds, expense-wide discounts, or merchant identity.
* The full expense discount set is submitted together. A valid replacement overwrites the old set; an invalid replacement leaves it unchanged.
* Item discounts are positive decimal amounts, in the expense currency, and keep their entered decimal scale.
* Multiple discounts can target one item; their sum cannot exceed that item's positive line total.
* Exact reconciliation is `sum(item line totals) - sum(item discounts) = charged total`.
* A reduction is the concrete amount; a promotion is the rule. Do not add a `"store"` actor string or a promotion reference to `PurchaseItem`.
* `.feature` is the behavior source of truth. Implement one acceptance scenario at a time: activate only that scenario, watch it fail, implement the minimum, pass it, then continue.
* Run only one database-backed suite at a time. Use `make -C backend bdd`; finish with `make -C backend ci`.
* Use fixed safe HTTP problem responses. Never expose SQL or runtime errors.

## File Structure

* `backend/internal/domain/purchase_item.go`, `expense.go`: item discount values and complete-set validation.
* `backend/internal/application/expense_service.go`, `ports.go`: replace-discounts use case and repository port.
* `backend/internal/adapters/http/expense_handler.go`, `router.go`: PUT request parsing, response serialization, and route.
* `backend/internal/adapters/postgres/expense_repository.go` and `migrations/0003_expense_item_discounts.{up,down}.sql`: ordered persistence and atomic replacement.
* `backend/features/expenses/record_store_discounts.feature`, `steps_test.go`: executable acceptance scenarios and bindings.
* Existing domain, HTTP, application, and repository test files: focused validation/round-trip coverage.

## Review Focus

* Missing or empty discounts on an expense whose items do not reconcile: return 422 without altering existing discounts.
* An item position outside the expense's items: return 422 and do not partially persist the set.
* A zero or negative discount amount, blank description, or malformed decimal: reject with the correct validation field/code.
* Multiple discounts whose total exceeds the target item's positive line total, including a zero/negative return item: reject even if the overall expense would otherwise reconcile.
* Replacing an existing set with a mismatched set: reject and preserve the previous set; repeating a valid identical set does not duplicate entries.

---

### Task 1: Record one item discount

**Files:**
- Modify: `backend/features/expenses/record_store_discounts.feature`
- Modify: `backend/features/steps_test.go`
- Modify: `backend/internal/domain/purchase_item.go`, `expense.go`, `expense_test.go`
 - Modify: `backend/internal/domain/purchase_item.go` to add `ID uuid.UUID`, `expense.go`, `expense_test.go`
- Modify: `backend/internal/application/expense_service.go`, `ports.go`
- Modify: `backend/internal/adapters/http/expense_handler.go`, `router.go`, `expense_handler_test.go`
- Create: `backend/internal/adapters/postgres/migrations/0003_expense_item_discounts.up.sql`
- Create: `backend/internal/adapters/postgres/migrations/0003_expense_item_discounts.down.sql`
- Modify: `backend/internal/adapters/postgres/expense_repository.go`
 - Create: `backend/internal/adapters/postgres/migrations/0003_add_item_id_to_expense_items.up.sql`
 - Create: `backend/internal/adapters/postgres/migrations/0003_add_item_id_to_expense_items.down.sql`
 - Create: `backend/internal/adapters/postgres/migrations/0004_expense_item_discounts.up.sql`
 - Create: `backend/internal/adapters/postgres/migrations/0004_expense_item_discounts.down.sql`
 - Modify: `backend/internal/adapters/postgres/expense_repository.go`

**Interfaces:**
- `domain.ItemDiscountInput { ItemPosition *int; Description string; Amount *decimal.Decimal }`
 - `domain.ItemDiscountInput { ItemID *uuid.UUID; ItemPosition *int; Description string; Amount *decimal.Decimal }`
- `domain.ItemDiscount { Description string; Amount decimal.Decimal }`
- `func (e *Expense) ReplaceItemDiscounts(inputs []ItemDiscountInput) *ValidationError`; Task 1 validates field values and targets, Task 3 adds exact total reconciliation, and Task 6 adds the per-item cap.
- Add `Discounts []ItemDiscount` to `PurchaseItem`.
- Add `SaveItemDiscounts(ctx context.Context, e *domain.Expense) error` to `application.ExpenseRepository`, replacing the full persisted set for that expense.
- Add `ExpenseService.ReplaceItemDiscounts(ctx context.Context, id uuid.UUID, inputs []domain.ItemDiscountInput) (*domain.Expense, error)`.
- HTTP route: `PUT /api/v1/expenses/{id}/item-discounts`.
 - Add `SaveItemDiscounts(ctx context.Context, e *domain.Expense) error` to `application.ExpenseRepository`, replacing the full persisted set for that expense.
 - Add `ExpenseService.ReplaceItemDiscounts(ctx context.Context, id uuid.UUID, inputs []domain.ItemDiscountInput) (*domain.Expense, error)`.
 - Modify `backend/internal/adapters/http/createRequest` item shape to accept optional per-item `discounts` array during expense registration.
 - HTTP route: `PUT /api/v1/expenses/{id}/item-discounts`.

- [ ] **Step 1: Remove `@todo` from only `Scenario: Record a store discount on one item` and implement its step bindings using the receipt values: charged total 59,690.00 ARS, Coca Cola line total 12,750.00 ARS, other item totals 48,640.00 ARS, and discount 1,700.00 ARS.**
- [ ] **Step 2: Run `make -C backend bdd`; verify that scenario fails because no item-discount route/behavior exists.**
- [ ] **Step 3: Add failing domain tests for a valid 5.00 discount, required description/amount, positive amount, and an item position that does not exist. Verify with `cd backend && go test ./internal/domain -run 'TestExpenseReplaceItemDiscounts'`.**
- [ ] **Step 4: Implement the domain types and replacement method. Build and validate a candidate set before assigning it to the expense.**
- [ ] **Step 5: Add the migration and repository replacement method. Use one transaction to delete the prior set and insert the ordered candidate; rollback on any error.**
- [ ] **Step 6: Add the application operation, PUT handler, route, and nested `discounts` response serialization. Parse request decimal strings and merge parse errors with domain errors; always serialize an empty `discounts` array when an item has none.**
- [ ] **Step 7: Run `cd backend && go test ./internal/... ./cmd/...` and `make -C backend bdd`; verify the activated scenario passes.**
- [ ] **Step 8: Commit** with `feat: record item discounts for an expense`.

### Task 2: Review the discounts on an expense

**Files:**
- Modify: `backend/features/expenses/record_store_discounts.feature`, `backend/features/steps_test.go`
- Modify: `backend/internal/adapters/postgres/expense_repository.go`
- Test: `backend/internal/adapters/postgres/expense_repository_test.go`, `backend/internal/adapters/http/expense_handler_test.go`

- [ ] **Step 1: Remove `@todo` from only `Scenario: Review the store discounts of an expense`.**
- [ ] **Step 2: Run the BDD scenario; verify it fails because `GET` does not load and show discounts.**
- [ ] **Step 3: Add a repository test that loads discounts for their item in stored order and preserves decimal scale.**
- [ ] **Step 4: Load discounts ordered by item and discount position; attach them to each `PurchaseItem`.**
- [ ] **Step 5: Verify HTTP output always includes each item's `discounts` array (empty if none) and includes `amount` and `description` for a recorded discount.**
- [ ] **Step 6: Run focused Go tests and the BDD suite; verify the review scenario passes.**
- [ ] **Step 7: Commit** with `feat: review expense item discounts`.

### Task 3: Reject an unreconciled complete set

**Files:**
- Modify: `backend/features/expenses/record_store_discounts.feature`, `backend/features/steps_test.go`
- Modify: `backend/internal/domain/expense.go`, `expense_test.go`
- Modify: `backend/internal/adapters/http/expense_handler_test.go`

- [ ] **Step 1: Remove `@todo` from only `Scenario: Reject an item discount when the itemized expense does not reconcile`.**
- [ ] **Step 2: Run the scenario; verify that a 3.00 ARS discount against 100.00 ARS item total and 95.00 ARS charged total is incorrectly accepted before the fix.**
- [ ] **Step 3: Add a domain test that rejects a set unless total item line totals minus all discounts equals the charged total.**
- [ ] **Step 4: Implement the exact comparison with `decimal.Decimal`; report the mismatch as a validation error on `discounts`.**
- [ ] **Step 5: Verify the handler returns 422 and the repository is not called when validation fails.**
- [ ] **Step 6: Run focused Go tests and the BDD suite; verify the mismatch scenario passes.**
- [ ] **Step 7: Commit** with `test: require item discounts to reconcile`.

### Task 4: Prove multiple discounts on one item

**Files:**
- Modify: `backend/features/expenses/record_store_discounts.feature`, `backend/features/steps_test.go`
- Test: `backend/internal/domain/expense_test.go`, `backend/internal/adapters/postgres/expense_repository_test.go`

- [ ] **Step 1: Remove `@todo` from only `Scenario: Record multiple discounts on one item together`.**
- [ ] **Step 2: Add the BDD step binding to submit both 3.00 and 2.00 ARS entries in one complete-set request.**
- [ ] **Step 3: Add domain and repository tests asserting both ordered discounts are attached to the target item and `100.00 - 3.00 - 2.00 = 95.00` is accepted.**
- [ ] **Step 4: Run focused tests and BDD; verify the scenario passes. If it already passes with Task 1's list model, keep this task test-only rather than changing production code.**
- [ ] **Step 5: Commit** with `test: cover multiple discounts on one item`.

### Task 5: Replace the complete set atomically

**Files:**
- Modify: `backend/features/expenses/record_store_discounts.feature`, `backend/features/steps_test.go`
- Test: `backend/internal/application/expense_service_test.go`, `backend/internal/adapters/postgres/expense_repository_test.go`

- [ ] **Step 1: Remove `@todo` from only `Scenario: Replace the complete item discount set`.**
- [ ] **Step 2: Run the scenario; verify it fails unless an invalid replacement is rejected and the prior set remains.**
- [ ] **Step 3: Add a domain/application test that invalid input preserves the existing set and does not call the repository, plus a repository test that a failed transaction rolls back deletion of the old set.**
- [ ] **Step 4: Ensure candidate domain validation happens before repository writes and repository delete/insert is transactional.**
- [ ] **Step 5: Run focused tests and BDD; verify the scenario passes and previous discounts remain after rejected replacement.**
- [ ] **Step 6: Commit** with `feat: replace item discount sets atomically`.

### Task 6: Prove valid replacement does not append

**Files:**
- Modify: `backend/features/expenses/record_store_discounts.feature`, `backend/features/steps_test.go`
- Test: `backend/internal/adapters/http/expense_handler_test.go`, `backend/internal/adapters/postgres/expense_repository_test.go`

- [ ] **Step 1: Remove `@todo` from only `Scenario: Replace an item discount set without duplicating entries`.**
- [ ] **Step 2: Add the BDD steps to submit `[3.00, 2.00]`, replace with `[4.00, 1.00]`, and submit that same replacement again.**
- [ ] **Step 3: Add assertions that the resulting item has exactly two discounts, in order and with the replacement values.**
- [ ] **Step 4: Run focused tests and BDD; verify old values are removed and the repeated request does not duplicate rows.**
- [ ] **Step 5: Commit** with `test: prove item discount replacement is idempotent`.

### Task 7: Enforce each item's discount limit

**Files:**
- Modify: `backend/features/expenses/record_store_discounts.feature`, `backend/features/steps_test.go`
- Test: `backend/internal/domain/expense_test.go`

- [ ] **Step 1: Remove `@todo` from only `Scenario: Reject discounts that exceed the item's line total`.**
- [ ] **Step 2: Run the scenario; verify that 110.00 ARS discounts on a 100.00 ARS item are rejected, even though the expense-wide total reconciles.**
- [ ] **Step 3: Add domain tests rejecting cumulative discounts above line total and any discount on zero/negative items; allow a discount equal to a positive line total.**
- [ ] **Step 4: Implement per-item cumulative validation before expense-wide reconciliation.**
- [ ] **Step 5: Run focused tests and BDD; verify the scenario passes.**
- [ ] **Step 6: Commit** with `feat: cap item discounts at line total`.

### Task 8: Validate the S4b slice

- [ ] **Step 1: Run `make -C backend test`.**
- [ ] **Step 2: Run `make -C backend bdd` by itself against the isolated test database.**
- [ ] **Step 3: Run `make -C backend ci`.**
- [ ] **Step 4: Confirm all S4b scenarios pass without `@todo`, no other slice's scenarios were activated, and the glossary/ADR/decisions links still agree.**
