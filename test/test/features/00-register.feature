Feature: User register to Fitbyte Service
  As a user
  I want to register and log in
  So that I can access and manage my activity

  Scenario: Successful admin registration
    When I prepare a POST request to "/v1/register"
    And I use register payload
    And send
    Then the response should be 201
    And the response body matches the auth schema
    And save the "email" value as emailFromRegister
    And save the "token" value as tokenFromRegister

  # I prepare another account
    When I use registerAnother payload
    And send
    Then the response should be 201
    And the response body matches the auth schema
    And save the "token" value as tokenNegativeScenario

  Scenario: Registration fails with invalid payloads
    When I prepare a POST request to "/v1/register"
    And I inject register violations
    And send violations
    Then the violation responses should be 400

  Scenario: Registration fails with username conflict
    When I prepare a POST request to "/v1/register"
    And I use register payload
    And send
    Then the response should be 409

  Scenario: Successful admin login
    When I prepare a POST request to "/v1/login"
    And I use login payload
    And send
    Then the response should be 200
    And the response body matches the auth schema

  # Negatives

  Scenario: Login fails with wrong password
    When I prepare a POST request to "/v1/login"
    And I use loginInvalidPass payload
    And send
    Then the response should be 400

  Scenario: Login fails with invalid payload
    When I prepare a POST request to "/v1/login"
    And I inject login violations
    And send violations
    Then the violation responses should be 400

