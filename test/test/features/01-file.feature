Feature: User can upload file
  As User
  I can post picture

  Scenario: I can post new picture file and look he file
    Given tokenFromRegister exists
    When I prepare a POST request to "/v1/file"
    And I set a bearer authentication header with tokenFromRegister value
    And I use formdata filePas payload
    And send
    Then the response should be 200
    And save the "uri" value as profileImageUri

    Given profileImageUri exists
    When I prepare a GET request to saved URI profileImageUri
    And I set a bearer authentication header with tokenFromRegister value
    And send
    Then the response should be 200
    And the file should be equal as filePas

  Scenario: I cannot post huge picture
    Given tokenFromRegister exists
    When I prepare a POST request to "/v1/file"
    And I set a bearer authentication header with tokenFromRegister value
    And I use formdata fileHuge payload
    And send
    Then the response should be 400

