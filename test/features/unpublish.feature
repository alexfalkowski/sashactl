@unpublish
Feature: Unpublish articles
  Unpublish existing articles

  Scenario: Successfully unpublish an article with a slug
    Given we have a published article with slug "1984"
    When we unpublish an article with slug "1984"
    Then it should run successfully
    And I should have an unpublished article with slug "1984"
    And I should see a log entry of "unpublished article" in the file "reports/unpublish.log"
    And the unpublished article with slug "1984" should be removed from the file system

  Scenario: Unsuccessfully unpublish an article with a slug when we have issues with the bucket
    Given we have a published article with slug "1984"
    And I set the proxy for service "aws" to "close_all"
    When we unpublish an article with slug "1984"
    Then it should not run successfully
    And I should reset the proxy for service "aws"
