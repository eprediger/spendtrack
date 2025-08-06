@todo
Feature: Other Users Profile Viewing
  As a User
  I want to be able to view other users' data
  So that I can know their contact information and public details

  Background:
    Given the user email "user@domain.com" and password "password"
    And the registration is successful
    And the user email "other-user@domain.com" and password "password"
    And the registration is successful
    And its logged in

  Scenario: Viewing public data of other users
    Given the user "user@domain.com" accesses "other-user@domain.com" user's profile
    When it views their profile
    Then the system should show only the public data configured by the profile owner
    And display information such as name, profile photo, and biography

  Scenario: Restriction of private data
    Given a user's profile is being viewed
    When the profile loads
    Then the system should not show private or sensitive information that is not intended to be public

  Scenario: Handling non-existent profiles
    Given the user attempts to access a profile that does not exist or has been deleted
    When I try to view the profile
    Then the system should show an appropriate message indicating that the profile is not available
