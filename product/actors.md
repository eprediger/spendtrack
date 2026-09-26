# Actors

## Purpose

This document defines the relevant product actors and their responsibilities.

Actors describe who interacts with the product and why. They are not equivalent to authentication identities, technical users, permissions, or data ownership.

The same person may perform multiple roles.

---

## Financial Manager

A Financial Manager is a person authorized to register, manage, and review financial information within a Financial Context.

Typical responsibilities include:

* registering income;
* registering expenses;
* recording savings;
* recording investments;
* managing budgets;
* managing financial objectives;
* reviewing financial information;
* managing shared financial activity when authorized.

A Financial Manager may operate within:

* a personal Financial Context;
* a shared Financial Context.

Being a member of a shared Financial Context does not imply ownership of all financial resources represented within that context.

---

## Promotion Manager

A Promotion Manager is a person authorized to create, maintain, and manage commercial promotions.

Typical responsibilities include:

* registering promotions;
* defining promotion conditions;
* defining validity periods;
* associating promotions with Commerce, Establishments, Rubrics, or other applicable scopes;
* defining applicable payment methods;
* recording reimbursement or claim conditions;
* maintaining promotion information.

Promotion management is conceptually different from personal financial management.

---

## Price Analyst

A Price Analyst is a person or future business user authorized to analyze aggregated product, price, commerce, and consumption information.

A Price Analyst may need access to information such as:

* products;
* observed prices;
* effective prices;
* commerce;
* establishments;
* geographical information;
* historical price information;
* aggregated consumption patterns.

A Price Analyst should not necessarily have access to individual users' personal financial information or identifiable individual expenses.

The product must therefore distinguish between:

* access to aggregated analytical information;
* access to individual financial information.

---

## Identity, Role, Permission and Data Scope

The product must keep the following concepts separate:

### Identity

Represents the person or account interacting with the system.

### Role

Represents a functional responsibility performed by an identity.

A person may have multiple roles.

### Permission

Represents an authorized operation or capability available to a role.

### Data Scope

Represents which financial or analytical information an identity is authorized to access or modify.

Examples of data scope include:

* personal information;
* a specific shared Financial Context;
* aggregated market information;
* aggregated price information.

These concepts must not be collapsed into a single concept such as "user".

---

## Relationship With Financial Context

A Financial Context defines the scope in which financial information is managed.

A person may participate in one or more Financial Contexts.

Participation in a Financial Context does not by itself determine:

* economic ownership;
* custody of assets;
* responsibility for every expense;
* entitlement to every income;
* access to every piece of information.

Those concepts must be represented independently when required by the business rules.

---

## Multiple Roles

The product must allow the same person to perform different functional roles.

For example, the same identity could act as:

* Financial Manager for a personal Financial Context;
* Financial Manager for a shared family Financial Context;
* Promotion Manager for commercial information;
* Price Analyst for aggregated analytical information.

The role determines the type of responsibility being exercised.

The data scope determines which information is accessible within that responsibility.

---

## Product Principle

Authorization must be designed around:

**Identity → Role → Permission → Data Scope**

rather than assuming:

**Identity → unrestricted access to all related data**

This distinction is particularly important for shared finances and future price-analysis capabilities.

