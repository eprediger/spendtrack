@expense-registration @s4e @todo
Feature: Identify the commerce and establishment of an expense
  As a Financial Manager
  I want to say where an expense happened
  So that I can see which businesses and places my money goes to

  Commerce is the business or brand. Establishment is the physical place.
  Both are optional. An establishment keeps its exact address, so that its
  coordinates (latitude and longitude) can be found later. The rest of the
  information to capture is still to be refined (legal entity, brand, branch).

  Scenario: Identify the commerce and establishment
    Given a Financial Context
    When the Financial Manager registers an expense at the commerce "Supermarket Chain" and the establishment "Main Avenue 1234, Buenos Aires"
    Then the expense is recorded with the commerce "Supermarket Chain"
    And the expense is recorded with the establishment "Main Avenue 1234, Buenos Aires"

  Scenario: An establishment on a highway keeps its exact address
    Given a Financial Context
    When the Financial Manager registers an expense at the establishment "National Highway 99, km 100, Example City, Example Province"
    Then the expense is recorded with the establishment "National Highway 99, km 100, Example City, Example Province"

  Scenario: Register an expense without commerce or establishment
    Given a Financial Context
    When the Financial Manager registers an expense without commerce and without establishment
    Then the expense is recorded
