@todo
Feature: Login with Biometric Data
  As a User
  I want to log into my account on the platform using my biometric data
  So that I can access the system securely and conveniently

  Background:
    Given the user email "user@domain.com" and password "password"
    And the registration is successful

  Scenario: Successful login with biometric data
    When the user logs in using biometric data
      # such as face or fingerprint recognition
    And the biometric log in is correct
    Then the system should allow access to the system

  Scenario: Failed biometric login
    Given the user login using biometric data is incorrect
    And the biometric log in is correct
    Then an appropriate error message should be shown to the user

  Scenario: Login failure
    When the user logs in using biometric data
    And it fails due to a service error
    Then the system should indicate an appropriate error message to the user

  Scenario: Session expiration
    When the session exceeds the expiration time
    Then the system should request the user to log in again

  Scenario: Blocked user
    Given the user account is blocked by administrators
    When the user logs in using biometric data
    Then the system should not allow it to log in
