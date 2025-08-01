@todo
Feature: Health check
  As a system admin
  I want to check the health of the app
  So that I can monitor it's status

  Scenario: Check app health
    Given the app is running
    When I request the health status
    Then I should receive a healthy status
