# Expense Registration — Decisions

How to read this file:

* The story map is the index. Each slice below has its own section.
* In each section, **Agreements** are approved. **Open decisions** are unresolved and must not be treated as requirements.
* Decisions that apply to more than one slice are in **Feature-Wide Agreements** or in the **Cross-Cutting** section.
* Open decisions are resolved through the relevant Product Discovery and Feature refinement process, before the slice is scheduled.

---

## Story Map

The following implementation slicing is agreed. Each slice activates the domain concepts its Acceptance Criteria and BDD scenarios actually require.

| Slice | Story                                                                                                                                                            | Concepts activated                              |
| ----- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------- |
| S1    | [Register a basic expense (date and time, amount, currency, description) and review it](../../../../../backend/features/expenses/register_basic_expense.feature) | Expense, Financial Context                      |
| S2    | Classify the expense                                                                                                                                             | Classification                                  |
| S3    | Associate supporting documentation (expenses may have none)                                                                                                      | Document                                        |
| S4a   | [Break the expense down into items](../../../../../backend/features/expenses/itemize_expense.feature)                                                            | Purchase, Item, Product                         |
| S4b   | [Record store discounts](../../../../../backend/features/expenses/record_store_discounts.feature)                                                                | Reduction                                       |
| S4c   | [Review the effective cost](../../../../../backend/features/expenses/review_effective_cost.feature)                                                              | Charged Total, Effective Cost                   |
| S4d   | [Record a reduction from another actor (expected, then received)](../../../../../backend/features/expenses/record_expected_reduction.feature)                    | Reduction, Promotion, Reimbursement             |
| S4e   | [Identify commerce and establishment](../../../../../backend/features/expenses/identify_commerce_establishment.feature)                                          | Commerce, Establishment                         |
| S4f   | [Correct a reduction applied wrongly](../../../../../backend/features/expenses/correct_reduction.feature)                                                        | Reduction                                       |
| S5    | Indicate people and economic allocation                                                                                                                          | Person, Economic Ownership, Economic Allocation |
| S6    | Represent payment mechanism / settlement                                                                                                                         | Movement, Financial Account                     |
| S7    | Installment purchase generates future payment commitments                                                                                                        | Payment Commitment, Obligation                  |

S1 is the walking skeleton: an expense exists without payment or document, so the first slice registers the economic fact and reads it back.

Each slice is specified by an executable `.feature` file — the single source for its User Story, Acceptance Criteria, and BDD scenarios — linked above once drafted. Status of the stories:

* S1 is approved and implemented.
* S4a–S4f are drafted broadly and tagged `@todo`. Their scenarios are a first outline and are refined before each slice is scheduled. The first target is S4a–S4d plus S4f.
* S2, S3 and S5–S7 are intentionally not drafted yet: each `.feature` is written when its slice is scheduled, applying the last responsible moment principle. Their absence is a deferral, not outstanding work.

The promotions behind the bank reductions belong to [Promotion Management](../../../promotion-management/epic.md). The Price Analyst needs for items and establishments belong to [Price Analysis](../../../price-analysis/epic.md).

---

## Feature-Wide Agreements

### Expense Is Economic Consumption

An expense represents economic consumption. It is not the same thing as a payment or a cash movement.

### Purchase Is Part of the Expense Representation

A purchase provides detailed information about an expense resulting from a commercial transaction. It must not create a second independent expense.

### Supporting Documents Are Evidence

A document provides evidence of an expense but is not required for registration.

### Shared Context Does Not Imply Ownership

Being part of a shared Financial Context does not automatically make an expense economically shared.

### Future Payment Commitments Are Distinct

Installments and other future settlements must remain distinguishable from the expense itself.

---

## S1 — Register a Basic Expense

### Agreements

#### S1 Required Attributes and Validation Rules

For the first slice, an expense requires occurrence date and time, amount, currency, and description at registration. Registration is rejected when required information is missing or invalid: the amount must be positive, the currency must be in the supported currency set, and the occurrence must not be in the future. The occurrence preserves both the absolute instant and the local wall-clock time, so time-windowed promotions can be analyzed later.

