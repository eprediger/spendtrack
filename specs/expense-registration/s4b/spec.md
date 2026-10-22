# S4b — Record Store Discounts: Specification

**Status:** approved

## Traceability

* Story and scenarios: [`record_store_discounts.feature`](../../../backend/features/expenses/record_store_discounts.feature) (`@s4b`)
* Product decisions: [`decisions.md`](../../../product/epics/financial-management/features/expense-registration/decisions.md), section S4b
* Language decision: [`0003-domain-language.md`](../../../adr/0003-domain-language.md)
* Receipt example: [`supermarket-receipt.md`](../../../product/examples/supermarket-receipt.md)
* Builds on: [S4a specification](../s4a/spec.md)

## Scope

Record and review discount amounts already applied to purchase-item prices and shown on a receipt. Discounts are sent at registration time, nested inside the purchase item they apply to. Every itemized expense must reconcile exactly: the sum of item line totals minus the sum of all item discounts equals the charged total.

This slice does not record a promotion rule, detect an applicable promotion, model a later refund, record expense-wide discounts, or link a discount to its merchant.

## Domain terms

* A **Promotion** is a commercial rule.
* A **Reduction** is the concrete amount applied to an expense or item.
* An item discount is a receipt-visible reduction already applied to that item's price. Its amount is positive in the API and domain; it is subtracted from the printed line total when reconciling the charged total.
* Do not record a summary line again when it repeats an applied discount. Do not record savings already embedded in item prices and not subtracted from the charged total.

In the supermarket example, the 1,700.00 "2DO AL 40%" line is attached to the Coca Cola purchase item. "AHORRO OFERTAS" repeats that amount and is not a second discount. "AHORRO CLUBDIA" is not itemized or deducted from the total and is not recorded.

## API contract

### `POST /api/v1/expenses`

The existing registration endpoint accepts an optional `discounts` array inside each `items[]` object:

```json
{
  "occurred_at": "2026-10-05T16:53:34-03:00",
  "amount": "59690.00",
  "currency": "ARS",
  "description": "Supermarket",
  "items": [
    {
      "description": "COCA COLA R 2L",
      "quantity": "1",
      "unit_price": "12750.00",
      "line_total": "12750.00",
      "discounts": [
        { "amount": "1700.00", "description": "2DO AL 40% LINEA COC" }
      ]
    }
  ]
}
```

* Each discount has `amount` (a positive decimal string in the expense currency) and `description` (required, trimmed).
* An item may have several discount entries. A missing or empty `discounts` array means the item has none.
* A discount needs no item reference: it belongs to the item that carries it.

### `POST` and `GET` responses

Each `items[]` element in the response includes `discounts`, always an array (empty when the item has none). Each entry has `amount` and `description`. Amounts keep the entered decimal scale, for example `"1700.00"`. Entries keep request order.

### Reconciliation

Every expense with purchase items must satisfy:

`sum(item.line_total) - sum(all item discounts) = amount`

The comparison uses `decimal.Equal`, so it ignores scale. An expense without items is not checked. The check runs only when no other field error exists.

## Field errors

Invalid data returns the existing fixed `422` validation problem response. Error fields are prefixed with the item index:

| Field | Code | Cause |
| --- | --- | --- |
| `items` | `not_reconciled` | item line totals minus discounts do not equal `amount` (detail: "item line totals minus discounts must equal amount") |
| `items[i].discounts[j].description` | `required` | missing or blank description |
| `items[i].discounts[j].amount` | `required` | missing amount |
| `items[i].discounts[j].amount` | `not_positive` | zero or negative amount |
| `items[i].discounts[j].amount` | `invalid` | not a decimal number (reported by the HTTP adapter and merged through `withParseErrors`) |
| `items[i].discounts` | `exceeds_line_total` | all entries are valid and their sum exceeds `max(line_total, 0)`; this also rejects any discount on a zero or negative item |

## Domain rules

* Each discount amount is positive; zero and negative amounts are invalid.
* Cumulative discounts on an item cannot exceed that item's positive `line_total`. A discount equal to the line total is allowed. Items with a zero or negative `line_total` cannot have a discount.
* Exact reconciliation is required for every itemized expense, with or without discounts.
* Do not add a merchant or actor field to the item discount in S4b.

## Persistence

* Migration `0003_expense_item_discounts` creates `expense_item_discounts` with columns `expense_id UUID NOT NULL`, `item_position INTEGER NOT NULL`, `position INTEGER NOT NULL`, `description TEXT NOT NULL`, `amount NUMERIC NOT NULL CHECK (amount > 0)`; primary key `(expense_id, item_position, position)`; foreign key `(expense_id, item_position)` referencing `expense_items(expense_id, position)` with `ON DELETE CASCADE`.
* `Save` inserts the discounts in the same transaction, after each item.
* `FindByID` loads discounts with a second query ordered by `item_position, position`, attaches them to `PurchaseItemInput.Discounts`, and validates the expense through `NewExpense` as it does today.
* There is no `item_id` column and no new port or service method.
* Decimals persist as scale-preserving text, following the existing repository.

Consequence: a stored itemized expense that does not reconcile cannot be read; the API answers the generic internal error. This is accepted because the product has no production data yet.

## Deferred

Editing the discounts of an existing expense (a `PUT` replacement route) is an open decision in [`decisions.md`](../../../product/epics/financial-management/features/expense-registration/decisions.md). The two replacement scenarios stay `@todo` in the feature file and are not implemented in this slice.

## Error handling

Use the existing RFC 9457 problem response at the HTTP boundary. Keep internal database errors out of responses and log them using the existing handler pattern.

## Verification

Run the focused domain and HTTP tests during implementation, then run `make -C backend bdd` serially against the isolated test database, followed by `make -C backend ci`.
