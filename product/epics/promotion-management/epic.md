# Epic: Promotion Management

**Status:** draft, outline only

## Purpose

Let the system know the commercial promotions that reduce the cost of a purchase, find the ones that apply to an expense, and tell the Financial Manager which promotion to use and what its limits are.

## Actors

* Promotion Manager: creates and maintains promotions (see [`actors.md`](../../actors.md)).
* Financial Manager: sees the promotions that apply, the refund limits and the pending refunds.

## Principles

* A promotion is not ordinary income. A refund from a promotion is a reduction of the expense, not new income.
* A promised refund is a forecast until it arrives. The Financial Manager confirms it by hand because the system has no bank integration.
* The system finds and applies promotions. The exact flow is a delayed decision.

## Stories

Draft stories live in [`specs/promotion-management/outline.md`](../../../specs/promotion-management/outline.md). They move into `.feature` files when the location of feature files for this context is decided.

## Related

The reduction of an expense by an actor is part of [Expense Registration](../financial-management/features/expense-registration/feature.md).
