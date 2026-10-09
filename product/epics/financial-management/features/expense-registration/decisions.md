# Expense Registration — Decisions

## Firm Agreements

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

### S1 Required Attributes and Validation Rules

For the first slice, an expense requires occurrence date and time, amount, currency, and description at registration. Registration is rejected when required information is missing or invalid: the amount must be positive, the currency must be in the supported currency set, and the occurrence must not be in the future. The occurrence preserves both the absolute instant and the local wall-clock time, so time-windowed promotions can be analyzed later.

### Supported Currencies Are a Curated Set

S1 accepts a deliberately small set of currency codes (`ARS`, `USD`), not the full ISO 4217 list. The set grows as the product needs more currencies. Every accepted code is a valid ISO 4217 alphabetic code; codes outside the set are rejected as `unsupported_currency`. Maintaining (or depending on) the full ISO list is avoided: such lists drift and add maintenance with no current product need.

This agreement covers S1 only; later slices may introduce additional attributes and their own validation rules.

---

## Story Map

The following implementation slicing is agreed. Each slice activates the domain concepts its Acceptance Criteria and BDD scenarios actually require.

| Slice | Story                                                                      | Concepts activated                                              |
| ----- | -------------------------------------------------------------------------- | --------------------------------------------------------------- |
| S1    | [Register a basic expense (date and time, amount, currency, description) and review it](../../../../../backend/features/expenses/register_basic_expense.feature) | Expense, Financial Context |
| S2    | Classify the expense                                                       | Classification                                                  |
| S3    | Associate supporting documentation (expenses may have none)                | Document                                                        |
| S4    | Detail the expense as a purchase (commerce, establishment, items)          | Purchase, Commerce, Establishment, Product                      |
| S5    | Indicate people and economic allocation                                    | Person, Economic Ownership, Economic Allocation                 |
| S6    | Represent payment mechanism / settlement                                   | Movement, Financial Account                                     |
| S7    | Installment purchase generates future payment commitments                  | Payment Commitment, Obligation                                  |

S1 is the walking skeleton: an expense exists without payment or document, so the first slice registers the economic fact and reads it back.

Each slice is specified by an executable `.feature` file — the single source for its User Story, Acceptance Criteria, and BDD scenarios — linked above once drafted.

Stories S2–S7 are intentionally not drafted yet: each `.feature` is written when its slice is scheduled, applying the last responsible moment principle. Their absence is a deferral, not outstanding work.

---

## Open Decisions

The following decisions remain unresolved and must not be treated as requirements:

* Exact expense classification structure.
* Whether fixed/variable and essential/discretionary are independent classification dimensions.
* How customizable classifications should be.
* When explicit economic allocation between people is required.
* Exact economic allocation mechanisms.
* Exact payment model.
* Payment methods supported by the first implementation slice.
* Exact installment/payment-commitment behavior.
* Supported document types in the first implementation slice.
* Whether manual transcription from documents is part of the first implementation slice.
* Exact purchase item information required.
* Product identification behavior.
* Commerce and Establishment registration behavior.
* Acceptance Criteria and BDD scenarios per story in the agreed Story Map (S1 approved, S2–S7 pending).
* Location of BDD feature files if additional applications are added to the repository (currently `backend/features/`).
* Whether enforcing per-currency decimal exponents (ISO 4217 minor units — e.g., JPY has 0 decimals, BHD has 3) is required or overkill. S1 validates the currency code but not exponent-consistent amounts; revisit when analyzing `Money` representation.

Open decisions must be resolved through the relevant Product Discovery and Feature refinement process.

