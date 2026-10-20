# Example: Supermarket Receipt

A receipt with many items from a supermarket branch. See the [examples index](README.md) for the rules of these files, including the number notation.

## What it shows

* Ten item lines and one discount line.
* A discount on one item ("second unit at 40%").
* A negative item: the return of returnable bottles.
* A barcode printed under most items. One item has a short internal code. The deposit lines have none.
* A savings block with a loyalty club amount that is not subtracted from the total.
* A card payment. The receipt does not show that a QR payment was used (see [Supermarket promotion](bank-promotion-supermarket.md)).

## Receipt data

| Field | Value |
| ----- | ----- |
| Type | Factura B Original (cód. 006) |
| Date and time | Saturday 2026-10-10, 20:06:20, UTC-03:00 |
| Merchant | A supermarket branch |
| Address (printed on the receipt, omitted here) | A street and number, neighbourhood and city |
| Tax status of merchant | IVA Responsable Inscripto |
| Customer | Consumidor Final |
| CAE | Present, valid until 2026-10-20 |

Items, with the descriptions as printed. The barcode (EAN-13) is printed under the description. The check digits of the barcodes are valid, which supports the transcription.

| Description | Barcode | Quantity | Unit price | Line total | VAT rate |
| ----------- | ------- | -------- | ---------- | ---------- | -------- |
| AEROSOL LAVANDER | 7790520998262 | 2 | 3,090.00 | 6,180.00 | 21% |
| P HIG PREMIUM BALANC | 7790250015543 | 1 | 5,955.00 | 5,955.00 | 21% |
| ATUN LOM N CAM 170GR | 7790580131364 | 2 | 6,278.00 | 12,556.00 | 21% |
| JAB TOC COTTON FRESH | 7791293050935 | 2 | 4,345.00 | 8,690.00 | 21% |
| PALTA | 228219 (internal code) | 4 | 990.00 | 3,960.00 | 10.5% |
| COCA COLA R 2L | 7790895000218 | 3 | 4,250.00 | 12,750.00 | 21% |
| 2DO AL 40% LINEA COC (discount on the line above) | | | | -1,700.00 | |
| ENVASE RETORNABLE | | 3 | 3,001.00 | 9,003.00 | 21% |
| YERBA MATE CAMPO ORI | 7790387015607 | 1 | 4,770.00 | 4,770.00 | 21% |
| ACEITE DE GIRASOL | 7792180001665 | 1 | 6,529.00 | 6,529.00 | 21% |
| ENVASE RETORNABLE (returned) | | 3 | -3,001.00 | -9,003.00 | 21% |

| Field | Value |
| ----- | ----- |
| Subtotal | 61,390.00 |
| "DESCUENTOS" (the sum of the store's own promotions) | -1,700.00 |
| Total | 59,690.00 |
| Payment | Card, 59,690.00 |
| Change | 0.00 |
| VAT included | 10,048.44 |
| Other national indirect taxes | 794.50 |
| Savings from offers ("AHORRO OFERTAS") | 1,700.00 |
| Savings from the loyalty club ("AHORRO CLUBDIA") | 4,460.00 (not subtracted from the total; it seems to be inside the prices, which is an inference) |
| Total savings printed ("TOT AHORRO") | 6,160.00 |

The item lines add up to the subtotal: the deposit and its return cancel each other.

## What follows from it

The promotion in [Supermarket promotion](bank-promotion-supermarket.md) applies to this purchase (a Saturday, a QR payment with a credit card, above the minimum purchase):

| | Amount |
| - | ------ |
| Charged total (the expense amount) | 59,690.00 |
| Expected bank refund (20%, under the 20,000 monthly cap) | 11,938.00 |
| Effective cost (*costo de bolsillo*) | 47,752.00 |

Notes:

* The receipt says only "card". The promotion needs a QR payment through the bank app. The receipt cannot show that the promotion applied.
* The minimum purchase is met with both the subtotal (61,390.00) and the total (59,690.00). It is not known which one the Bank uses.

## Illustrates

* S4a: items, including a negative item. S4b: a discount on one item. S4c and S4d: effective cost and an expected bank reduction.
* S4b: the loyalty club savings line is not recorded because it is not allocated to an item or subtracted from the charged total.
* Price Analysis: barcodes under most items, and descriptions that are cut off.
