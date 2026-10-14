# S4a — Break the Expense Down Into Items: Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** An expense can carry optional items (description, quantity, unit price, line total), and reviewing the expense shows them.

**Architecture:** Hexagonal, as in S1. `domain.Item` is built only by `domain.NewItem`, called by `NewExpense` for registration and by the Postgres repository for loading. The HTTP adapter parses the decimals and merges its parse errors with the domain's. Items live in a new table, saved with the expense in one transaction.

**Tech Stack:** Go, `net/http`, `pgx/v5`, `shopspring/decimal`, `golang-migrate`, `godog`, `testify`.

**Spec:** [`spec.md`](spec.md) (approved). Scenarios: [`itemize_expense.feature`](../../../backend/features/expenses/itemize_expense.feature).

## Global Constraints

* Decimals are strings in requests and responses, and keep their scale (`"0.750"` stays `"0.750"`).
* Item errors are named `items[<position>].<field>`, position from 0, fields `description`, `quantity`, `unit_price`, `line_total`. Codes: `required`, `invalid`, `not_positive`.
* `NewItem` is the only way to build an `Item`, for registration and for loading.
* `quantity` > 0. `unit_price` and `line_total` take any sign, zero included. No check between item fields, and none between items and the amount.
* Item amounts are `decimal.Decimal`, not `Money`. An item has no currency.
* Internal errors answer the generic `500`. The details go to the log only.
* Run only one test suite at a time (shared Postgres, see `AGENTS.md`). Verification is `make -C backend ci`.
* Commits use Conventional Commits. The pre-commit hook runs the full CI, so every commit must be green.
* Do not create `docs/superpowers/`.

Shorthand used below: `BDD` = `DATABASE_URL="postgres://postgres:postgres@localhost:5432/spendtrack?sslmode=disable" go -C backend test -count=1 ./features/...`

## Review Focus

* `items` omitted or `null`: the expense has no items, and the response has `"items": []`, never `null`. Pinned in Task 3 (`the expense has no items` asserts an empty array, and `TestParseItems` / `TestToResponseItems`).
* `items: [null]`: one item at position 0 with all four fields `required`. Pinned in Task 1 and Task 4.
* A description of only spaces is `required`. Pinned in Task 1.
* Two invalid items in one request: both positions are reported, together with the expense's own errors. Pinned in Task 2.
* `items` that is not an array (`"items": "x"`) answers `400` malformed JSON. Pinned in Task 4.

---

### Task 0: Commit the approved design

**Files:**
- Modify (already edited, uncommitted): `backend/features/expenses/itemize_expense.feature`, `product/epics/financial-management/features/expense-registration/decisions.md`
- Create (already written, uncommitted): `specs/expense-registration/s4a/spec.md`, `specs/expense-registration/s4a/tasks.md`

- [ ] **Step 1: Commit**

```bash
git add backend/features/expenses/itemize_expense.feature product/epics/financial-management/features/expense-registration/decisions.md specs/expense-registration/s4a
git commit -m "docs: add S4a items specification, agreements and task plan"
```

Expected: the hook passes (nothing in the Go code changed).

### Task 1: `NewItem`

**Files:**
- Create: `backend/internal/domain/item.go`
- Test: `backend/internal/domain/item_test.go`

**Interfaces:**
- Produces:
  - `type ItemInput struct { Description string; Quantity, UnitPrice, LineTotal *decimal.Decimal }` (a nil pointer means the field was missing)
  - `type Item struct { Description string; Quantity, UnitPrice, LineTotal decimal.Decimal }`
  - `func NewItem(in ItemInput) (Item, []FieldError)`. Field names in the errors: `description`, `quantity`, `unit_price`, `line_total` (no position). The description is trimmed in the result.

- [ ] **Step 1: Write the failing tests** in `item_test.go` (use the `dec` and `ptr` helpers of `expense_test.go`)
  - `TestNewItemValid`: input `"  Candies "`, quantity `5`, unit price `150.00`, line total `750.00` gives description `"Candies"` and equal amounts, no errors.
  - `TestNewItemSignedAmounts`: quantity `3`, unit price `-3001.00`, line total `-9003.00` is valid. Unit price `0` and line total `0` are valid.
  - `TestNewItemValidation` (table, compares `field -> code`): description `"   "` gives `description: required`; quantity nil gives `quantity: required`; quantity `0` and `-1` give `quantity: not_positive`; unit price nil gives `unit_price: required`; line total nil gives `line_total: required`; an empty `ItemInput{}` gives all four `required` together.

