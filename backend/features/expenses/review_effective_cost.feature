@expense-registration @s4c @todo
Feature: Review the effective cost of an expense
  As a Financial Manager
  I want to see how much money really left my pocket for an expense
  So that I can understand the real cost after discounts and refunds

  The charged total is the amount the store charged. The effective cost is
  the charged total minus the reductions from other actors, such as a bank
  refund. A reduction that is only expected makes the effective cost an
  estimate until the reduction is received.

  Scenario: An expense without reductions
    Given an expense for 14000.00 "ARS" was registered
    When the Financial Manager reviews the expense
    Then the expense shows a charged total of 14000.00 "ARS"
    And the expense shows an effective cost of 14000.00 "ARS"
    And the effective cost is not marked as an estimate

  Scenario: An expense with an expected bank reduction
    Given an expense for 14000.00 "ARS" was registered
    And an expected bank reduction of 3500.00 "ARS" was recorded on the expense
    When the Financial Manager reviews the expense
    Then the expense shows a charged total of 14000.00 "ARS"
    And the expense shows an effective cost of 10500.00 "ARS"
    And the effective cost is marked as an estimate

  Scenario: An expense with a received bank reduction
    Given an expense for 14000.00 "ARS" was registered
    And a bank reduction of 3500.00 "ARS" was recorded and confirmed as received
    When the Financial Manager reviews the expense
    Then the expense shows an effective cost of 10500.00 "ARS"
    And the effective cost is not marked as an estimate

  Scenario: A receipt with a store discount and a bank promotion
    Given an expense for 59690.00 "ARS" was registered
    And an expected bank reduction of 11938.00 "ARS" was recorded on the expense
    When the Financial Manager reviews the expense
    Then the expense shows a charged total of 59690.00 "ARS"
    And the expense shows an effective cost of 47752.00 "ARS"

  Scenario: Reductions are listed by actor
    Given an expense with a store discount and a bank reduction was registered
    When the Financial Manager reviews the expense
    Then the expense lists each reduction with its actor, amount, description and status

  Scenario: A store discount is already part of the charged total
    Given an expense for 59690.00 "ARS" with a store discount of 1700.00 "ARS" was registered
    When the Financial Manager reviews the expense
    Then the expense shows a charged total of 59690.00 "ARS"
    And the expense shows an effective cost of 59690.00 "ARS"