#### Supported Currencies Are a Curated Set

S1 accepts a deliberately small set of currency codes (`ARS`, `USD`), not the full ISO 4217 list. The set grows as the product needs more currencies. Every accepted code is a valid ISO 4217 alphabetic code; codes outside the set are rejected as `unsupported_currency`. Maintaining (or depending on) the full ISO list is avoided: such lists drift and add maintenance with no current product need.

#### Amounts Are Not Checked Against the Currency Decimal Exponent

S1 validates that the currency is supported, but not that the number of decimals in the amount matches the ISO 4217 exponent of the currency (for example JPY has 0 decimals, BHD has 3). Both supported currencies (`ARS`, `USD`) use 2 decimals, so the check adds code and no current benefit.

Consequence: an amount with extra decimals (for example `10.123 ARS`) is accepted and stored as given. If the check is added later, existing data may need cleanup or rounding.

Revisit when any of these happens:

* a currency with an exponent other than 2 is added to the supported set;
* totals from items and discounts need exact reconciliation;
* a client or report depends on a fixed number of decimals.

If the check is added, apply it to amounts and line totals only. Unit prices and quantities on receipts can have more decimals than the currency exponent.

These agreements cover S1 only; later slices may introduce additional attributes and their own validation rules.

### Open Decisions

None.

---

## S2 — Classify the Expense

### Agreements

None yet.

### Open Decisions

* Exact expense classification structure.
* Whether fixed/variable and essential/discretionary are independent classification dimensions.
* How customizable classifications should be.

---

## S3 — Associate Supporting Documentation

### Agreements

None yet. The feature-wide agreement "Supporting Documents Are Evidence" applies.

### Open Decisions

* Supported document types in the first implementation slice.
* Whether manual transcription from documents is part of the first implementation slice.

---

## S4a — Break the Expense Down Into Items

### Agreements

#### Items Carry Description, Quantity, Unit Price and Line Total

An expense can be broken down into items. Each item has a description, a quantity, a unit price and a line total. The line total is the amount the receipt prints. It can differ from quantity times unit price. Items are optional. A line can be negative (for example a returned deposit).

### Open Decisions

* Any other item information required beyond description, quantity, unit price and line total (for example tax rate).
* Product identification behavior, including when the receipt prints a barcode.

---

## S4b — Record Store Discounts

### Agreements

#### Reductions Are One Concept, Attributed to an Actor

Store discounts and reductions from a bank or another actor are one concept: a reduction. Each reduction names the actor that applied it. A reduction applies to the whole expense or to one item.

This agreement also applies to S4d and S4f.

### Open Decisions

* Whether items and store discounts must add up to the charged total, and what happens when they do not.
* Whether informational savings printed on a receipt (for example the loyalty club savings line of the [supermarket receipt](../../../../examples/supermarket-receipt.md), already inside prices) are recorded or ignored.
* Whether a store sale price cut is the same concept as a program promotion.
* Whether "Reduction" stays one concept or splits into the store's discount (*descuento*) and the bank's refund (*reintegro*). In Spanish, users keep them apart. See the **En español** line of Reduction in the glossary. Deferred: a concrete example is needed before deciding. Revisit when S4b and S4d are refined. "Reduction" is also not a common English term for this (English uses discount, rebate, cashback), so the English name may change too.

---

## S4c — Review the Effective Cost

### Agreements

#### Effective Cost Is the Headline Number

The Financial Manager sees the effective cost: the money that really left the pocket. It is the charged total minus the reductions from other actors.

#### The Expense Amount Is the Charged Total

The amount registered on an expense (S1) is the charged total: what the store charged, as printed on the receipt, after the store's own discounts. Reductions from other actors, such as a bank refund, are not part of that amount. The effective cost is derived from it and shown next to it.

Example ([supermarket receipt](../../../../examples/supermarket-receipt.md) and [bank promotion](../../../../examples/bank-promotion-supermarket.md)): the expense amount is 59,690.00 ARS. The bank promotion refunds 20%, which is 11,938.00 ARS. The effective cost is 47,752.00 ARS.

