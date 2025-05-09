Feature: Publish articles
  Publish created articles

  Scenario: Successfully publish an article with a slug
    When we publish an article with slug "1984"
    Then it should run successfully
    And I should have a published article with slug "1984"
    And I should see a log entry of "published article" in the file "reports/publish.log"

  Scenario: Unsuccessfully publish an article with a slug when we have issues with the bucket
    Given I set the proxy for service "aws" to "close_all"
    When we publish an article with slug "1984"
    Then it should not run successfully
    And I should not have a published article with slug "1984"
    And I should reset the proxy for service "aws"
