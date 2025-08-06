@todo
Feature: Password Recovery
  As a User
  I want to be able to recover my password in case I don't remember it
  So that I can regain access to my account

  Background:
    Given the user email "user@domain.com" and password "password"
    And the registration is successful

  Scenario: Successful password recovery
    When the user asks for a password recovery
    Then an email should be sent with a password recovery link
    And the link should have only one use and an expiration time

  Scenario: Recovery failure
    Given the recovery time has expired
    When the user asks for a password recovery
    Then the system should not be able to perform the recovery
    And the user should start the process again
