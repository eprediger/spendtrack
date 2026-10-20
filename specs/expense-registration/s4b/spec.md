# S4b — Record Store Discounts: Specification

**Status:** draft, awaiting product owner review

## Traceability

* Story and scenarios: [`record_store_discounts.feature`](../../../backend/features/expenses/record_store_discounts.feature) (`@s4b`)
* Product decisions: [`decisions.md`](../../../product/epics/financial-management/features/expense-registration/decisions.md), section S4b
* Language decision: [`0003-domain-language.md`](../../../adr/0003-domain-language.md)
* Receipt example: [`supermarket-receipt.md`](../../../product/examples/supermarket-receipt.md)
* Builds on: [S4a specification](../s4a/spec.md)

## Scope

Record and review discount amounts already applied to purchase-item prices and shown on a receipt. The operation submits the complete set of item discounts for one expense at once. It replaces the previously recorded set and validates the entire resulting expense before saving.

This slice does not record a promotion rule, detect an applicable promotion, model a later refund, record expense-wide discounts, or link a discount to its merchant. Actor means the real party accountable for applying or paying a reduction. The receipt discount can be linked to its merchant after S4e identifies the merchant.

## Domain terms

* A **Promotion** is a commercial rule.
* A **Reduction** is the concrete amount applied to an expense or item.
* An item discount is a receipt-visible reduction already applied to that item's price. Its amount is positive in the API and domain; it is subtracted from the printed line total when reconciling the charged total.
* Do not record a summary line again when it repeats an applied discount. Do not record savings already embedded in item prices and not subtracted from the charged total.

In the supermarket example, the 1,700.00 "2DO AL 40%" line is attached to the Coca Cola purchase item. "AHORRO OFERTAS" repeats that amount and is not a second discount. "AHORRO CLUBDIA" is not itemized or deducted from the total and is not recorded.

## API contract

### `PUT /api/v1/expenses/{id}/item-discounts`

The request replaces the complete item-discount set for the expense:

```json
{
  "discounts": [
    {
      "item_id": "5f8d3b8e-...",
      "amount": "1700.00",
      "description": "SECOND UNIT AT 40%"
    }
  ]
}
```

* `item_id` is the UUID of the target purchase item as returned in the expense `GET` response; clients may also submit `item_position` for convenience (the server will resolve it to the corresponding `item_id`).

### `POST /api/v1/expenses` (create with per-item discounts)

The existing expense registration endpoint may include per-item discounts in the `items[]` objects when the receipt shows them at registration time. Example excerpt of the create request:

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

* For created expenses, the server generates `item_id` UUIDs for each created purchase item. The response `GET /api/v1/expenses/{id}` will include these `item_id` values. Clients that submit a subsequent `PUT /api/v1/expenses/{id}/item-discounts` may reference items by `item_id` (recommended) or `item_position`.
* The create-time discounts are validated with the same domain rules as replacements: per-item positive amounts, cumulative item cap, and exact expense reconciliation.
* `amount` is a positive decimal string in the expense currency. A request does not choose a separate currency.
* `description` is required and trimmed.
* An expense may have multiple discount entries for one item.
* A missing or empty `discounts` array means the complete set is empty. It is accepted only if the item line totals already equal the charged total.
* A valid request returns `200 OK` with the expense resource. Each `items[]` element includes `discounts`, always an array (empty when the item has none); each discount has `amount` and `description`.
* A valid replacement is atomic. Repeating the same complete request has the same result and does not duplicate discounts.
* If validation fails, the existing set remains unchanged.
* An unknown expense returns the existing fixed `404` problem response.
* A malformed body returns the existing fixed `400` malformed-JSON problem response.
* Invalid discount data or failed reconciliation returns the existing fixed `422` validation problem response. Error fields identify the entry, for example `discounts[0].amount`, `discounts[0].description`, or `discounts[0].item_position`; total mismatch uses `discounts`.

### `GET /api/v1/expenses/{id}`

The existing expense response includes each item's `discounts` array. Entries keep request order and decimal scale, for example `"1700.00"`.

## Domain rules

* Each discount targets exactly one existing purchase item.
* Each discount amount is positive; zero and negative amounts are invalid.
* Cumulative discounts on an item cannot exceed that item's positive `line_total`. Items with zero or negative `line_total` cannot have a discount.
* The expense must have purchase items to accept non-empty discounts.
* Exact reconciliation is required:

  `sum(item.line_total) - sum(item discounts) = expense charged total`

* Validate the complete replacement set before changing the `Expense`. If any validation fails, return all field errors and leave the old discount set unchanged.
* Do not add a merchant or actor field to the item discount in S4b.

## Persistence

* Add `item_id` UUID to `expense_items` in migration `0003_add_item_id_to_expense_items`, populate existing rows with generated UUIDs, and add a unique constraint `(expense_id, item_id)`. Keep `(expense_id, position)` as the table primary key.
* Add `expense_item_discounts` in numbered migration `0004_expense_item_discounts`.
* Store `expense_id UUID NOT NULL`, `item_id UUID NOT NULL`, `discount_position INTEGER NOT NULL`, `amount NUMERIC NOT NULL CHECK (amount > 0)`, and `description TEXT NOT NULL`. The foreign key `(expense_id, item_id)` references `(expense_id, item_id)` in `expense_items` with `ON DELETE CASCADE`. The primary key is `(expense_id, item_id, discount_position)`.
* Save the full replacement in a transaction: delete the prior set, insert the new ordered set, then commit. A failed insert or commit leaves the previous set intact.
* Load discounts in item and discount position order and attach them to the corresponding purchase item. Each returned purchase item in the API includes its `item_id` UUID.
* Persist decimals as scale-preserving text, following the existing expense repository.

## Error handling

Use the existing RFC 9457 problem response at the HTTP boundary. Keep internal database errors out of responses and log them using the existing handler pattern.

## Verification

Run the focused domain and HTTP tests during implementation, then run `make -C backend bdd` serially against the isolated test database, followed by `make -C backend ci`.
