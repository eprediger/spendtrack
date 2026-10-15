# Domain Glossary

## Term Status

Every term carries an implementation status:

* **provisional** — the concept is defined, but no approved implementation slice has activated it yet.
* **validated** — the concept's behavior is implemented and verified against an approved specification.

All terms begin as provisional and become validated as approved User Stories activate them through implementation slices.

Approval status (firm agreement vs. open decision) is a separate axis, tracked in each Feature's `decisions.md`. Where a term is backed by a Firm Agreement, this is noted on the term.

## Language

Terms and definitions are in English, as are code and documents ([`adr/0003-domain-language.md`](../adr/0003-domain-language.md)). Each core term has an **En español** line: the Spanish word people use, and a note when the English term does not map one to one. Terms from the Argentine fiscal receipt keep the Spanish text printed on the receipt.

When the Spanish line says the wording is provisional, it needs confirmation by the product owner.

## Expense

**Status:** validated by S1 (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

Economic consumption recognized as a financial fact.

An expense is not synonymous with a cash outflow or payment.

**En español:** gasto. Es el hecho económico (el consumo), no el pago. La salida de dinero es un *movimiento*.

## Movement

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

A financial movement or settlement involving financial resources.

A movement is not necessarily an expense.

**En español:** movimiento (de fondos). Puede ser un pago, una transferencia o una extracción. No siempre es un gasto.

## Purchase

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

A detailed representation of an expense resulting from a commercial transaction.

A purchase does not create a second independent expense.

**En español:** compra. Es el detalle de un gasto que viene de una transacción comercial. No crea un segundo gasto.

## Payment Commitment

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

A future commitment to settle an amount related to an expense or other obligation.

A payment commitment is not itself the economic consumption.

**En español:** compromiso de pago (por ejemplo, las cuotas de una compra). No es el gasto.

## Obligation

**Status:** provisional

A future or existing financial responsibility that may require settlement.

An obligation is distinct from the expense that may have originated it.

**En español:** obligación (deuda o responsabilidad financiera que puede requerir un pago). No es el gasto que la originó.

## Financial Context

**Status:** provisional — activated implicitly by S1 (single context, not yet persisted); firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md)

An environment in which financial information is managed and shared among one or more people.

Membership in a Financial Context does not imply economic ownership.

**En español:** contexto financiero (por ejemplo, personal o familiar). Participar en uno no implica ser dueño económico de todo lo que contiene.

## Economic Ownership

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

The person or persons to whom an asset, position, income, expense, or economic responsibility belongs.

Economic ownership is distinct from account custody.

**En español:** titularidad económica (a quién le pertenece de verdad). No es lo mismo que quién tiene la cuenta (*custodia*).

## Custody

**Status:** provisional

The financial account or institution through which an asset or position is held.

The custodian does not necessarily represent the economic owner.

**En español:** custodia (la cuenta o la institución donde está el activo). Quien custodia no es necesariamente el dueño.

## Financial Position

**Status:** provisional

An economic ownership position in an asset.

A financial position may be distributed among multiple economic owners.

**En español:** posición financiera (la parte de un activo que tiene cada titular económico).

## Asset

**Status:** provisional

An economic resource with value.

**En español:** activo (recurso con valor económico).

## Liability

**Status:** provisional

An economic responsibility or amount owed.

**En español:** pasivo (monto que se debe).

## Net Worth

**Status:** provisional

Derived economic value calculated from economic assets/positions minus relevant liabilities and obligations.

Net worth is not a primary historical fact.

**En español:** patrimonio neto (activos menos pasivos). Se calcula; no es un hecho histórico.

## Savings

**Status:** provisional

Accumulation of financial resources that are not necessarily invested.

Savings and investments are distinct concepts.

**En español:** ahorros (dinero apartado, no necesariamente invertido). No confundir con la línea "Ahorro" de los tickets, que muestra los descuentos obtenidos en esa compra.

## Investment

**Status:** provisional

A financial asset or position acquired or maintained with an investment purpose.

**En español:** inversión.

## Objective

