# AI Engineering Instructions

## Project Context

This is a greenfield financial management product for individuals and families.

The product enables users to manage and understand:

* income;
* expenses;
* savings;
* investments;
* shared financial contexts;
* financial objectives;
* budgets;
* future financial commitments.

The product is developed iteratively. Business requirements must be understood and approved before technical implementation.

## Engineering Principles

* Start from the business problem, not from technology.
* Do not assume architecture, frameworks, languages, infrastructure, or implementation strategies before requirements justify them.
* Prefer the simplest solution that satisfies validated requirements.
* Avoid premature optimization and speculative abstractions.
* Preserve explicit distinctions between business concepts.
* Do not invent unresolved business decisions.
* Expected behavior must be documented and approved — as User Story + Acceptance Criteria + BDD, and later as specification — before implementation begins. Code that implements undocumented or unapproved behavior is out of process.
* Identify ambiguity and request clarification when it materially affects behavior.
* The human is the final decision maker for significant product, architecture, scope, and risk decisions.

## Product-to-Implementation Flow

```text
Business Need
→ Epic
→ Feature
→ Actors
→ Story Map
→ User Story + Acceptance Criteria + BDD (executable .feature file)
→ Domain Discovery
→ Domain Model
→ Technical Design
→ Implementation
→ Validation
→ Iterate
```

User Stories, Acceptance Criteria, and BDD are treated as one requirements unit.

Executable `.feature` files are the source of truth for User Stories, Acceptance Criteria, and BDD: the Feature description carries the narrative, scenarios carry the Acceptance Criteria, and step bindings make them executable.

`.feature` files currently live in `backend/features/`. Their placement relative to future applications is deferred until a second application needs to bind the same scenarios (last responsible moment).

GitHub Issues and Projects are an optional workflow layer for status, prioritization, and discussion. They hold no requirement content and link to the corresponding `.feature` file.

The repository stores durable product, domain, architecture, and engineering knowledge.

Do not duplicate User Story or Acceptance Criteria content outside `.feature` files.

## Spec-Driven Development

The project adopts Spec-Driven Development (SDD) for AI-assisted implementation.

SDD governs the transition from approved product requirements to implementation.

```text
Approved Requirements
→ Specification
→ Technical Plan
→ Implementation Tasks
→ Implementation
→ Verification
→ Validation
```

The SDD process must:

* make implementation intent explicit;
* prevent AI agents from inventing unresolved requirements;
* keep specifications traceable to approved requirements;
* derive technical plans from approved specifications;
* derive implementation tasks from approved plans;
* validate implementation against the specification;
* propagate approved changes back to the appropriate artifacts.

SDD does not replace Product Discovery, Story Mapping, User Stories, Acceptance Criteria, BDD, or human approval.

## Proposed SDD Tooling

Spec Kit is the proposed tooling for operationalizing SDD.

Its adoption and exact integration into this repository are still open decisions.

Do not assume specific Spec Kit commands, workflows, artifact locations, or GitHub integrations until those decisions are explicitly approved.

## Decision Handling

Distinguish clearly between:

* approved decisions;
* assumptions;
* recommendations;
* open decisions.

Never represent an open decision as an approved requirement.

Significant decisions require explicit human approval.

## Repository Context

Relevant project knowledge should be loaded progressively:

```text
product/
    Product intent and requirements

domain/
    Business domain model and rules

specs/
    Approved implementation specifications

architecture/
    Technical architecture

adr/
    Significant architectural decisions

src/
    Implementation
```

`AGENTS.md` contains instructions for AI agents. It must not become a duplicate repository of all product knowledge.

## Quality

AI-generated code is subject to the same engineering standards as human-written code.

Consider, according to feature risk and value:

* correctness;
* automated testing;
* security;
* accessibility;
* performance;
* reliability;
* observability;
* maintainability.

Implementation is not complete merely because the code compiles or tests pass. The approved behavior must also be validated.

## Domain Integrity

The following distinctions are fundamental:

```text
Expense ≠ Movement
Purchase ≠ Expense
Financial Context ≠ Ownership
Custody ≠ Economic Ownership
Obligation ≠ Expense
Payment Commitment ≠ Expense
Savings ≠ Investment
Objective ≠ Asset/Position
Promotion ≠ Ordinary Income
```

Do not collapse these concepts merely for implementation convenience.
