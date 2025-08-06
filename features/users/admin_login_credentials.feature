@todo
Feature: Administrator Login
  As an Administrator
  I want to log into my account on the platform
  So that I can access the system functionalities

  Background:
    Given the admin email "admin@domain.com" and password "password"
    And the registration is successful

  Scenario: Successful login
    When the administrator attempts to log in is correct
    Then the system should allow access to the system

  Scenario: Failed login
    When the administrator attempts to log in is incorrect
    Then an appropriate error message should be shown to the user

  Scenario: Login failure
    Given the user login fails due to a service error
    When the user attempts to log in
    Then the system should indicate an appropriate error message to the user
