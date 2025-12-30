Feature: User manage their account
  As a User
  I can choose my preference

  Scenario: When I open my account, I should see an empty preferences
    Given tokenFromRegister exists
    And emailFromRegister exists
    When I prepare a GET request to "/v1/user"
    And I set a bearer authentication header with tokenFromRegister value
    And send
    Then the response should be 200
    And the response body matches the user schema
    And "preference" value is null
    And "weightUnit" value is null
    And "heightUnit" value is null
    And "weight" value is null
    And "height" value is null
    And "email" value is equal to saved emailFromRegister
    And "name" value is null
    And "imageUri" value is null

  Scenario: I add my account preferences and see once again the changes
    Given tokenFromRegister exists
    When I prepare a PATCH request to "/v1/user"
    And I use userPreference payload
    And I set a bearer authentication header with tokenFromRegister value
    And using userPreference payload, I add profileImageUri to "imageUri"
    And send
    Then the response should be 200
    And the response body matches the user schema
    And "preference" value is equal to "WEIGHT"
    And "weightUnit" value is equal to "LBS"
    And "heightUnit" value is equal to "INCH"
    And "weight" value is equal to 75
    And "height" value is equal to 167
    And "email" value is equal to saved emailFromRegister
    And "imageUri" value is equal to saved profileImageUri
    And save the "name" value as nameFromPatchUser

    Given emailFromRegister exists
    When I prepare a GET request to "/v1/user"
    And I set a bearer authentication header with tokenFromRegister value
    And send
    Then the response should be 200
    And the response body matches the user schema
    And "preference" value is equal to "WEIGHT"
    And "weightUnit" value is equal to "LBS"
    And "heightUnit" value is equal to "INCH"
    And "weight" value is equal to 75
    And "height" value is equal to 167
    And "email" value is equal to saved emailFromRegister
    And "name" value is equal to saved nameFromPatchUser
    And "imageUri" value is equal to saved profileImageUri

    # negatives cases
  Scenario: I didn't have account, so I should not able to see my account
    When I prepare a GET request to "/v1/user"
    And I set a bearer authentication header with false token
    And send
    Then the response should be 401

  Scenario: I didn't have account, so I should not able to patch my account
    When I prepare a PATCH request to "/v1/user"
    And I set a bearer authentication header with false token
    And send
    Then the response should be 401

  Scenario: I have an account, but inject negative payload
    Given tokenFromRegister exists
    When I prepare a PATCH request to "/v1/user"
    And I inject userPreference violations
    And I set a bearer authentication header with tokenFromRegister value
    And using userPreference payload, I add profileImageUri to "imageUri"
    And send violations
    Then the response should be 400






