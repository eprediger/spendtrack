@todo
Feature: PIN Notification in Registration Process
  As a user who has registered on the platform
  I want to receive a verification PIN via email or WhatsApp/SMS
  So that I can confirm my account and validate the authenticity of my registration securely

  Background:
    Given the user email "user@domain.com" and password "password"
    And the registration is successful

  Scenario: PIN verification sending
    Then the system should generate a unique PIN and send it to the selected medium (email, WhatsApp or SMS)
    And show an informative message indicating that the PIN has been sent

  Scenario: Valid PIN validation
    When the user enters the received PIN in the verification form
    Then the system validates if the PIN is correct and entered within the expiration time (e.g., 10 minutes)
    Then the account should be activated and the user should be able to access the platform

  Scenario: Incorrect PIN validation
    When the user enters the received PIN in the verification form
    Then the system validates if the PIN is correct and entered within the expiration time (e.g., 10 minutes)
    Then  if the PIN is incorrect or has expired, an error message should be shown

  Scenario: Expired PIN validation
    When the user enters the received PIN in the verification form
    Then the system validates if the PIN is correct and entered within the expiration time (e.g., 10 minutes)
    Then the option to resend a new PIN should be offered

  Scenario: PIN security and single use
    When the PIN is created
    Then it should be valid for only one verification attempt and have a defined expiration time
    And once used or expired, the system should invalidate the PIN to prevent reuse and ensure the process
