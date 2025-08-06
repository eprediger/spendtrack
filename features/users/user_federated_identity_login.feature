@todo
Feature: User Login with Federated Identity Providers
  As a user
  I want to log into the platform using federated identity providers (e.g., Google, Facebook, etc.)
  So that I can access quickly and securely without manually registering my data

  Background:
    Given the user wants to log into the platform through a federated identity provider

  Scenario: Successful login through federated identity
    When the user selects a federated identity provider
    And the provider authorizes the access
    Then the system receives a valid authentication token
    And the system should start my session
    And I should be taken to the platform

  Scenario Outline: Authentication error handling
    When the provider fails due to "<error>"
    Then the system should display a clear error message indicating authentication failed
    And offer me the option to retry or use another authentication method

    Examples:
      | error               |
      | invalid credentials |
      | connection issues   |

  Scenario: Profile data synchronization
    When the user logs in through a federated provider
    Then the system should extract and synchronize basic profile information
    And create or update the user account with name, email, and profile photo
    And allow the user to modify this data later if desired

  Scenario: Maintaining user system uniqueness
    Given the user previously registered with email and password
    When it logs in using a federated provider with the same email
    Then the system should ask if it wants to combine the information
    And request its original credentials for verification
    And allow it to log in using both methods in the future

  Scenario: Secure authentication handling
    Given the user logs in using a federated provider
    When the authentication is successful
    Then the system should securely process my authentication
    And maintain the integrity and confidentiality of my session
