# S4a — Break the Expense Down Into Items: Specification

**Status:** approved

## Traceability

* Story (US + AC + BDD): [`backend/features/expenses/itemize_expense.feature`](../../../backend/features/expenses/itemize_expense.feature) (`@s4a`)
* Agreements: [`decisions.md`](../../../product/epics/financial-management/features/expense-registration/decisions.md), section S4a
* Builds on: [S1 specification](../s1/spec.md) and [S1 technical plan](../s1/plan.md)
* Examples: [gas station receipt](../../../product/examples/gas-station-receipt.md), [supermarket receipt](../../../product/examples/supermarket-receipt.md), [bakery receipt](../../../product/examples/bakery-receipt.md)

## Why

A Financial Manager wants to see what an expense was made of, not only its total. The receipts have the lines (description, quantity, unit price, line total). This slice stores and shows them. Items also prepare the later slices: discounts per item (S4b), effective cost (S4c) and price analysis.

## Scope

Register an expense with optional items, and review the items of a registered expense. The client sends the items already in structured form. Reading a photo of a receipt is not part of this slice and does not change it.

## Decided by the product owner

* An item has a description, a quantity, a unit price and a line total. Items are optional.
* The quantity is a positive decimal number (greater than 0).
* The unit price and the line total keep the sign the receipt prints. A returned deposit is a positive quantity with a negative unit price.
* The line total is stored as printed. It is not checked against quantity times unit price.
* Items are not checked against the amount of the expense. The product owner prefers that they add up exactly, which belongs to S4b.
* An item is built in one way only: through its constructor. Registration and loading both use it, so an invalid item can never exist in the model. A stored item that fails the check is not served.
* All items are in the currency of the expense. Items in different currencies are not supported.
* Editing or deleting items, tax rate, barcode, unit of measure, product identification and the photo flow are outside this slice.

## API Contract

### `POST /api/v1/expenses`

The S1 request gets one new optional key, `items`. Decimals are strings, like `amount`.

```json
{
  "occurred_at": "2026-10-05T16:53:34-03:00",
  "amount": "14000.00",
  "currency": "ARS",
  "description": "Gas station store",
  "items": [
    { "description": "Candies", "quantity": "5", "unit_price": "150.00", "line_total": "750.00" },
    { "description": "Bottled water", "quantity": "1", "unit_price": "3400.00", "line_total": "3400.00" }
  ]
}
```

* `items` missing or `[]` means no items. The S1 behavior is unchanged.
* `201 Created` returns the S1 resource plus `items` (see below).
* `422 Unprocessable Entity` uses the S1 problem body. Item errors are added to the same `errors` list, with the position of the item (starting at 0) in the field name, for example `items[1].quantity`.

### `GET /api/v1/expenses/{id}`

The S1 resource gets `items`. It is always present, and it is `[]` when the expense has no items. Items keep the order they were sent in. Decimals keep the scale they were sent with (`"0.750"` stays `"0.750"`).

```json
{
  "id": "expense-uuid",
  "occurred_at": "2026-10-05T16:53:34-03:00",
  "amount": "14000.00",
  "currency": "ARS",
  "description": "Gas station store",
  "items": [
    { "description": "Candies", "quantity": "5", "unit_price": "150.00", "line_total": "750.00" }
  ]
}
```

An item has no currency of its own.

## Validation Rules

| Field of an item | Rule | Error codes |
| ---------------- | ---- | ----------- |
| `description` | required, not empty after trimming | `required` |
| `quantity` | required, a decimal, greater than 0 | `required`, `invalid`, `not_positive` |
| `unit_price` | required, a decimal, any sign (zero allowed) | `required`, `invalid` |
| `line_total` | required, a decimal, any sign (zero allowed) | `required`, `invalid` |

* Item errors and expense errors are reported together, as in S1.
* There is no rule between the fields of an item, and none between the items and the amount.
* `items` that is not a JSON array answers `400` with the S1 malformed JSON problem.

## Domain Model

* `Item` is a value inside the `Expense`: description, quantity, unit price, line total. The amounts are `decimal.Decimal`. They are not `Money`: the currency is the expense's, and `Money` keeps its rule that an amount is positive.
* `NewItem` is the only way to build an `Item`. `NewExpense` calls it for each item and returns every violated field. The repository calls it for each stored item.
* The HTTP adapter parses the decimals (as it does for `amount`) and merges its parse errors with the domain's.

## Persistence

* A new table `expense_items`: `expense_id` (references `expenses`, deleted with it), `position` (the order), `description`, `quantity` (greater than 0), `unit_price`, `line_total`. The primary key is `(expense_id, position)`.
* The expense and its items are saved in one transaction.
* Decimals are sent to Postgres as text with their scale, as for the amount, so `0.750` is not stored as `0.75`.
* Loading reads the items in order and builds each one with `NewItem`. If a stored item fails, the review answers the generic `500` and the details go to the log, as for a stored `Money`. The database checks (`quantity > 0`, not null) are a second guard.

## Error Semantics

The same as S1: validation failures answer `422` with `application/problem+json`, internal errors answer the generic `500` and the details go to the log.

## Verification

`make -C backend ci` passes, including the `@s4a` scenarios with `@todo` removed.

## Calls Made During Design (checked and agreed by the product owner)

* `items` is always in the response, `[]` when empty.
* Items keep the order they were sent in.
* Zero is allowed for the unit price and the line total (a free item).
* The description is trimmed, as in S1.
* Item errors use the field name `items[<position>].<field>`.
* There is no limit on the number of items.
* The items go in a separate table, not in a JSON column. The migration is `0002_expense_items`.
* The BDD step `the system indicates the item "<field>" is required` (or `is invalid`) matches an error for that field on any item position.
