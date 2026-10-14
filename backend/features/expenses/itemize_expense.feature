@expense-registration @s4a @todo
Feature: Itemize an expense
  As a Financial Manager
  I want to break an expense down into items
  So that I can see what I bought and how much each item cost

  An item has a description, a quantity, a unit price and a line total.
  Items are optional: an expense without items stays valid.
  The quantity is a positive decimal number: a weight is a fraction.
  The unit price and the line total keep the sign the receipt prints.
  The line total is the amount the receipt prints. It can differ from
  quantity times unit price. Items are in the currency of the expense.

  Scenario: Register an expense with items
    Given a Financial Context
    When the Financial Manager registers an expense at "2026-10-05T16:53:34-03:00" for 14000.00 "ARS" described as "Gas station store" with items:
      | description   | quantity | unit price | line total |
      | Candies       | 5        | 150.00     | 750.00     |
      | Bottled water | 1        | 3400.00    | 3400.00    |
      | Sandwich      | 1        | 9850.00    | 9850.00    |
    Then the expense is recorded in the Financial Context
    And the expense has 3 items

  Scenario: Register an expense without items
    Given a Financial Context
    When the Financial Manager registers an expense at "2026-09-20T14:30:00-03:00" for 45.90 "USD" described as "Groceries"
    Then the expense is recorded in the Financial Context
    And the expense has no items

  Scenario: The printed line total differs from quantity times unit price
    When the Financial Manager registers an expense with an item of quantity 3, unit price 33.33 and line total 100.00
    Then the expense is recorded
    And the item keeps the line total 100.00 as printed

  Scenario: A returned deposit is a negative item
    When the Financial Manager registers an expense with an item "RETURNABLE BOTTLE" of quantity 3, unit price -3001.00 and line total -9003.00
    Then the expense is recorded
    And the item shows quantity 3 and line total -9003.00

  Scenario: A quantity can be a fraction
    When the Financial Manager registers an expense with an item of quantity 0.750, unit price 3200.00 and line total 2400.00
    Then the expense is recorded
    And the item shows quantity 0.750 and line total 2400.00

  Scenario Outline: A quantity that is not positive is rejected
    When the Financial Manager registers an expense with an item of quantity <quantity>, unit price 150.00 and line total 150.00
    Then the registration is rejected
    And the system indicates the item "quantity" is invalid

    Examples:
      | quantity |
      | 0        |
      | -1       |

  Scenario Outline: Incomplete item information is rejected
    When the Financial Manager registers an expense with an item without "<field>"
    Then the registration is rejected
    And the system indicates the item "<field>" is required

    Examples:
      | field       |
      | description |
      | quantity    |
      | unit price  |
      | line total  |

  Scenario: Review the items of a registered expense
    Given an expense with items was registered
    When the Financial Manager reviews the expense
    Then the expense shows each item with its description, quantity, unit price and line total
