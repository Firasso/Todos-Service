Feature: add todo
  In order to add a todo
  As an API user
  I need to be able to post a todo

  Scenario: should add a todo
    When I request REST endpoint with method "POST" and path "/todos" with body
    """
    {
      "text":"added by dog",
      "completed":true
    }
    """
    Then there should be no error
    And the response status code should be "201"