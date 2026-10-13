# Example: Gas Station Receipt

A simple receipt from the convenience store of a gas station. See the [examples index](README.md) for the rules of these files, including the number notation.

## What it shows

* A short list of items with no discount.
* No payment promotion on the receipt. The bank refund comes later (see [Gas station promotion](bank-promotion-gas-station.md)).
* The receipt names the company that runs the station (a legal entity). It does not name the fuel brand or the store brand.
* The address is a national highway and a kilometer, not a street and a number.
* Paid through the merchant's own app.
* No CAE printed. Only a QR code.

## Receipt data

| Field | Value |
| ----- | ----- |
| Type | Tique Factura "B" (cód. 082) |
| Date and time | Monday 2026-10-05, 16:53:34, UTC-03:00 |
| Merchant | A convenience store in a gas station complex |
| Address (printed on the receipt, omitted here) | A national highway and kilometer, city, province and postal code |
| Tax status of merchant | IVA Responsable Inscripto |
| Customer | Consumidor Final |

Items, with the descriptions as printed. Each description ends with a VAT rate, `(21)`. The last two have a second number in brackets, `[17,15]` and `[17,28]`, whose meaning is not known. Quantities are printed as `5,0000 u` (units, with four decimals).

| Description | Quantity | Unit price | Line total |
| ----------- | -------- | ---------- | ---------- |
| Caramelos Arcor Frutal 1 Un | 5 | 150.00 | 750.00 |
| Agua Villavicencio Sin | 1 | 3,400.00 | 3,400.00 |
| Triple Blanco Jamon y | 1 | 9,850.00 | 9,850.00 |

| Field | Value |
| ----- | ----- |
| Total | 14,000.00 |
| Payment | Through the merchant's app, "other payment methods", 14,000.00 |
| Merchant discount printed | 0.00 |
| Change | 0.00 |
| VAT included | 2,416.27 |
| Other national indirect taxes | 77.68 |

## What follows from it

The promotion in [Gas station promotion](bank-promotion-gas-station.md) applies to this purchase (paid through the merchant's app, in the store, within the valid period):

| | Amount |
| - | ------ |
| Charged total (the expense amount) | 14,000.00 |
| Expected bank refund (25%, under the 10,000 monthly cap) | 3,500.00 |
| Effective cost (*costo de bolsillo*) | 10,500.00 |

With a receipt of 50,000.00 the refund would be 12,500.00, cut to 10,000.00 by the monthly cap.

## Illustrates

* S4a: items. S4d: an expected bank reduction.
* The legal entity on the receipt is not the brand or store that the promotion names (see S4e, commerce and establishment).
* An establishment address can be a highway and a kilometer.
* Price Analysis: descriptions with no barcode.
