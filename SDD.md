# Spec-Driven Development

## Purpose

This project uses Spec-Driven Development (SDD) for AI-assisted software implementation.

The purpose is to ensure that AI coding agents implement approved behavior from explicit, reviewable specifications rather than from ad-hoc implementation prompts.

## Relationship With Product Development

Product Discovery remains responsible for defining business intent:

```text
Business Need
→ Epic
→ Feature
→ Actors
→ Story Map
→ User Story
→ Acceptance Criteria
→ BDD
```

SDD begins when an implementation slice has sufficiently defined and approved behavior:

```text
Approved Requirements
→ Specification
→ Clarification
→ Technical Plan
→ Tasks
→ Implementation
→ Verification
→ Validation
```

SDD does not replace Product Discovery.

## Principles

* Specifications must be traceable to approved product requirements.
* AI agents must not invent unresolved business decisions.
* Technical decisions must be derived from validated requirements.
* Implementation tasks must derive from the approved technical plan.
* Changes discovered during implementation must be reflected in the appropriate artifact.
* Implementation must be verified against the specification.
* Human approval remains required for significant product, architecture, scope, and risk decisions.

## Artifact Responsibilities

| Artifact       | Purpose                                | Source of Truth                  |
| -------------- | -------------------------------------- | -------------------------------- |
| `.feature` file | User Story + Acceptance Criteria + BDD | Product behavior                |
| `specs/`       | Detailed implementation specification  | Approved implementation behavior |
| Technical Plan | Technical design for the specification | Engineering approach             |
| Tasks          | Executable implementation units        | Approved technical plan          |
| ADR            | Significant architectural decision     | Architecture                     |
| `src/`         | Implementation                         | Executable system                |

The exact relationship between these artifacts remains subject to refinement.

## Repository Conventions

* Executable requirements (`.feature`): `backend/features/<domain>/<story>.feature`. Unimplemented stories and scenarios are tagged `@todo`; slice traceability uses tags like `@s1`.
* Specifications: `specs/<feature>/<slice>/` — `spec.md` (approved implementation behavior), `plan.md` (technical design), `tasks.md` (implementation checklist).
* Architectural decisions: `adr/` (created on the first recorded decision).
* Domain model documentation: `domain/`.
* Approval is a human decision, recorded in the relevant `decisions.md` or in the artifact's own status.
* An implementation slice is validated when its scenarios run green with `@todo` removed and the activated domain terms are marked `validated` in `product/glossary.md`.

## Spec Kit

Spec Kit is the proposed tooling for operationalizing SDD.

Its formal adoption and integration are not yet final decisions.

Before adoption, the project must determine:

* supported AI coding agents;
* Spec Kit version;
* installation strategy;
* repository integration;
* specification location;
* relationship between Spec Kit artifacts and repository `.feature` files and `specs/`;
* relationship between Spec Kit tasks and `tasks.md` (and optional GitHub Issues);
* required workflow stages;
* quality gates;
* relationship with ADRs and architecture documentation.

## Proposed Pilot

The first SDD pilot should use a small, approved implementation slice of `Expense Registration`.

The pilot should evaluate:

* requirement traceability;
* specification quality;
* AI implementation behavior;
* ambiguity detection;
* duplication between artifacts;
* developer experience;
* reviewability;
* maintenance cost;
* integration with GitHub;
* verification quality.

The objective of the pilot is to validate the development process, not merely the generated code.

## Next Steps

Resolved: the SDD lifecycle and artifact locations are defined in this document; GitHub Issues are an optional workflow layer holding no requirement content; project-level AI engineering principles live in `AGENTS.md`.

1. Draft the specification for the first slice of `Expense Registration` (S1).
2. Run the pilot: specification → technical plan → tasks → implementation → validation.
3. Evaluate the pilot against the criteria listed above.
4. Evaluate Spec Kit against this lifecycle — adopt only if it adds value over the conventions already defined.
5. Adjust the workflow and adopt it for subsequent Features if approved.

