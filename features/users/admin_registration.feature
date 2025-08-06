@todo
Feature: Administrator Registration
  As an administrator
  I want to be able to register other administrators on the platform
  So that they can access administrative functionalities

  Background:
    Given the admin email "admin@domain.com" and password "password"
    And the registration is successful

  Scenario: Successful registration
    When the admin register the email "new-admin@domain.com"
    Then the system should register a new administrator with their respective credentials
    And allow access to the system

  Scenario: Missing or invalid required fields
    Given the administrator registration is incorrect due to missing or invalid field entry
    When I attempt to register
    Then the system should indicate the missing and/or incorrect fields
    And not allow access to the system

  Scenario: Registration failure
    When the administrator registration fails due to a service error
    Then the system should indicate an appropriate error message to the user
