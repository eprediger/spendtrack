# Epic: Price Analysis

**Status:** draft, outline only

## Purpose

Let a Price Analyst compare the price of products over time, across businesses and across places, using aggregated information from registered expenses.

## Actors

* Price Analyst (see [`actors.md`](../../actors.md)).

## Principles

* The Price Analyst gets aggregated information. The Price Analyst does not get the personal financial information of a user or identifiable individual expenses.
* Analysis needs structured data. Items keep their unit price, and the establishment keeps its address, so that this analysis is possible later.

## Stories

Draft stories live in [`specs/price-analysis/outline.md`](../../../specs/price-analysis/outline.md). They move into `.feature` files when the location of feature files for this context is decided.

## Related

Items and establishments come from [Expense Registration](../financial-management/features/expense-registration/feature.md).
