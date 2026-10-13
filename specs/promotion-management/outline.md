# Promotion Management: Story Outline

**Status:** draft, not approved

This outline holds draft stories. It is not a specification. Move each story into a `.feature` file when the location of feature files for this context is decided (see [`decisions.md`](../../product/epics/financial-management/features/expense-registration/decisions.md), Cross-Cutting section). Epic: [`promotion-management`](../../product/epics/promotion-management/epic.md).

## Evidence

Two bank promotions, both valid 2026-10-01 to 2026-10-31. Full details: [gas station promotion](../../product/examples/bank-promotion-gas-station.md) and [supermarket promotion](../../product/examples/bank-promotion-supermarket.md).

| | Gas station store | Supermarket |
| --- | --- | --- |
| Discount rate (*porcentaje de descuento*) | 25% | 20% |
| Monthly refund cap (*tope de reintegro*) | $10,000 | $20,000 |
| Minimum purchase | none | $35,000 |
| Days | every day | Friday and Saturday |
| Prerequisite | subscription to the Bank's program | holder of a bank account at the Bank |
| Payment channel | the merchant's app only | QR payment through the Bank's app or a QR wallet app |
| Cards | debit and credit, one card network | debit and credit, several card networks |
| Store scope | only the convenience store, not the fuel | all branches, in person and online |

Shared rules:

* No combination with other promotions.
* Excluded: transfers, account money, virtual wallets, Cuota a Cuota.
* The refund goes to the holder's bank account within 30 days.
* The monthly cap adds up all purchases of the holder.

Open points seen in the examples:

* The supermarket minimum purchase: the subtotal or the charged total?
* The [supermarket receipt](../../product/examples/supermarket-receipt.md) shows "card" while the promotion needs a QR payment. The receipt alone cannot prove a promotion applied.
* The [gas station receipt](../../product/examples/gas-station-receipt.md) shows the legal entity of the company that runs the station, but the promotion names a brand and a store.

## Actor: Promotion Manager

* As a Promotion Manager, I want to register a promotion with its saving, so that it is available to the system.
* As a Promotion Manager, I want to define the conditions of a promotion (minimum purchase, days, prerequisite), so that the system applies it only when it is valid.
* As a Promotion Manager, I want to define the validity period, so that the system ignores expired promotions.
* As a Promotion Manager, I want to define the payment methods and payment channel, so that the system matches the way the user paid.
* As a Promotion Manager, I want to define the commerce, establishment or product scope, so that the promotion targets the right purchases.
* As a Promotion Manager, I want to define the refund cap, the refund delay and the refund destination, so that the user knows the limits.
* As a Promotion Manager, I want to define whether a promotion combines with others, so that the system does not apply two that exclude each other.

## Actor: Financial Manager

* As a Financial Manager, I want the system to find the promotions that apply to my expense, so that I record the expected reduction without searching by hand.
* As a Financial Manager, I want to see how much of the monthly refund cap I already used, so that I know the refund I can still get.
* As a Financial Manager, I want to see my pending refunds and their due dates, so that I can confirm them when they arrive.
* As a Financial Manager, I want to know which promotion to use for a planned purchase, so that I pay the lowest effective cost.

## To refine

* How the system finds and applies a promotion (the exact flow is a delayed decision).
* Which payment information the expense must hold to match a promotion (depends on the payment slice of Expense Registration).
* Whether a promotion is a snapshot on the expense or a reference to the catalog.