**Status:** provisional

A financial goal with a target and potentially a time horizon.

An objective is distinct from the asset or position used to fund it.

**En español:** objetivo (o meta) financiero. No es lo mismo que el activo que lo financia.

## Budget

**Status:** provisional

A set of financial targets or allocation rules used to guide and evaluate financial behavior.

**En español:** presupuesto.

## Commerce

**Status:** provisional

A business or brand associated with a commercial transaction.

**En español:** comercio (el negocio o la marca, por ejemplo una cadena de supermercados). El ticket muestra la razón social, que puede ser distinta de la marca (ver el [ticket de la estación de servicio](examples/gas-station-receipt.md)).

## Establishment

**Status:** provisional

A physical location where a commercial transaction occurs.

An establishment is recorded with its exact address, so that its coordinates (latitude and longitude) can be found later.

**En español:** establecimiento (el lugar físico: local o sucursal), con su dirección exacta.

## Product

**Status:** provisional

A product or service acquired in a purchase.

**En español:** producto (o servicio).

## Purchase Item

**Status:** provisional

One line of a purchase. A purchase item has a description, a quantity, a unit price and a line total. The line total is the printed amount. It can differ from quantity times unit price. A purchase item can be negative (for example a returned deposit).

**En español:** ítem o renglón (cada línea de un ticket o de una compra).

## Charged Total

**Status:** provisional

The amount the store charged for an expense, after the store's own discounts. It is the amount registered on the expense.

**En español:** total cobrado (lo que el comercio cobró; el "TOTAL" del ticket). No usar "total pagado": en este producto, lo que realmente se pagó después de los reintegros es el *costo de bolsillo*.

## Reduction

**Status:** provisional

A decrease of the cost of an expense, applied by an actor (the store, a bank, or another). It applies to the whole expense or to one item. A reduction from another actor is expected until the Financial Manager confirms it as received.

A reduction is not ordinary income.

**En español:** no tiene una sola palabra en español, y "reducción" no se usa así en la vida diaria. Es el monto que se aplica a un gasto, no la regla de la promoción (el porcentaje y el tope pertenecen a la *promoción*). Según cómo llega se dice distinto: *descuento* (se resta en la caja y figura en el ticket, por ejemplo "DESCUENTOS: -1700") o *reintegro* (el dinero que se devuelve después, con tope). El ticket también tiene una línea de *ahorro*. Ojo: "descuento" también se usa para el porcentaje de una promoción. Redacción provisional, a confirmar. Puede que haya que separar estos conceptos (ver decisiones abiertas de S4b).

## Effective Cost

**Status:** provisional

The money that really left the pocket for an expense: the charged total minus the reductions from other actors. It is an estimate while a reduction is only expected.

**En español:** costo de bolsillo: el monto neto que se terminó pagando, luego de los reintegros. Es una estimación mientras algún reintegro esté pendiente. Se evitan "costo efectivo" (*efectivo* también significa dinero en efectivo), "costo real" (suele significar ajustado por inflación) y "costo neto" (en las facturas, *neto* suele significar antes de IVA).

## Promotion

**Status:** provisional

A commercial condition that can provide a discount, benefit, or reimbursement when its conditions are satisfied.

A promotion is the rule. A reduction is what the rule produces for one expense (for example a 25% discount rate with a monthly refund cap of 10,000 gives a refund of 3,500 on a 14,000 expense, and 10,000 on a 50,000 expense).

A promotion can have a discount rate (a percentage), a refund cap, a minimum purchase, valid days, a validity period, payment methods, a payment channel and a store scope. See [Promotion Management](epics/promotion-management/epic.md).

**En español:** promoción o beneficio (por ejemplo, el reintegro de un banco por pagar con QR, ver la [promoción del supermercado](examples/bank-promotion-supermarket.md)). Es la regla comercial, no el monto que se aplica a un gasto. Tiene un *porcentaje de descuento* (por ejemplo 25%) y un *tope de reintegro* (por ejemplo $10.000 por mes).

## Reimbursement

**Status:** provisional

