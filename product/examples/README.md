# Examples

Real receipts and promotions, written down as text. They give context to the definitions in [`glossary.md`](../glossary.md), the [Expense Registration decisions](../epics/financial-management/features/expense-registration/decisions.md) and the draft outlines in `specs/`.

## Examples

| Example | Type | Date | What it shows |
| ------- | ---- | ---- | ------------- |
| [Gas station receipt](gas-station-receipt.md) | Receipt | 2026-10-05 | A simple receipt. Three items, no discount. The receipt does not show the bank promotion. |
| [Supermarket receipt](supermarket-receipt.md) | Receipt | 2026-10-10 | A complex receipt. Many items, a discount on one item, a negative item, a loyalty club savings line. |
| [Bakery receipt](bakery-receipt.md) | Receipt | 2026-10-11 | A receipt with no payment lines. Quantities with three decimals. |
| [Gas station promotion](bank-promotion-gas-station.md) | Promotion | October 2026 | A bank refund: percentage, monthly cap, payment channel, one store only. |
| [Supermarket promotion](bank-promotion-supermarket.md) | Promotion | October 2026 | A bank refund: minimum purchase, valid days, payment channel. |

The gas station and supermarket receipts match the two promotions. Each file shows the refund and the cost of the expense that follow from them.

## Rules for these files

* **Anonymous merchants.** The brand, legal name, tax ID and address of a merchant, document numbers, cashier names and loyalty card numbers are left out. A merchant is described by its type (supermarket, gas station, bakery) and the bank is "the Bank". The real receipt prints the exact address of the branch, and the system must keep it (see S4e in the Expense Registration decisions).
* **Item descriptions are kept as printed.** They are short, in capital letters and often cut off. This is on purpose: it shows how hard it is to identify a product from its description. Product names can be brands.
* **Numbers make sense, they are not always exact.** Quantity times unit price equals the line total, and the lines add up to the total. A digit can be wrong, because the data was transcribed by eye.
* **Number notation.** The receipts use a dot as the thousands separator and a comma as the decimal separator (`25.400,00` is twenty-five thousand four hundred). These files use the English notation (`25,400.00`). Quantities are units: `2,000` on a receipt means 2 units, not two thousand.
* **Images are not versioned.** The photos and screenshots stay outside the repository. They show personal identifiers (for example a cashier name, a loyalty card number and a transaction number).
* **Promotions change.** A promotion example is a snapshot with a date. It is not a record of the current terms.
* **Amounts** are in Argentine pesos (ARS).
