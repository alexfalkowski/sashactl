@clear_pool
Feature: New article
  Create articles

  @operational
  Scenario: Successfully create an article with a name
    Given I start the system
    When we create an article with name "new article"
    Then it should run successfully
    And I should have an article named "new article"
    And I should see a log entry of "created article" in the file "reports/new.log"

  @missing
  Scenario: Unsuccessfully create an article with a name as the config is missing
    Given I start the system
    When we create an article with name "new article"
    Then it should not run successfully
    And I should not have an article named "new article"
    And I should see a log entry of "not found" in the file "reports/new.log"

  @erroneous
  Scenario: Unsuccessfully create an article with a name as the config is broken
    Given I start the system
    When we create an article with name "new article"
    Then it should not run successfully
    And I should not have an article named "new article"
    And I should see a log entry of "internal server error" in the file "reports/new.log"

  @operational
  Scenario: Unsuccessfully create an article with a name as the config is down
    Given I start the system
    And I set the proxy for server "bucket" to "close_all"
    When we create an article with name "new article"
    Then it should not run successfully
    And I should not have an article named "new article"
    And I should reset the proxy for server "bucket"
