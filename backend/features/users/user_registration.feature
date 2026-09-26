Feature: User Registration
  As a user
  I want to register in the platform
  So that I can access the system functionalities

  @todo
  Scenario: Successful registration
    Given the user email "user@domain.com" and password "password"
    When the registration is successful
    Then the system will create a new user with their profile
#    And allow access to the system

  @todo
  Scenario: Location request
    Given the user registration is correct
    Then the system should request a default location from the user

  @todo
  Scenario: Missing fields
    Given the user registration is incorrect due to missing field entry
    When I attempt to register
    Then the system should indicate the missing and/or incorrect fields
    And not allow access to the system

  @todo
  Scenario: Invalid fields
    Given the user registration is incorrect due to invalid field entry
    When I attempt to register
    Then the system should indicate the missing and/or incorrect fields
    And not allow access to the system

  @todo
  Scenario: Duplicated user email
    Given a registered user with email "user@domain.com" and password "password" exists
    When a new user registers with email "user@domain.com"
    Then the system should respond a "REGISTRATION_ERROR" message

  @todo
  Scenario: Failed registration
    Given the user email "user@domain.com" and password "password"
    When the registration fails due to a system error
    Then the system should respond a "REGISTRATION_ERROR" message
