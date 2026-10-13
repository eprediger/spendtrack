@expense-registration @s4f @todo
Feature: Correct a reduction that was applied wrongly
  As a Financial Manager
  I want to correct a reduction when the real outcome differs from what was expected
  So that the effective cost shows what really happened

  The system has no bank integration and can suggest a reduction that does
  not apply. The Financial Manager decides the real outcome.
  The record keeps the expected amount next to the outcome.

  Scenario: The received amount differs from the expected amount
    Given an expected reduction of 3500.00 "ARS" was recorded on an expense
    When the Financial Manager confirms the reduction as received for 2800.00 "ARS"
    Then the reduction shows an expected amount of 3500.00 "ARS"
    And the reduction shows a received amount of 2800.00 "ARS"
    And the effective cost uses 2800.00 "ARS"

  Scenario: The reduction never arrives
    Given an expected reduction of 3500.00 "ARS" was recorded on an expense
    When the Financial Manager marks the reduction as not received
    Then the reduction shows an expected amount of 3500.00 "ARS"
    And the reduction shows it was not received
    And the effective cost does not include the reduction

  Scenario: The reduction never applied
    Given an expected reduction was recorded on an expense
    When the Financial Manager dismisses the reduction
    Then the reduction no longer reduces the effective cost