- [ ] **Step 2: Run to verify they fail**

Run: `go -C backend test ./internal/domain -run TestNewItem`
Expected: build failure, `undefined: NewItem` and `undefined: ItemInput`.

- [ ] **Step 3: Implement `NewItem(in ItemInput) (Item, []FieldError)` in `item.go`**

Follow the style of `NewMoney`: collect every error, return `Item{}` when there is any.

- [ ] **Step 4: Run to verify they pass**

Run: `go -C backend test -count=1 ./internal/domain`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/domain/item.go backend/internal/domain/item_test.go
git commit -m "feat: add the Item value built only through NewItem"
```

### Task 2: `NewExpense` takes items

**Files:**
- Modify: `backend/internal/domain/expense.go`, `backend/internal/domain/expense_test.go`, `backend/internal/application/expense_service.go`

**Interfaces:**
- Consumes: `ItemInput`, `Item`, `NewItem` (Task 1).
- Produces:
  - `Expense` gets `Items []Item` (nil when there are none).
  - `func NewExpense(occurredAt time.Time, amount *decimal.Decimal, currency *string, description string, items []ItemInput, now time.Time) (*Expense, *ValidationError)`
  - `application.RegisterInput` gets `Items []domain.ItemInput`. `Register` passes it on.

- [ ] **Step 1: Write the failing tests** in `expense_test.go`
  - Update every existing `NewExpense` call to pass `nil` for `items`.
  - `TestNewExpenseItems`: two valid items keep their order in `e.Items`; no items gives `len(e.Items) == 0`.
  - `TestNewExpenseItemErrors` (field to code): item 1 with quantity `0` plus an invalid currency gives `items[1].quantity: not_positive` and `currency: unsupported_currency` together. Item 0 with an empty description and item 1 with a nil line total give `items[0].description: required` and `items[1].line_total: required` together.

- [ ] **Step 2: Run to verify they fail**

Run: `go -C backend test ./internal/domain`
Expected: build failure on the `NewExpense` arguments.

- [ ] **Step 3: Implement.** `NewExpense` calls `NewItem` for each input, prefixes each error field with `items[<position>].`, and appends them to the expense's own errors. Pass `in.Items` through in `Register`. The handler needs no change yet: the new `RegisterInput` field stays nil until Task 3.

- [ ] **Step 4: Run to verify they pass**

Run: `go -C backend test -count=1 ./internal/...`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add backend/internal
git commit -m "feat: validate items when building an expense"
```

### Task 3: Items round trip through the API

**Files:**
- Create: `backend/internal/adapters/postgres/migrations/0002_expense_items.up.sql`, `0002_expense_items.down.sql`
- Modify: `backend/internal/adapters/postgres/expense_repository.go`, `backend/internal/adapters/http/expense_handler.go`, `backend/internal/adapters/http/expense_handler_test.go`, `backend/features/steps_test.go`, `backend/features/expenses/itemize_expense.feature`

**Interfaces:**
- Consumes: Task 2.
- Produces in the handler:
  - `type itemRequest struct { Description *string; Quantity, UnitPrice, LineTotal *string }` with JSON keys `description`, `quantity`, `unit_price`, `line_total`; `createRequest` gets `Items []itemRequest` (`items`).
  - `func parseItems(reqs []itemRequest) ([]domain.ItemInput, []domain.FieldError)`. Unparsable decimal text gives a nil pointer plus an `invalid` error named `items[i].<field>`. A missing field gives a nil pointer and no error.
  - `func withParseErrors(fields, parseErrs []domain.FieldError) []domain.FieldError`. It replaces domain errors that have the same field name. It generalises `withAmountError`, which is removed.
  - `type itemResponse struct { Description, Quantity, UnitPrice, LineTotal string }` with the same JSON keys; `expenseResponse` gets `Items []itemResponse` (`items`), never nil.
