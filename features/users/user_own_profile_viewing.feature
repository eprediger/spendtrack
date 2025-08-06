@todo
Feature: Own Profile Viewing
  As a User
  I want to be able to view my user profile data
  So that I can review my personal information and confirm it is correct and update it if necessary

  Background:
    Given the user email "user@domain.com" and password "password"
    And the registration is successful
    And its logged in

  Scenario: Correct profile visualization
    Given the user accesses the "My Profile" section
    When the user attempts to view its profile
    Then the system should show the profile information in an organized manner
    And display both public and private data associated with the account

  Scenario: Error handling
    Given an error occurs while loading the profile (e.g., connection failure)
    When the user attempts to view its profile
    Then the system should show an appropriate error message and indicate options to retry loading