### Open Decisions

* Which number budgets and reports use: the charged total (the expense amount) or the effective cost.
* Whether a reduction can be larger than the charged total.

---

## S4d — Record a Reduction From Another Actor

### Agreements

#### Reductions Are Expected Until Confirmed

A reduction from another actor (for example a bank refund) starts as expected. The Financial Manager confirms it as received. While a reduction is only expected, the effective cost is an estimate. The system has no bank integration, so the confirmation is manual.

### Open Decisions

* Which payment information (instrument, channel, holder) the expense holds, so that promotions can be matched. This pulls part of S6 forward.

---

## S4e — Identify Commerce and Establishment

### Agreements

#### An Establishment Keeps Its Exact Address

An establishment is recorded with its exact address, as printed on the receipt, and not only its city or region. The address must be exact enough to find its coordinates (latitude and longitude) later. The coordinates themselves are not captured now.

The receipts show two forms of address: a street and a number, and a national highway with a kilometer (see the [gas station receipt](../../../../examples/gas-station-receipt.md)).

### Open Decisions

* Commerce and Establishment registration behavior.
* How the coordinates of an establishment are found and stored from its address, and when.
* How an address is structured (street and number, highway and kilometer, city, province, postal code).
* Which fiscal receipt fields (tax id, receipt number, CAE, taxes) are captured.

---

## S4f — Correct a Reduction Applied Wrongly

### Agreements

#### Reductions Can Change After Registration

The Financial Manager can add a reduction after the expense was registered, and can correct one that was applied wrongly. When the outcome differs from the expectation, the record keeps the expected amount next to the outcome (different amount, not received). A reduction that never applied can be dismissed. The effective cost uses the actual outcome. A full change history is not required for now.

### Open Decisions

None yet.

---

## S5 — Indicate People and Economic Allocation

### Agreements

None yet. The feature-wide agreement "Shared Context Does Not Imply Ownership" applies.

### Open Decisions

* When explicit economic allocation between people is required.
* Exact economic allocation mechanisms.

---

## S6 — Represent Payment Mechanism and Settlement

### Agreements

#### One Expense Has One Payment Method, For Now

For now, an expense is paid with one and only one payment method. Receipts that show the sum of several payments ("Suma de sus pagos") are not modeled yet. The payment lines of a receipt are optional: some receipts print none (for example the [bakery receipt](../../../../examples/bakery-receipt.md)).

### Open Decisions

* Whether one expense can be paid with several payment methods at the same time (a feature to explore).
* Exact payment model.
* Payment methods supported by the first implementation slice.
* How much of the payment information needed by S4d (instrument, channel, holder) belongs here.

---

## S7 — Installment Purchase and Payment Commitments

### Agreements

None yet. The feature-wide agreement "Future Payment Commitments Are Distinct" applies.

### Open Decisions

* Exact installment/payment-commitment behavior.

---

## Cross-Cutting

### Agreements

#### Artifact Placement for New Contexts

Stories for the expense extensions are `.feature` files in `backend/features/expenses/`. Promotion Management and Price Analysis have no implementation, and the application that will host them is not decided. Their stories are draft outlines in `specs/promotion-management/` and `specs/price-analysis/`. This is a temporary exception to the rule that stories live only in `.feature` files.

#### Examples Are Written Down, Anonymous, Without Images

Real receipts and promotions that give context to the definitions are transcribed as text in [`product/examples/`](../../../../examples/README.md). They are anonymous: merchants are described by type (supermarket, gas station, bakery), the bank is "the Bank", and personal identifiers, document numbers and the exact address of the merchant are left out. Item descriptions are kept as printed. Numbers are consistent (quantity times unit price equals the line total; the lines add up to the total) but a digit can be wrong. The original photos and screenshots are not versioned and stay outside the repository.

### Open Decisions

* Location of the feature files for Promotion Management and Price Analysis (a technical decision, deferred).
* Location of BDD feature files if additional applications are added to the repository (currently `backend/features/`).
