@todo
Feature: User Login with Username and Password
  As a User
  I want to log into my account on the platform using username and password
  So that I can access the system functionalities

  Background:
    Given the user email "user@domain.com" and password "password"
    And the registration is successful

  Scenario: Successful login with email and password
    Given the user login using email and password is correct
    When the user attempts to log in
    Then the system should allow access to the system

  Scenario: Failed login with incorrect credentials
    Given the user login using email and password is incorrect
    When the user attempts to log in
    Then an appropriate error message should be shown to the user

  Scenario: Login failure
    Given the user login fails due to a service error
    When the user attempts to log in
    Then the system should indicate an appropriate error message to the user

  Scenario: Session expiration
    When the user session exceeds the expiration time
    Then the system should request the user to log in again

  Scenario: Blocked user
    Given the user account is blocked by administrators
    When the user attempts to log in
    Then the system should not allow the user to log in

  Scenario: Account lockout due to failed attempts
    Given the user login using email and password is incorrect
    When the system detects multiple failed login attempts in a short period of time
    Then the system should temporarily block the user account and show a warning message
    And indicate the need to wait or contact support to unlock it
