@delete
Feature: Delete articles
  Delete existing articles

  Scenario: Successfully delete an article with a slug
    Given we have a published article with slug "1984"
    When we delete an article with slug "1984"
    Then it should run successfully
    And I should have a deleted article with slug "1984"
    And I should see a log entry of "deleted article" in the file "reports/delete.log"
    And the article with slug "1984" should be removed from the file system

  Scenario: Unsuccessfully delete an article with a slug when we have issues with the bucket
    Given we have a published article with slug "1984"
    And I set the proxy for service "aws" to "close_all"
    When we delete an article with slug "1984"
    Then it should not run successfully
    And I should reset the proxy for service "aws"
