# Example: Bakery Receipt

A receipt from a bakery in the city of Buenos Aires. See the [examples index](README.md) for the rules of these files, including the number notation.

## What it shows

* No payment lines: no payment method, no amount paid, no change. Payment lines are optional on a receipt.
* No cashier, no store number and no POS number.
* Quantities printed with three decimals (`2,000` means 2 units).
* A line with a local gross income tax rate (1.50%).
* No promotion is known for this purchase.

## Receipt data

| Field | Value |
| ----- | ----- |
| Type | Factura B Original (cód. 006) |
| Date and time | Sunday 2026-10-11, 17:51, UTC-03:00 |
| Merchant | A bakery |
| Address (printed on the receipt, omitted here) | A street and number, in the city of Buenos Aires (CABA) |
| Tax status of merchant | Responsable Inscripto |
| Customer | A consumidor final |
| CAE | Present, valid until 2026-10-21 |

Items, with the descriptions as printed. Quantities are printed as `2,000` (2 units). Unit prices are printed with four decimals (`4.100,0000`).

| Description | Quantity | Unit price | Line total |
| ----------- | -------- | ---------- | ---------- |
| T-queso y Huevo | 2 | 4,100.00 | 8,200.00 |
| T-jam y Que BLANCO | 2 | 4,100.00 | 8,200.00 |
| Alfajor Maicena | 1 | 4,500.00 | 4,500.00 |
| Alfa Coco | 1 | 4,500.00 | 4,500.00 |

| Field | Value |
| ----- | ----- |
| Total | 25,400.00 |
| VAT included | 4,054.87 |
| Other taxes | 0.00 |
| Gross income tax rate (CABA) | 1.50% |

## Illustrates

* S4a: items. A quantity is a decimal number (printed with three decimals), not only a whole number.
* S6: for now an expense is paid with one payment method, and the receipt does not always say which.
* The glossary note that payment lines are optional.
* Price Analysis: descriptions that are abbreviated in a different way than on other receipts.