- Produces in the repository: `func numericText(d decimal.Decimal) string` (scale-preserving text, replaces the inline `StringFixed` in `Save`).

- [ ] **Step 1: Prepare the scenarios.** Move `@todo` from the Feature line to every Scenario and Scenario Outline of `itemize_expense.feature`, then remove it from the four scenarios of this task: register with items, register without items, review the items, and a new one:

```gherkin
  Scenario: Reviewing an expense whose stored item is invalid
    Given an expense with an item with an empty description is stored in the database
    When the Financial Manager reviews the expense
    Then the system indicates an internal error
```

- [ ] **Step 2: Add the step bindings** in `steps_test.go`:
  - `the Financial Manager registers an expense at "…" for … "…" described as "…" with items:` builds the payload from the table (columns `description`, `quantity`, `unit price`, `line total`, sent as strings).
  - `the expense is recorded` (status 201, keeps the id).
  - `the expense has (\d+) items` and `the expense has no items` (the response `items` must be an array, an empty one for the second).
  - `an expense with items was registered` (two fixed items, keeps the id).
  - `the expense shows each item with its description, quantity, unit price and line total` (same items, same order, same text as sent).
  - `an expense with an item with an empty description is stored in the database` (inserts the expense row, then an `expense_items` row with an empty description, directly through `pool`, and keeps the id).

- [ ] **Step 3: Write the failing handler tests** in `expense_handler_test.go`
  - `TestParseItems`: `nil` input gives no inputs and no errors; `quantity: "abc"` at position 1 gives nil `Quantity` and `{items[1].quantity, invalid, "must be a decimal number"}`; a missing `unit_price` gives a nil pointer and no error; a decimal keeps its scale.
  - `TestWithParseErrors`: replaces the domain's `items[1].quantity: required` by the parse error and leaves other fields; the input slice is not modified (replaces `TestWithAmountError`, keeping its amount case).
  - `TestToResponseItems`: an expense without items gives `Items` as an empty, non-nil slice; amounts are formatted with their scale (`"0.750"`).

- [ ] **Step 4: Run to verify they fail**

Run: `go -C backend test ./internal/adapters/http` then `BDD`
Expected: the handler tests fail to build (`undefined: parseItems`); the scenarios fail because the API ignores `items` (for example `the expense has 3 items` does not hold, and the stored-invalid-item review answers `200`).

- [ ] **Step 5: Implement**
  - Migration `0002_expense_items`: table `expense_items` with `expense_id UUID NOT NULL REFERENCES expenses(id) ON DELETE CASCADE`, `position INTEGER NOT NULL`, `description TEXT NOT NULL`, `quantity NUMERIC NOT NULL CHECK (quantity > 0)`, `unit_price NUMERIC NOT NULL`, `line_total NUMERIC NOT NULL`, `PRIMARY KEY (expense_id, position)`. The down file drops the table.
  - `Save`: one transaction, the expense first, then each item with its position, amounts as `numericText`.
  - `FindByID`: after the money check, read the items `ORDER BY position` and build each with `NewItem`. If one fails, return an error (the handler already answers `500` and logs it).
  - Handler: parse with `parseItems`, merge parse errors with `withParseErrors` (amount error included), format the response items with the same scale-preserving formatter as the amount.

- [ ] **Step 6: Run to verify they pass**

Run: `go -C backend test -count=1 ./internal/...` then `BDD`
Expected: PASS, and the four scenarios run (no `@todo`).

- [ ] **Step 7: Run the whole verification**

Run: `make -C backend ci`
Expected: 0 lint issues, all tests pass.

- [ ] **Step 8: Commit**

```bash
git add backend
git commit -m "feat: register and review expense items"
```

### Task 4: Item validation through the API

**Files:**
- Modify: `backend/features/expenses/itemize_expense.feature`, `backend/features/steps_test.go`, `backend/internal/adapters/http/expense_handler_test.go`

**Interfaces:** Consumes Task 3. Adds no production interface.

- [ ] **Step 1: Prepare.** Add this scenario to the feature, tagged `@todo`, then remove `@todo` from it and from the two outlines `A quantity that is not positive is rejected` and `Incomplete item information is rejected`:

