@expense-registration @s4b
Feature: Record store discounts
  As a Financial Manager
  I want to record the discounts the store applied on the receipt
  So that I can see how much the store reduced each purchase

  S4b records a discount amount already applied on the receipt and associated
  with one purchase item. It does not record the promotion rule that produced
  the discount. Discounts applied to the whole expense are outside this slice.

  @todo
  Scenario: Record a store discount on one item
    Given an itemized expense with charged total 59690.00 "ARS", item "COCA COLA R 2L" line total 12750.00 "ARS" and other items totaling 48640.00 "ARS" was registered
    When the Financial Manager records a store discount of 1700.00 "ARS" described as "2DO AL 40% LINEA COC" on the item "COCA COLA R 2L" identified by its item id
    Then the discount is recorded on the item "COCA COLA R 2L"

  @todo
  Scenario: Review the store discounts of an expense
    Given an expense with a store discount was registered
    When the Financial Manager reviews the expense
    Then the expense shows the discount amount and description on the item it applies to

  @todo
  Scenario: Reject an item discount when the itemized expense does not reconcile
    Given an itemized expense with a charged total of 95.00 "ARS" and item "SOFT DRINK 2L" line total of 100.00 "ARS" was registered
    When the Financial Manager records item discounts totaling 3.00 "ARS" on that item
    Then the discounts are rejected because the itemized expense does not reconcile

  @todo
  Scenario: Record multiple discounts on one item together
    Given an itemized expense with a charged total of 95.00 "ARS" and an item "SOFT DRINK 2L" line total of 100.00 "ARS" was registered
    When the Financial Manager records item discounts of 3.00 and 2.00 "ARS" on that item together
    Then both discounts are recorded on the item "SOFT DRINK 2L"

  @todo
  Scenario: Reject invalid replacement of item discounts
    Given an itemized expense with a charged total of 95.00 "ARS" and item "SOFT DRINK 2L" line total of 100.00 "ARS" was registered
    And item discounts totaling 5.00 "ARS" were recorded on that item
    When the Financial Manager replaces that item's discounts with 4.00 "ARS"
    Then the discount replacement is rejected because the itemized expense does not reconcile
    And the previously recorded discounts remain unchanged

  @todo
  Scenario: Replace an item discount set without duplicating entries
    Given an itemized expense with a charged total of 95.00 "ARS" and an item "SOFT DRINK 2L" line total of 100.00 "ARS" was registered
    And item discounts of 3.00 and 2.00 "ARS" were recorded on that item
    When the Financial Manager replaces the item discounts with 4.00 and 1.00 "ARS"
    Then the item has exactly discounts of 4.00 and 1.00 "ARS"
    When the Financial Manager submits the same item discount set again
    Then the item still has exactly 2 discounts

  @todo
  Scenario: Reject discounts that exceed the item's line total
    Given an itemized expense with a charged total of 40.00 "ARS" and item "SOFT DRINK 2L" line total of 100.00 "ARS" and another item line total of 50.00 "ARS" was registered
    When the Financial Manager records discounts totaling 110.00 "ARS" on the item "SOFT DRINK 2L"
    Then the discounts are rejected because they exceed that item's line total