A later financial benefit associated with a previous purchase or promotion. It starts as expected and is confirmed as received.

A reimbursement is not ordinary income.

**En español:** reintegro (la palabra que usa el banco). Empieza como esperado y se confirma al recibirlo. No es un ingreso.

## Forecast

**Status:** provisional

A derived expectation about future financial behavior.

A forecast must not modify historical facts.

**En español:** proyección (o pronóstico). Nunca modifica los hechos históricos.

## Recurring Rule

**Status:** provisional

A rule describing how an obligation or financial event may recur.

A recurring rule is distinct from a concrete occurrence.

**En español:** regla de recurrencia (por ejemplo, un pago mensual). No es lo mismo que cada ocurrencia concreta.

## Document

**Status:** provisional (firm agreement — see expense-registration [decisions.md](epics/financial-management/features/expense-registration/decisions.md))

Evidence associated with a financial event, such as a receipt, invoice, image, or PDF.

**En español:** documento o comprobante (ticket, factura, imagen, PDF).

---

## Argentine Fiscal Receipt Terms

**Status:** provisional (all terms in this section)

These terms describe the structure of an Argentine fiscal receipt ("ticket fiscal"). The receipt is a [Document](#document) that can feed an expense. Examples: the [gas station](examples/gas-station-receipt.md), [supermarket](examples/supermarket-receipt.md) and [bakery](examples/bakery-receipt.md) receipts (see [`examples/`](examples/README.md)), used in [Expense Registration](epics/financial-management/features/expense-registration/decisions.md).

How to read the tables:

* **Source term** is the text printed on the receipt (Spanish).
* **Reference key** is the identifier used in the original receipt glossary. It keeps the terms traceable. It is not an approved data model: which fields the product captures, and when, is decided per slice (see the open decisions about fiscal receipt fields and commerce details in the Expense Registration decisions).
* Not all fields are equally important or needed now.
* The receipt format is the most common one and the first to support. Other formats (utility bills, credit card statements) will come later.

Related terms above: Commerce and Establishment (merchant data), Purchase Item (line items), Charged Total (`total_amount`), Reduction (discounts), Document (the receipt itself).

### Merchant Data (Datos del Emisor)

| Term | Source term | Reference key | Meaning |
| ---- | ----------- | ------------- | ------- |
| Trade name | Nombre Comercial / Fantasía | `trade_name` | The brand name shown to customers (for example the name of a supermarket chain). |
| Legal name | Razón Social | `legal_name` | The official company name registered with the tax authority (for example the registered company behind a supermarket chain or a gas station complex). |
| Tax ID | C.U.I.T. | `tax_id` | The number the tax authority assigns to a person or a legal entity, to identify it for tax purposes. It is unique within the country that issues it. CUIT is the Argentine kind: a unique 11-digit federal taxpayer number (other countries use other names, for example NIT in Colombia). It is not the same as an identity document such as the DNI. On a receipt it identifies the legal entity of the merchant, not the brand or the store. |
| Gross income tax ID | Ingresos Brutos | `gross_income_tax_id` | The provincial gross income tax registration number. |
| Address | Domicilio / Dirección | `address` | The physical address of the store branch or establishment. |
| Tax status | Condición frente al IVA | `tax_status` | The tax profile of the merchant (for example "IVA Responsable Inscripto" means tax registered). |
| Business start date | Inicio de Actividades | `business_commencement_date` | The date the company legally started to operate. |

### Receipt Metadata (Datos del Comprobante)

| Term | Source term | Reference key | Meaning |
| ---- | ----------- | ------------- | ------- |
| Receipt type | Tipo de Comprobante | `receipt_type` | The tax document category set by the tax authority (AFIP), for example "Tique Factura 'B'" or "Factura B Original". |
| Receipt number | N° de Comprobante / Factura | `receipt_number` | The unique serial identifier of the transaction. |
| Purchase date | Fecha | `purchase_date` | The calendar date of the transaction (DD/MM/YYYY). |
| Purchase time | Hora | `purchase_time` | The time of the payment (HH:MM:SS). |
| Store number | Nro. de Tienda / Sucursal | `store_number` | The internal identifier of the physical branch. |
| POS number | Nro. de Caja | `pos_number` | The identifier of the point-of-sale terminal or checkout register. |
| Cashier | Cajero / Operador | `cashier` | The name or employee ID of the cashier. |

### Customer Data (Datos del Cliente)

| Term | Source term | Reference key | Meaning |
| ---- | ----------- | ------------- | ------- |
| Customer tax category | Condición de Compra / Categoría Fiscal | `customer_tax_category` | The tax category of the buyer (for example "Consumidor Final" means end consumer). |

### Line Items (Detalle de Productos)

| Term | Source term | Reference key | Meaning |
| ---- | ----------- | ------------- | ------- |
| Quantity | Cantidad | `quantity` | The number of units bought on a line. |
| Unit price | Precio Unitario | `unit_price` | The price of one unit. Whether it includes tax depends on regional rules. |
| Description | Descripción / Ítem | `description` | The commercial name, retail description or barcode of the product. |
| VAT rate | Alícuota IVA | `vat_rate` | The percentage of value added tax on the line (for example 21.00%, 10.50%). |
| Line total | Importe / Subtotal Ítem | `line_total` | The total cost of the line (quantity times unit price). |

### Totals, Discounts and Payments (Totales y Descuentos)

| Term | Source term | Reference key | Meaning |
| ---- | ----------- | ------------- | ------- |
| Subtotal | Subtotal | `subtotal` | The sum of all lines before general discounts. |
| Discounts | Descuentos / Ahorro | `discounts` | The sum of deductions, loyalty discounts or promotion reductions. |
| Total amount | TOTAL | `total_amount` | The final net amount charged to the customer. |
| Payment method | Medio de Pago | `payment_method` | The financial instrument used to pay (for example "Tarjeta" means card, "Efectivo" means cash). |
| Amount paid | Suma de sus pagos / Recibí | `amount_paid` | The total amount the customer handed over or authorized. |
| Change due | Su Vuelto | `change_due` | The change returned to the customer when cash exceeds the total amount. |

The payment lines (payment method, amount paid, change due) are not on every receipt. For example, the [bakery receipt](examples/bakery-receipt.md) prints none of them. For now, an expense is paid with one payment method (see S6 in the Expense Registration decisions). Several payment methods on one expense, which the line "Suma de sus pagos" suggests, is a feature to explore later.

On the [supermarket receipt](examples/supermarket-receipt.md), "DESCUENTOS" is the sum of the store's own promotions, applied before the payment processor handles the transaction. The label is how that store shows the sum. A promotion that depends on the payment (for example the [bank promotion](examples/bank-promotion-supermarket.md) for paying with a credit card through a QR) is not in that sum.

Whether informational savings printed in this block (for example the loyalty club savings line of the [supermarket receipt](examples/supermarket-receipt.md), already inside prices) are recorded is an open decision (see S4b in the Expense Registration decisions).

### Tax Breakdown (Desglose de Impuestos)

| Term | Source term | Reference key | Meaning |
| ---- | ----------- | ------------- | ------- |
| Fiscal receipt | Ticket Fiscal | `fiscal_receipt` | The generic name for a government-regulated official receipt. |
| VAT included | IVA Contenido | `vat_included` | The part of the total amount that is value added tax. |
| Other indirect taxes | Otros Impuestos / Tributos | `other_indirect_taxes` | Additional municipal fees, internal taxes or national indirect taxes. |

### Fiscal Validation (Validación Fiscal)

| Term | Source term | Reference key | Meaning |
| ---- | ----------- | ------------- | ------- |
| CAE | CAE | `cae` | "Código de Autorización Electrónico". The mandatory authorization code issued by AFIP. |
| CAE expiration date | F. vto (Vencimiento CAE) | `cae_expiration_date` | The deadline that sets the legal validity of the CAE. |
| Fiscal QR code | Código QR | `fiscal_qr_code` | The 2D barcode with the fiscal data needed for public validation. |
