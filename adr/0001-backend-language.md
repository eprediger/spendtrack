# 0001. Backend language and framework

## Status

Accepted (supersedes an earlier Python decision reached during the same evaluation)

## Context

* Product goal: a personal financial management tool, production-grade, also serving as a portfolio piece.
* Technical goal: deliberately outside the Node/TypeScript/React ecosystem.
* A mobile app is a real roadmap item (register an expense on the spot, consult the budget anywhere). Clients communicate with the backend over an API.
* The product vision includes aggregated price analysis (Price Analyst actor), which favors a strong data ecosystem.
* Single developer; portfolio signaling toward backend/infrastructure roles is a stated goal.
* The backend had a Python 3.13 + FastAPI scaffold (hexagonal architecture, pytest, behave, Docker, CI) with only a health endpoint implemented — the cheapest possible moment to change language.

## Decision

**Go** is the backend language for the core product API.

**Python** remains reserved for data-analysis workloads (the Price Analyst capability and any analytics/reporting features), where its ecosystem (pandas, notebooks) is unmatched.

Client applications are built in the best-fit technology for their platform (e.g., native Kotlin for Android). Backend and clients communicate through an API and share no code.

## Options considered

* **Go** (chosen): strong portfolio signal for backend/infrastructure roles; single static binary deployment; excellent for API services and concurrent I/O; the domain-modeling weaknesses are mitigated by convention — unexported fields + validating constructors, marker interfaces for closed sets, and the `exhaustive` linter for switch coverage.
* **Python + FastAPI**: fast iteration and the best data ecosystem; weaker domain-invariant enforcement (discipline over compiler); the existing scaffold made it the path of least resistance, which is not a decision driver. Retained for analytics only.
* **Kotlin (Ktor/Spring)**: strongest domain modeling of the candidates (sealed/value classes, exhaustive `when`) and Kotlin Multiplatform sharing with a native app. Rejected because backend and clients do not need to share a language, and JVM ceremony outweighs the benefit at this stage. Revisit if a shared domain module proves valuable.
* **Rust**: rejected — correctness guarantees cost a development-velocity tax this API does not need.
* **C#/.NET**: viable, but offered no decisive advantage.

## Consequences

* Domain invariants are enforced by convention and tooling (constructors, marker interfaces, linters, tests), not by the compiler. This is the accepted cost of the Go choice.
* The Python scaffold in `backend/` is replaced by a Go implementation; `backend/features/*.feature` files are language-agnostic and carry over — BDD is re-bound via `godog`.
* CI, Dockerfile, Makefile, and hooks are migrated to Go tooling.
* Python remains in the monorepo for future analytics services.
