# Expense Registration — Decisions

## Firm Agreements

### Expense Is Economic Consumption

An expense represents economic consumption and must not be conflated with payment or cash movement.

### Purchase Is Part of the Expense Representation

A purchase provides detailed information about an expense resulting from a commercial transaction. It must not create a second independent expense.

### Supporting Documents Are Evidence

A document provides evidence of an expense but is not required for registration.

### Shared Context Does Not Imply Ownership

Being part of a shared Financial Context does not automatically make an expense economically shared.

### Future Payment Commitments Are Distinct

Installments and other future settlements must remain distinguishable from the expense itself.

### S1 Required Attributes and Validation Rules

For the first slice, an expense requires occurrence date, amount, currency, and description at registration. Registration is rejected when required information is missing or invalid: the amount must be positive, the currency must be a valid ISO 4217 code, and the occurrence date must not be in the future.

This agreement covers S1 only; later slices may introduce additional attributes and their own validation rules.

---

## Story Map

The following implementation slicing is agreed. Each slice activates the domain concepts its Acceptance Criteria and BDD scenarios actually require.

| Slice | Story                                                                      | Concepts activated                                              |
| ----- | -------------------------------------------------------------------------- | --------------------------------------------------------------- |
| S1    | [Register a basic expense (date, amount, currency, description) and review it](../../../../../backend/features/expenses/register_basic_expense.feature) | Expense, Financial Context |
| S2    | Classify the expense                                                       | Classification                                                  |
| S3    | Associate supporting documentation (optional)                              | Document                                                        |
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

Open decisions must be resolved through the relevant Product Discovery and Feature refinement process.

