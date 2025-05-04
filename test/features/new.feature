@new
Feature: New article
  Create articles

  Scenario: Successfully create an article with a name
    When we create an article with name "new article"
    Then it should run successfully
    And I should have an article named "new article"
    And I should see a log entry of "created article" in the file "reports/new.log"

  Scenario: Successfully create an article with a name when the config is missing
    When we create an article with name "new article"
    Then it should run successfully
    And I should have an article named "new article"
    And I should see a log entry of "created article" in the file "reports/new.log"

  Scenario: Unsuccessfully create an article with a name when we have issues with the bucket
    And I set the proxy for service "aws" to "close_all"
    When we create an article with name "new article"
    Then it should not run successfully
    And I should not have an article named "new article"
    And I should reset the proxy for service "aws"
