@expense-registration @s4d @todo
Feature: Record a reduction from another actor
  As a Financial Manager
  I want to record a reduction that a bank or another actor applies after the purchase
  So that I can see the effective cost of the expense and follow the refund

  A reduction from another actor does not appear on the receipt. It starts as
  expected. The Financial Manager confirms it as received when the money
  arrives. The system has no bank integration, so the confirmation is manual.

  Scenario: Record an expected bank reduction
    Given an expense for 14000.00 "ARS" was registered
    When the Financial Manager records an expected reduction of 3500.00 "ARS" by the bank "The Bank" described as "25% refund at the gas station store" on the whole expense
    Then the reduction is recorded as expected
    And the reduction shows "The Bank" as the actor

  Scenario: Confirm a reduction as received
    Given an expected reduction of 3500.00 "ARS" was recorded on an expense
    When the Financial Manager confirms the reduction as received
    Then the reduction is recorded as received

  Scenario: Add a reduction after the expense was registered
    Given an expense was registered yesterday
    When the Financial Manager records an expected reduction on the expense
    Then the expense has the reduction

  Scenario Outline: Missing required information
    When the Financial Manager records a reduction without "<field>"
    Then the reduction is rejected
    And the system indicates "<field>" is required

    Examples:
      | field       |
      | actor       |
      | amount      |
      | description |
