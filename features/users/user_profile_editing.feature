@todo
Feature: Profile Editing
  As a User
  I want to be able to edit my user profile data (e.g., name, email, profile photo, and other relevant details)
  So that I can keep my personal information updated and reflect my changes over time

  Background:
    Given the user email "user@domain.com" and password "password"
    And the registration is successful

  Scenario: Successful data editing
    When the user modifies and saves their profile data correctly
    Then the system should update the information and show a success confirmation

  Scenario: Field validation
    When the user modifies with invalid information or leaves required fields empty
    Then the system should show clear error messages indicating what data needs correction before saving

  Scenario: Real-time feedback
    When the user makes changes to their profile
    Then the system should show a preview or confirmation of the changes made before the final update

  Scenario: Server error handling
    When the user modifies and saves their profile data correctly
    And an error occurs during the update process
    Then the system should show an informative error message and suggest retrying the operation
