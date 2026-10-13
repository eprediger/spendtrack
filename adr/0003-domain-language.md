# 0003. Language of the domain terms

## Status

Accepted

## Context

* Domain-Driven Design says the ubiquitous language should be the language of the domain experts, and the same words should appear in conversation, documents, tests and code.
* The product is built in Argentina, by an Argentine developer, for an Argentine public for now. Many financial terms are clearer in Spanish.
* Some English terms do not map one to one to Spanish. Example: "Reduction" covers *descuento* (applied by the store), *reintegro* (paid back by the bank) and *ahorro* (the savings line on a receipt).
* [`adr/0001-backend-language.md`](0001-backend-language.md) says the project is also a portfolio piece, so an English-speaking audience matters.
* Many fiscal terms have an English equivalent (IVA is VAT, alícuota is rate, comprobante is receipt or invoice). A few are specific to the Argentine tax system and have no faithful equivalent: CAE and the receipt types (Factura A, B, C). CUIT is the Argentine kind of a concept that exists in every country, the Tax ID (NIT in Colombia, RUT in Chile, and so on).
* The change is cheapest now: only slice S1 is implemented.
* The plain-language rule in [`AGENTS.md`](../AGENTS.md) (STE100) exists for English only.

## Decision

English is the language of the domain terms, the documents, the `.feature` files and the code.

Spanish is an auxiliary language:

* Each core glossary term in [`product/glossary.md`](../product/glossary.md) has an **En español** line. It gives the Spanish word people use. It adds a note when the English term is ambiguous or does not map one to one.
* Terms from the Argentine fiscal receipt keep the Spanish text printed on the receipt as their source term.
* Fiscal concepts specific to the Argentine tax system, with no faithful English equivalent, keep the Spanish term as their name (for example CAE, Factura B). Fiscal terms that do have an English equivalent use it (for example VAT for IVA).
* A concept that exists in every country uses the English name. The Argentine kind keeps its local name as a kind of that concept. Example: the concept is Tax ID, and CUIT is the Argentine kind.
* The Spanish wording of terms created by the project is provisional until the product owner confirms it.

## Options considered

* **English everywhere, with Spanish auxiliary definitions (chosen)**: keeps code and documents consistent with the portfolio goal and the existing tooling. The glossary fixes the cases that are lost in translation.
* **Spanish for the domain model, English for technical layers**: closest to the DDD ideal for the domain experts. Rejected for now: it mixes the two languages in code (for example a Spanish entity with an English repository), and it conflicts with the portfolio goal.
* **Spanish everywhere**: the most consistent for local users. Rejected for now: the highest cost, and the same portfolio conflict.

## Consequences

* New glossary terms need an **En español** line. This is written in [`AGENTS.md`](../AGENTS.md).
* An English term that hides distinct Spanish concepts is a signal to revisit the concept itself. Example: "Reduction" may split into discount and refund (open decision in S4b).
* UI language is not part of this decision. A Spanish interface does not need Spanish domain code.
* Revisit if localisation becomes a feature, if the audience of the repository changes, or if the English terms cause repeated confusion with domain experts. The cost of the change grows with each implemented slice.