```gherkin
  Scenario Outline: A non-numeric item amount is rejected
    When the Financial Manager registers an expense with an item whose "<field>" is "abc"
    Then the registration is rejected
    And the system indicates the item "<field>" is invalid

    Examples:
      | field      |
      | quantity   |
      | unit price |
      | line total |
```

- [ ] **Step 2: Add the step bindings:** `registers an expense with an item of quantity (\S+), unit price (\S+) and line total (\S+)`, the same with `an item "([^"]*)" of quantity …`, `registers an expense with an item without "([^"]*)"`, `registers an expense with an item whose "([^"]*)" is "([^"]*)"`, and `the system indicates the item "([^"]*)" is (required|invalid)` (an error on any position whose field ends with the field name, spaces turned into underscores; `required` needs code `required`, `invalid` accepts any code).

- [ ] **Step 3: Write the failing handler tests**
  - `TestParseItemsNullItem`: JSON `{"items":[null]}` decoded into `createRequest` and parsed gives one `ItemInput` with an empty description and three nil pointers.
  - `TestCreateItemsNotAnArray`: `NewExpenseHandler(nil).Create` with body `{"items":"x"}` answers `400` and the `malformed-json` problem (the handler returns before using the service).

- [ ] **Step 4: Run to verify the new tests behave**

Run: `go -C backend test ./internal/adapters/http` then `BDD`
Expected: the handler tests pass or fail only if Task 3 missed a case; the scenarios pass if Task 3 is complete. A failure means a gap in Task 3: fix it here, with the failing test first.

- [ ] **Step 5: Run the whole verification**

Run: `make -C backend ci`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add backend
git commit -m "test: cover item validation through the API"
```

### Task 5: Sign, total and scale scenarios

**Files:**
- Modify: `backend/features/expenses/itemize_expense.feature`, `backend/features/steps_test.go`

**Interfaces:** Consumes Task 3 and Task 4.

- [ ] **Step 1:** Remove `@todo` from the scenarios `The printed line total differs from quantity times unit price`, `A returned deposit is a negative item` and `A quantity can be a fraction`.

- [ ] **Step 2: Add the step bindings** `the item keeps the line total (\S+) as printed` and `the item shows quantity (\S+) and line total (\S+)` (exact text of the first item in the last response).

- [ ] **Step 3: Run**

Run: `BDD`
Expected: PASS with no new production code. These scenarios pin behavior the earlier tasks already give (no cross-check of line total, negative price accepted, scale kept). Their items also do not add up to the expense amount, which pins the rule that items are not checked against it. If one fails, a forbidden check or a scale loss exists: fix the code, not the scenario.

- [ ] **Step 4: Commit**

```bash
git add backend/features
git commit -m "test: pin item sign, printed total and fractional quantity"
```

### Task 6: Close the slice

**Files:**
- Modify: `backend/features/expenses/itemize_expense.feature`, `product/glossary.md`, `specs/expense-registration/s4a/spec.md`

- [ ] **Step 1:** Make sure no `@todo` is left in `itemize_expense.feature`. The first line is `@expense-registration @s4a`.

- [ ] **Step 2:** In `product/glossary.md`, change the status of `Item` to `validated by S4a (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))`. Leave `Purchase` and `Product` as they are (S4a does not activate them).

- [ ] **Step 3:** In `spec.md`, change the status to `implemented`.

- [ ] **Step 4: Run the whole verification**

Run: `make -C backend ci`
Expected: 0 lint issues, all scenarios run (none skipped as `@todo` in `itemize_expense.feature`), all tests pass.

- [ ] **Step 5: Commit**

```bash
git add backend/features product/glossary.md specs/expense-registration/s4a/spec.md
git commit -m "docs: mark S4a as implemented and Item as validated"
```

## Left Out On Purpose

* A limit on the number of items, and a limit on the size of decimal text. A decimal such as `1e999999999` can use a lot of memory when it is formatted or stored. This already holds for `amount`, so it is a separate change.
* Editing or deleting items, tax rate, barcode, unit of measure and the photo flow (see the spec).
