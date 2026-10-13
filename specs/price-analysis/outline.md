# Price Analysis: Story Outline

**Status:** draft, not approved

This outline holds draft stories. It is not a specification. Move each story into a `.feature` file when the location of feature files for this context is decided (see [`decisions.md`](../../product/epics/financial-management/features/expense-registration/decisions.md), Cross-Cutting section). Epic: [`price-analysis`](../../product/epics/price-analysis/epic.md).

## Evidence

* The [supermarket receipt](../../product/examples/supermarket-receipt.md) prints a barcode under most items. One item has a short internal code. The [gas station receipt](../../product/examples/gas-station-receipt.md) prints none.
* Item descriptions are short, in capital letters and cut off. The same product can have a different description in each business.
* The same product can have a different price in a different place. The receipts show the address of the establishment.

## Actor: Price Analyst

* As a Price Analyst, I want to identify the same product across businesses, so that I can compare its price.
* As a Price Analyst, I want to see the price of a product over time, so that I can see how it changes.
* As a Price Analyst, I want to compare the price of a product by place, so that I can see geographic differences.
* As a Price Analyst, I want to see aggregated information only, so that I do not access the personal financial information of a user.

## To refine

* How to identify a product when no barcode is printed.
* How coordinates (latitude and longitude) are found from the exact address of the establishment. The address is kept (see S4e in the Expense Registration decisions). The coordinates come later.
* Which price counts: the unit price printed, or the price after store discounts.
* How the Price Analyst gets data from the expenses of many users without seeing individual expenses.
