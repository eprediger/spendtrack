# 0002. API protocol for clients

## Status

Accepted

## Context

* A mobile app is a real roadmap item, raising the question of whether gRPC should be the client-facing protocol (binary protobuf, HTTP/2, lower serialization cost on constrained devices).
* The backend follows hexagonal architecture: transport protocols are driving adapters and do not touch the domain.

## Decision

REST + JSON is the client-facing protocol for now.

## Rationale

* At this product's scale (small payloads, few calls), gRPC's CPU/battery advantages are real in principle but negligible in practice — battery cost is dominated by radio usage, not JSON parsing.
* REST keeps free tooling: OpenAPI docs, `curl` debugging, no proxy needed if a web client ever exists.
* gRPC adds a protoc/codegen toolchain per client platform and removes easy introspection.

## Consequences

* gRPC can be added later as a second driving adapter alongside HTTP without touching the domain — this decision is deliberately reversible.
* If gRPC is adopted for learning/portfolio reasons, that should be recorded as the driver rather than technical necessity.
