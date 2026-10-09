@expense-registration @s1 @todo
Feature: Register a basic expense
  As a Financial Manager
  I want to register an expense as an economic event
  So that I can track where my money is consumed

  An expense represents economic consumption, not a payment.
  This slice covers the economic fact alone: no payment mechanism,
  no supporting document, no purchase detail, no installments.

  Required information: occurrence date and time, amount, currency, description.

  Scenario: Register an expense
    Given a Financial Context
    When the Financial Manager registers an expense at "2026-09-20T14:30:00-03:00" for 45.90 "USD" described as "Groceries"
    Then the expense is recorded in the Financial Context
    And the expense occurred at "2026-09-20T14:30:00-03:00" with amount 45.90 "USD" and description "Groceries"

  Scenario: Review a registered expense
    Given an expense at "2026-09-20T14:30:00-03:00" for 45.90 "USD" described as "Groceries" was registered
    When the Financial Manager reviews the expense
    Then the expense shows occurred at "2026-09-20T14:30:00-03:00" with amount 45.90 "USD" and description "Groceries"

  Scenario Outline: Missing required information
    When the Financial Manager registers an expense without "<field>"
    Then the registration is rejected
    And the system indicates "<field>" is required

    Examples:
      | field       |
      | occurred at |
      | amount      |
      | currency    |
      | description |

  Scenario Outline: Invalid information is rejected
    When the Financial Manager registers an expense at "<occurred at>" for <amount> "<currency>" described as "<description>"
    Then the registration is rejected
    And the system indicates the <field> is invalid

    Examples:
      | occurred at              | amount | currency | description | field    |
      | 2026-09-20T14:30:00-03:00| -45.90 | USD      | Groceries   | amount   |
      | 2026-09-20T14:30:00-03:00| 0      | USD      | Groceries   | amount   |
      | 2026-09-20T14:30:00-03:00| 45.90  | XX1      | Groceries   | currency |
      | 2026-09-20T14:30:00-03:00| 45.90  | US       | Groceries   | currency |

  Scenario: Occurrence cannot be in the future
    When the Financial Manager registers an expense that occurs in the future for 45.90 "USD" described as "Groceries"
    Then the registration is rejected
    And the system indicates the occurrence cannot be in the future

  Scenario: A missing field and an invalid field are reported together
    When the Financial Manager registers an expense with:
      | occurred at | 2026-09-20T14:30:00-03:00 |
      | currency    | XX1                       |
      | description | Groceries                 |
    Then the registration is rejected
    And the system indicates "amount" is required
    And the system indicates the currency is invalid

  Scenario: Reviewing an expense that does not exist
    When the Financial Manager reviews an expense with an unknown id
    Then the system indicates the expense was not found

  Scenario: Reviewing an expense with a malformed id
    When the Financial Manager reviews an expense with a malformed id
    Then the system indicates the id is invalid
