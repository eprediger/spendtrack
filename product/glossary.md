# Domain Glossary

## Term Status

Every term carries an implementation status:

* **provisional** — the concept is defined, but no approved implementation slice has activated it yet.
* **validated** — the concept's behavior is implemented and verified against an approved specification.

All terms begin as provisional and become validated as approved User Stories activate them through implementation slices.

Approval status (firm agreement vs. open decision) is a separate axis, tracked in each Feature's `decisions.md`. Where a term is backed by a Firm Agreement, this is noted on the term.

## Expense

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

Economic consumption recognized as a financial fact.

An expense is not synonymous with a cash outflow or payment.

## Movement

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

A financial movement or settlement involving financial resources.

A movement is not necessarily an expense.

## Purchase

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

A detailed representation of an expense resulting from a commercial transaction.

A purchase does not create a second independent expense.

## Payment Commitment

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

A future commitment to settle an amount related to an expense or other obligation.

A payment commitment is not itself the economic consumption.

## Obligation

**Status:** provisional

A future or existing financial responsibility that may require settlement.

An obligation is distinct from the expense that may have originated it.

## Financial Context

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

An environment in which financial information is managed and shared among one or more people.

Membership in a Financial Context does not imply economic ownership.

## Economic Ownership

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

The person or persons to whom an asset, position, income, expense, or economic responsibility belongs.

Economic ownership is distinct from account custody.

## Custody

**Status:** provisional

The financial account or institution through which an asset or position is held.

The custodian does not necessarily represent the economic owner.

## Financial Position

**Status:** provisional

An economic ownership position in an asset.

A financial position may be distributed among multiple economic owners.

## Asset

**Status:** provisional

An economic resource with value.

## Liability

**Status:** provisional

An economic responsibility or amount owed.

## Net Worth

**Status:** provisional

Derived economic value calculated from economic assets/positions minus relevant liabilities and obligations.

Net worth is not a primary historical fact.

## Savings

**Status:** provisional

Accumulation of financial resources that are not necessarily invested.

Savings and investments are distinct concepts.

## Investment

**Status:** provisional

A financial asset or position acquired or maintained with an investment purpose.

## Objective

**Status:** provisional

A financial goal with a target and potentially a time horizon.

An objective is distinct from the asset or position used to fund it.

## Budget

**Status:** provisional

A set of financial targets or allocation rules used to guide and evaluate financial behavior.

## Commerce

**Status:** provisional

A business or brand associated with a commercial transaction.

## Establishment

**Status:** provisional

A physical location where a commercial transaction occurs.

## Product

**Status:** provisional

A product or service acquired in a purchase.

## Promotion

**Status:** provisional

A commercial condition that can provide a discount, benefit, or reimbursement when its conditions are satisfied.

## Reimbursement

**Status:** provisional

A later financial benefit associated with a previous purchase or promotion.

A reimbursement is not ordinary income.

## Forecast

**Status:** provisional

A derived expectation about future financial behavior.

A forecast must not modify historical facts.

## Recurring Rule

**Status:** provisional

A rule describing how an obligation or financial event may recur.

A recurring rule is distinct from a concrete occurrence.

## Document

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

Evidence associated with a financial event, such as a receipt, invoice, image, or PDF.
