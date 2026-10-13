@expense-registration @s4b @todo
Feature: Record store discounts
  As a Financial Manager
  I want to record the discounts the store applied on the receipt
  So that I can see how much the store reduced each purchase

  A store discount is a reduction applied by the store and printed on the
  receipt. It applies to the whole expense or to one item.

  Scenario: Record a store discount on one item
    Given an expense with an item "SOFT DRINK 2L" was registered
    When the Financial Manager records a store discount of 1700.00 "ARS" described as "SECOND UNIT AT 40%" on the item "SOFT DRINK 2L"
    Then the discount is recorded on the item "SOFT DRINK 2L"
    And the discount shows the store as the actor

  Scenario: Record a store discount on the whole expense
    Given an expense was registered
    When the Financial Manager records a store discount of 500.00 "ARS" described as "Store promotion" on the whole expense
    Then the discount is recorded on the whole expense
    And the discount shows the store as the actor

  Scenario: Review the store discounts of an expense
    Given an expense with a store discount was registered
    When the Financial Manager reviews the expense
    Then the expense shows the discount, its amount, its description, what it applies to and the store as the actor
